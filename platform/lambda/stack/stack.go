// Package stack provides CloudFormation stack support.
package stack

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"

	"github.com/apex/log"
	"github.com/apex/up"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/event"
)

// TODO: refactor a lot
// TODO: backoff
// TODO: profile changeset name and description flags
// TODO: flags for changeset name / description

// defaultChangeset name.
var defaultChangeset = "changes"

// Stack represents a single CloudFormation stack.
type Stack struct {
	client  *cloudformation.CloudFormation
	lambda  *lambda.Lambda
	route53 *route53.Route53
	events  event.Events
	zones   []*route53.HostedZone
	config  *up.Config
}

// New stack.
func New(c *up.Config, events event.Events, zones []*route53.HostedZone, region string) *Stack {
	sess := session.New(aws.NewConfig().WithRegion(region))
	return &Stack{
		client:  cloudformation.New(sess),
		lambda:  lambda.New(sess),
		route53: route53.New(sess),
		events:  events,
		zones:   zones,
		config:  c,
	}
}

// template returns a configured resource template.
func (s *Stack) template() Map {
	return template(&Config{
		Config: s.config,
		Zones:  s.zones,
	})
}

// Create the stack.
func (s *Stack) Create(version string) error {
	c := s.config
	tmpl := s.template()
	name := c.Name

	b, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	_, err = s.client.CreateStack(&cloudformation.CreateStackInput{
		StackName:        &name,
		TemplateBody:     aws.String(string(b)),
		TimeoutInMinutes: aws.Int64(60),
		DisableRollback:  aws.Bool(true),
		Capabilities:     aws.StringSlice([]string{"CAPABILITY_NAMED_IAM"}),
		Parameters: []*cloudformation.Parameter{
			{
				ParameterKey:   aws.String("Name"),
				ParameterValue: &name,
			},
			{
				ParameterKey:   aws.String("FunctionName"),
				ParameterValue: &name,
			},
			{
				ParameterKey:   aws.String("FunctionVersionStaging"),
				ParameterValue: &version,
			},
			{
				ParameterKey:   aws.String("FunctionVersionProduction"),
				ParameterValue: &version,
			},
		},
	})

	if err != nil {
		return errors.Wrap(err, "creating stack")
	}

	if err := s.report(resourceStateFromTemplate(tmpl, CreateComplete)); err != nil {
		return errors.Wrap(err, "reporting")
	}

	stack, err := s.getStack()
	if err != nil {
		return errors.Wrap(err, "fetching stack")
	}

	status := Status(*stack.StackStatus)
	if status.State() == Failure {
		return errors.New(*stack.StackStatusReason)
	}

	return nil
}

// Delete the stack, optionally waiting for completion.
func (s *Stack) Delete(wait bool) error {
	_, err := s.client.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: &s.config.Name,
	})

	if err != nil {
		return errors.Wrap(err, "deleting")
	}

	if wait {
		tmpl := s.template()
		if err := s.report(resourceStateFromTemplate(tmpl, DeleteComplete)); err != nil {
			return errors.Wrap(err, "reporting")
		}
	}

	return nil
}

// Show resources.
func (s *Stack) Show() error {
	defer s.events.Time("platform.stack.show", nil)()

	// show stack status
	stack, err := s.getStack()
	if err != nil {
		return errors.Wrap(err, "fetching stack")
	}

	s.events.Emit("platform.stack.show.stack", event.Fields{
		"stack": stack,
	})

	// show nameservers
	if err := s.showNameservers(); err != nil {
		return errors.Wrap(err, "showing nameservers")
	}

	// skip events if everything is ok
	if Status(*stack.StackStatus).State() == Success {
		return nil
	}

	// show events
	events, err := s.getFailedEvents()
	if err != nil {
		return errors.Wrap(err, "fetching latest events")
	}

	for _, e := range events {
		if *e.LogicalResourceId == s.config.Name {
			continue
		}

		s.events.Emit("platform.stack.show.stack.event", event.Fields{
			"event": e,
		})
	}

	return nil
}

