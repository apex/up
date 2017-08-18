package stack

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"

	"github.com/apex/log"
	"github.com/apex/up"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform/event"
)

// TODO: refactor a lot
// TODO: backoff
// TODO: profile changeset name and description flags

// defaultChangeset name.
var defaultChangeset = "changes"

// Stack represents a single CloudFormation stack.
type Stack struct {
	client *cloudformation.CloudFormation
	events event.Events
	config *up.Config
}

// New stack.
func New(c *up.Config, events event.Events, region string) *Stack {
	sess := session.New(aws.NewConfig().WithRegion(region))
	return &Stack{
		client: cloudformation.New(sess),
		events: events,
		config: c,
	}
}

// Create the stack.
func (s *Stack) Create(version string) error {
	c := s.config
	tmpl := template(c)
	name := c.Name

	b, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	_, err = s.client.CreateStack(&cloudformation.CreateStackInput{
		StackName:        &name,
		TemplateBody:     aws.String(string(b)),
		TimeoutInMinutes: aws.Int64(15),
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
				ParameterKey:   aws.String("FunctionVersion"),
				ParameterValue: &version,
			},
		},
	})

	if err != nil {
		return errors.Wrap(err, "creating stack")
	}

	if err := s.report("create"); err != nil {
		return errors.Wrap(err, "reporting events")
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
		if err := s.report("delete"); err != nil {
			return errors.Wrap(err, "reporting")
		}
	}

	return nil
}

// Show resources.
func (s *Stack) Show() error {
	defer s.events.Time("platform.stack.show", nil)()

	stack, err := s.getStack()
	if err != nil {
		return errors.Wrap(err, "fetching stack")
	}

	s.events.Emit("platform.stack.show.stack", event.Fields{
		"stack": stack,
	})

	events, err := s.getLatestEvents()
	if err != nil {
		return errors.Wrap(err, "fetching latest events")
	}

	for _, e := range events {
		s.events.Emit("platform.stack.show.event", event.Fields{
			"event": e,
		})
	}

	return nil
}

// Plan changes.
func (s *Stack) Plan() error {
	c := s.config
	tmpl := template(c)
	name := c.Name

	b, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	defer s.events.Time("platform.stack.plan", nil)

	// TODO: don't change deployments

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
		Description:   aws.String("Managed by Up."), // TODO: flag?
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
				ParameterKey:   aws.String("FunctionVersion"),
				ParameterValue: aws.String("104"),
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

	return s.report("update")
}

// report events.
func (s *Stack) report(state string) error {
	hit := make(map[string]bool)
	tmpl := template(s.config)

	defer s.events.Time("platform.stack."+state, event.Fields{
		"resources": len(tmpl["Resources"].(Map)),
	})()

	for range time.Tick(time.Second) {
		stack, err := s.getStack()

		if util.IsNotFound(err) {
			return nil
		}

		if err != nil {
			return errors.Wrap(err, "fetching stack")
		}

		status := Status(*stack.StackStatus)

		if status.IsDone() {
			return nil
		}

		events, err := s.getEvents()

		if util.IsNotFound(err) {
			return nil
		}

		if err != nil {
			return errors.Wrap(err, "fetching events")
		}

		for _, e := range events {
			if hit[*e.EventId] {
				continue
			}
			hit[*e.EventId] = true

			s.events.Emit("platform.stack."+state+".event", event.Fields{
				"event": e,
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

		for _, e := range res.StackEvents {
			events = append(events, e)
		}

		next = res.NextToken

		if next == nil {
			break
		}
	}

	return
}

// getEventsByState returns events by state.
func (s *Stack) getEventsByState(state State) (v []*cloudformation.StackEvent, err error) {
	events, err := s.getEvents()
	if err != nil {
		return
	}

	for _, e := range events {
		s := Status(*e.ResourceStatus)
		if s.State() == state {
			v = append(v, e)
		}
	}

	return
}

// hasChange returns true if id matches a physical or logical id.
func hasChange(id string, changes []*cloudformation.Change) bool {
	for _, c := range changes {
		cid := *c.ResourceChange.LogicalResourceId

		if s := c.ResourceChange.PhysicalResourceId; s != nil {
			cid = *s
		}

		if cid == id {
			return true
		}
	}

	return false
}