// Plan changes.
func (s *Stack) Plan() error {
	c := s.config
	tmpl := s.template()
	name := c.Name

	b, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	defer s.events.Time("platform.stack.plan", nil)

	prod, err := s.lambda.GetAlias(&lambda.GetAliasInput{
		FunctionName: &name,
		Name:         aws.String("production"),
	})

	if err != nil {
		return errors.Wrap(err, "fetching production alias")
	}

	stage, err := s.lambda.GetAlias(&lambda.GetAliasInput{
		FunctionName: &name,
		Name:         aws.String("staging"),
	})

	if err != nil {
		return errors.Wrap(err, "fetching staging alias")
	}

	log.Debug("deleting changeset")
	_, err = s.client.DeleteChangeSet(&cloudformation.DeleteChangeSetInput{
		StackName:     &name,
		ChangeSetName: &defaultChangeset,
	})

	if err != nil {
		return errors.Wrap(err, "deleting changeset")
	}

	log.Debug("creating changeset")
	_, err = s.client.CreateChangeSet(&cloudformation.CreateChangeSetInput{
		StackName:     &name,
		ChangeSetName: &defaultChangeset,
		TemplateBody:  aws.String(string(b)),
		Capabilities:  aws.StringSlice([]string{"CAPABILITY_NAMED_IAM"}),
		ChangeSetType: aws.String("UPDATE"),
		Description:   aws.String("Managed by Up."),
		Parameters: []*cloudformation.Parameter{
			{
				ParameterKey:   aws.String("Name"),
				ParameterValue: &name,
			},
			{
				ParameterKey:   aws.String("FunctionName"),
				ParameterValue: &name,
			},
			{
				ParameterKey:   aws.String("FunctionVersionStaging"),
				ParameterValue: stage.FunctionVersion,
			},
			{
				ParameterKey:   aws.String("FunctionVersionProduction"),
				ParameterValue: prod.FunctionVersion,
			},
		},
	})

	if err != nil {
		return errors.Wrap(err, "creating changeset")
	}

	var next *string

	for {
		log.Debug("describing changeset")
		res, err := s.client.DescribeChangeSet(&cloudformation.DescribeChangeSetInput{
			StackName:     &name,
			ChangeSetName: &defaultChangeset,
			NextToken:     next,
		})

		if err != nil {
			return errors.Wrap(err, "describing changeset")
		}

		status := Status(*res.Status)

		if status.State() == Failure {
			if _, err := s.client.DeleteChangeSet(&cloudformation.DeleteChangeSetInput{
				StackName:     &name,
				ChangeSetName: &defaultChangeset,
			}); err != nil {
				return errors.Wrap(err, "deleting changeset")
			}

			return errors.New(*res.StatusReason)
		}

		if !status.IsDone() {
			log.Debug("waiting for completion")
			time.Sleep(750 * time.Millisecond)
			continue
		}

		for _, c := range res.Changes {
			s.events.Emit("platform.stack.plan.change", event.Fields{
				"change": c,
			})
		}

		next = res.NextToken

		if next == nil {
			break
		}
	}

	return nil
}

// Apply changes.
func (s *Stack) Apply() error {
	c := s.config
	name := c.Name

	res, err := s.client.DescribeChangeSet(&cloudformation.DescribeChangeSetInput{
		StackName:     &name,
		ChangeSetName: &defaultChangeset,
	})

	if isNotFound(err) {
		return errors.Errorf("changeset does not exist, run `up stack plan` first")
	}

	if err != nil {
		return errors.Wrap(err, "describing changeset")
	}

	defer s.events.Time("platform.stack.apply", event.Fields{
		"changes": len(res.Changes),
	})()

	_, err = s.client.ExecuteChangeSet(&cloudformation.ExecuteChangeSetInput{
		StackName:     &name,
		ChangeSetName: &defaultChangeset,
	})

	if err != nil {
		return errors.Wrap(err, "executing changeset")
	}

	if err := s.report(resourceStateFromChanges(res.Changes)); err != nil {
		return errors.Wrap(err, "reporting")
	}

	return nil
}

// report events with a map of desired stats from logical or physical id,
// any resources not mapped are ignored as they do not contribute to changes.
func (s *Stack) report(states map[string]Status) error {
	defer s.events.Time("platform.stack.report", event.Fields{
		"total":    len(states),
		"complete": 0,
	})()

	for range time.Tick(time.Second) {
		stack, err := s.getStack()

		if util.IsNotFound(err) {
			return nil
		}

		if util.IsThrottled(err) {
			time.Sleep(3 * time.Second)
			continue
		}

		if err != nil {
			return errors.Wrap(err, "fetching stack")
		}

		status := Status(*stack.StackStatus)

		if status.IsDone() {
			return nil
		}

		res, err := s.client.DescribeStackResources(&cloudformation.DescribeStackResourcesInput{
			StackName: &s.config.Name,
		})

		if util.IsThrottled(err) {
			time.Sleep(time.Second * 3)
			continue
		}

		if err != nil {
			return errors.Wrap(err, "describing stack resources")
		}

		complete := len(resourcesCompleted(res.StackResources, states))

		s.events.Emit("platform.stack.report.event", event.Fields{
			"total":    len(states),
			"complete": complete,
		})
	}

	return nil
}

// showNameservers emits events for listing name servers.
func (s *Stack) showNameservers() error {
	s.events.Emit("platform.stack.show.nameservers", nil)

	for _, stage := range s.config.Stages.List() {
		if stage.Domain == "" {
			continue
		}

		res, err := s.route53.ListHostedZonesByName(&route53.ListHostedZonesByNameInput{
			DNSName:  &stage.Domain,
			MaxItems: aws.String("1"),
		})

		if err != nil {
			return errors.Wrap(err, "listing hosted zone")
		}

		if len(res.HostedZones) == 0 {
			continue
		}

		z := res.HostedZones[0]
		if stage.Domain+"." != *z.Name {
			continue
		}

		zone, err := s.route53.GetHostedZone(&route53.GetHostedZoneInput{
			Id: z.Id,
		})

		if err != nil {
			return errors.Wrap(err, "fetching hosted zone")
		}

		defer s.events.Time("platform.stack.show.stage", event.Fields{
			"stage": stage,
		})()

		for _, ns := range zone.DelegationSet.NameServers {
			s.events.Emit("platform.stack.show.nameserver", event.Fields{
				"nameserver": *ns,
			})
		}
	}

	return nil
}

// getStack returns the stack.
func (s *Stack) getStack() (*cloudformation.Stack, error) {
	res, err := s.client.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: &s.config.Name,
	})

	if err != nil {
		return nil, err
	}

	stack := res.Stacks[0]
	return stack, nil
}

// getLatestEvents returns the latest events for each resource.
func (s *Stack) getLatestEvents() (v []*cloudformation.StackEvent, err error) {
	events, err := s.getEvents()
	if err != nil {
		return
	}

	hit := make(map[string]bool)

	for _, e := range events {
		id := *e.LogicalResourceId
		if hit[id] {
			continue
		}

		hit[id] = true
		v = append(v, e)
	}

	return
}

// getFailedEvents returns failed events.
func (s *Stack) getFailedEvents() (v []*cloudformation.StackEvent, err error) {
	events, err := s.getEvents()
	if err != nil {
		return
	}

	for _, e := range events {
		if Status(*e.ResourceStatus).State() == Failure {
			v = append(v, e)
		}
	}

	return
}

// getEvents returns events.
func (s *Stack) getEvents() (events []*cloudformation.StackEvent, err error) {
	var next *string

	for {
		res, err := s.client.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
			StackName: &s.config.Name,
			NextToken: next,
		})

		if err != nil {
			return nil, err
		}

		events = append(events, res.StackEvents...)

		next = res.NextToken

		if next == nil {
			break
		}
	}

	return
}

// resourceStateFromTemplate returns a map of the logical ids from template t, to status s.
func resourceStateFromTemplate(t Map, s Status) map[string]Status {
	r := t["Resources"].(Map)
	m := make(map[string]Status)

	for id := range r {
		m[id] = s
	}

	return m
}

// TODO: ignore deletes since they're in cleanup phase?

// resourceStateFromChanges returns a map of statuses from a changeset.
func resourceStateFromChanges(changes []*cloudformation.Change) map[string]Status {
	m := make(map[string]Status)

	for _, c := range changes {
		var state Status
		var id string

		if s := c.ResourceChange.PhysicalResourceId; s != nil {
			id = *s
		}

		if id == "" {
			id = *c.ResourceChange.LogicalResourceId
		}

		switch a := *c.ResourceChange.Action; a {
		case "Add":
			state = CreateComplete
		case "Modify":
			state = UpdateComplete
		case "Remove":
			state = DeleteComplete
		default:
			panic(errors.Errorf("unhandled Action %q", a))
		}

		m[id] = state
	}

	return m
}

// resourcesCompleted returns a map of the completed resources. When the resource is not
// present in states, it is ignored as no changes are expected.
func resourcesCompleted(resources []*cloudformation.StackResource, states map[string]Status) map[string]*cloudformation.StackResource {
	m := make(map[string]*cloudformation.StackResource)

	for _, r := range resources {
		var expected Status
		var id string

		// try physical id first, this is necessary as
		// replacement of a logical id will cause the id
		// to appear twice (once for Add once for Remove).
		if s := r.PhysicalResourceId; s != nil {
			if _, ok := states[*s]; ok {
				id = *s
			}
		}

		// try logical id
		if s := *r.LogicalResourceId; id == "" {
			if _, ok := states[s]; ok {
				id = s
			}
		}

		// expected state
		if id != "" {
			expected = states[id]
		}

		// matched expected state
		if expected == Status(*r.ResourceStatus) {
			m[id] = r
		}
	}

	return m
}

// isNotFound returns true if the error indicates a missing changeset.
func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "ChangeSetNotFound")
}
