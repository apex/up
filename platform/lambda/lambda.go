package lambda

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/apex/apex/shim"
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"

	"github.com/apex/up"
	"github.com/apex/up/config"
	"github.com/apex/up/internal/proxy/bin"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/zip"
	"github.com/apex/up/platform"
	"github.com/apex/up/platform/event"
	"github.com/apex/up/platform/lambda/stack"
)

// errFirstDeploy is returned from .deploy() when a function is created.
var errFirstDeploy = errors.New("first deploy")

// assume policy for the lambda function.
var apiGatewayAssumePolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Service": "apigateway.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		},
		{
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
	]
}`

// policy for the lambda function.
var logWritePolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Resource": "*",
			"Action": [
				"logs:CreateLogGroup",
				"logs:CreateLogStream",
				"logs:PutLogEvents"
			]
		}
	]
}`

// TODO: concurrency option
// TODO: aggregate progress report for N regions or distinct progress bars
// TODO: refactor with another region-scoped struct to clean this up

// Platform implementation.
type Platform struct {
	config  *up.Config
	runtime string
	handler string
	zip     *bytes.Buffer
	events  event.Events
}

// New platform.
func New(c *up.Config, events event.Events) *Platform {
	return &Platform{
		config:  c,
		runtime: "nodejs6.10",
		handler: "_proxy.handle",
		events:  events,
	}
}

// Build implementation.
func (p *Platform) Build() error {
	start := time.Now()
	p.zip = new(bytes.Buffer)

	if err := p.injectProxy(); err != nil {
		return errors.Wrap(err, "injecting proxy")
	}
	defer p.removeProxy()

	r, stats, err := zip.Build(".")
	if err != nil {
		return errors.Wrap(err, "building")
	}

	if _, err := io.Copy(p.zip, r); err != nil {
		return errors.Wrap(err, "copying")
	}

	if err := r.Close(); err != nil {
		return errors.Wrap(err, "closing")
	}

	p.events.Emit("platform.build.zip", event.Fields{
		"files":             stats.FilesAdded,
		"size_uncompressed": stats.SizeUncompressed,
		"size_compressed":   p.zip.Len(),
		"duration":          time.Since(start),
	})

	return nil
}

// Deploy implementation.
func (p *Platform) Deploy(stage string) error {
	regions := p.config.Regions
	errc := make(chan error, len(regions))

	if err := p.createRole(); err != nil {
		return errors.Wrap(err, "iam")
	}

	for _, r := range regions {
		go func(region string) {
			version, err := p.deploy(region, stage)

			if err == nil {
				errc <- nil
				return
			}

			if err != errFirstDeploy {
				errc <- errors.Wrap(err, region)
				return
			}

			if err := p.CreateStack(region, version); err != nil {
				errc <- errors.Wrap(err, region)
				return
			}

			errc <- nil
		}(r)
	}

	for range p.config.Regions {
		if e := <-errc; e != nil {
			return e
		}
	}

	return nil
}

// Logs implementation.
func (p *Platform) Logs(region, query string) platform.Logs {
	return NewLogs(p, region, query)
}

// URL returns the stage url.
func (p *Platform) URL(region, stage string) (string, error) {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := apigateway.New(s)

	api, err := p.getAPI(c)
	if err != nil {
		return "", errors.Wrap(err, "fetching api")
	}

	if api == nil {
		return "", errors.Errorf("cannot find the API, looks like you haven't deployed")
	}

	id := fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/%s/", *api.Id, region, stage)
	return id, nil
}

// CreateStack implementation.
func (p *Platform) CreateStack(region, version string) error {
	return stack.New(p.config, p.events, region).Create(version)
}

// DeleteStack implementation.
func (p *Platform) DeleteStack(region string, wait bool) error {
	log.Debug("deleting stack")
	if err := stack.New(p.config, p.events, region).Delete(wait); err != nil {
		return errors.Wrap(err, "deleting stack")
	}

	log.Debug("deleting function")
	if err := p.deleteFunction(region); err != nil && !util.IsNotFound(err) {
		return errors.Wrap(err, "deleting function")
	}

	log.Debug("deleting role")
	if err := p.deleteRole(region); err != nil && !util.IsNotFound(err) {
		return errors.Wrap(err, "deleting role")
	}

	return nil
}

// ShowStack implementation.
func (p *Platform) ShowStack(region string) error {
	return stack.New(p.config, p.events, region).Show()
}

// deploy to the given region.
func (p *Platform) deploy(region, stage string) (version string, err error) {
	fields := event.Fields{
		"stage":  stage,
		"region": region,
	}

	defer p.events.Time("platform.deploy", fields)()

	ctx := log.WithField("region", region)
	s := session.New(aws.NewConfig().WithRegion(region))
	a := apigateway.New(s)
	c := lambda.New(s)

	ctx.Debug("fetching function config")
	_, err = c.GetFunctionConfiguration(&lambda.GetFunctionConfigurationInput{
		FunctionName: &p.config.Name,
	})

	if util.IsNotFound(err) {
		defer p.events.Time("platform.function.create", fields)
		return p.createFunction(c, a, stage)
	}

	if err != nil {
		return "", errors.Wrap(err, "fetching function config")
	}

	defer p.events.Time("platform.function.update", fields)
	return p.updateFunction(c, stage)
}

// createFunction creates the function.
func (p *Platform) createFunction(c *lambda.Lambda, a *apigateway.APIGateway, stage string) (version string, err error) {
retry:
	res, err := c.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: &p.config.Name,
		Handler:      &p.handler,
		Runtime:      &p.runtime,
		Role:         &p.config.Lambda.Role,
		MemorySize:   aws.Int64(int64(p.config.Lambda.Memory)),
		Timeout:      aws.Int64(int64(p.config.Lambda.Timeout)),
		Publish:      aws.Bool(true),
		Environment:  toEnv(p.config.Environment),
		Code: &lambda.FunctionCode{
			ZipFile: p.zip.Bytes(),
		},
	})

	// IAM is eventually consistent apparently, so we have to keep retrying
	if isCreatingRole(err) {
		log.Debug("waiting for role to be created")
		time.Sleep(500 * time.Millisecond)
		goto retry
	}

	if err != nil {
		return "", errors.Wrap(err, "creating function")
	}

	return *res.Version, errFirstDeploy
}

// updateFunction updates the function.
func (p *Platform) updateFunction(c *lambda.Lambda, stage string) (version string, err error) {
	var publish bool

	if stage != "development" {
		publish = true
		log.Debug("publishing new version")
	}

	_, err = c.UpdateFunctionConfiguration(&lambda.UpdateFunctionConfigurationInput{
		FunctionName: &p.config.Name,
		Handler:      &p.handler,
		Runtime:      &p.runtime,
		Role:         &p.config.Lambda.Role,
		MemorySize:   aws.Int64(int64(p.config.Lambda.Memory)),
		Timeout:      aws.Int64(int64(p.config.Lambda.Timeout)),
		Environment:  toEnv(p.config.Environment),
	})

	if err != nil {
		return "", errors.Wrap(err, "updating function config")
	}

	res, err := c.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
		FunctionName: &p.config.Name,
		Publish:      &publish,
		ZipFile:      p.zip.Bytes(),
	})

	if err != nil {
		return "", errors.Wrap(err, "updating function code")
	}

	if publish {
		log.Debugf("alias %s to %s", stage, *res.Version)
		_, err := c.UpdateAlias(&lambda.UpdateAliasInput{
			FunctionName:    &p.config.Name,
			FunctionVersion: res.Version,
			Name:            &stage,
		})

		if err != nil {
			return "", errors.Wrap(err, "creating function alias")
		}

		return *res.Version, nil
	}

	return "", nil
}

// deleteFunction deletes the lambda function.
func (p *Platform) deleteFunction(region string) error {
	// TODO: sessions all over... refactor
	c := lambda.New(session.New(aws.NewConfig().WithRegion(region)))

	_, err := c.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: &p.config.Name,
	})

	return err
}

// createRole creates the IAM role unless it is present.
func (p *Platform) createRole() error {
	s := session.New(aws.NewConfig())
	c := iam.New(s)

	name := fmt.Sprintf("%s-api-function", p.config.Name)
	desc := util.ManagedByUp("")

	if s := p.config.Lambda.Role; s != "" {
		log.Debugf("using role from config %s", s)
		return nil
	}

	log.Debug("checking for role")
	existing, err := c.GetRole(&iam.GetRoleInput{
		RoleName: &name,
	})

	// network or permission error
	if err != nil && !util.IsNotFound(err) {
		return errors.Wrap(err, "fetching role")
	}

	// use the existing role
	if err == nil {
		arn := *existing.Role.Arn
		log.Debugf("using existing role %s", arn)
		p.config.Lambda.Role = arn
		return nil
	}

	log.Debug("creating role")
	role, err := c.CreateRole(&iam.CreateRoleInput{
		RoleName:                 &name,
		Description:              &desc,
		AssumeRolePolicyDocument: &apiGatewayAssumePolicy,
	})

	if err != nil {
		return errors.Wrap(err, "creating role")
	}

	log.Debug("attaching policy")
	_, err = c.PutRolePolicy(&iam.PutRolePolicyInput{
		PolicyName:     &name,
		RoleName:       &name,
		PolicyDocument: &logWritePolicy,
	})

	if err != nil {
		return errors.Wrap(err, "attaching policy")
	}

	arn := *role.Role.Arn
	log.Debugf("set role to %s", arn)
	p.config.Lambda.Role = arn

	return nil
}

// deleteRole deletes the role and policy.
func (p *Platform) deleteRole(region string) error {
	name := fmt.Sprintf("%s-api-function", p.config.Name)
	c := iam.New(session.New(aws.NewConfig().WithRegion(region)))

	_, err := c.DeleteRolePolicy(&iam.DeleteRolePolicyInput{
		RoleName:   &name,
		PolicyName: &name,
	})

	if err != nil {
		return errors.Wrap(err, "deleting role policy")
	}

	_, err = c.DeleteRole(&iam.DeleteRoleInput{
		RoleName: &name,
	})

	if err != nil {
		return errors.Wrap(err, "deleting role")
	}

	return nil
}

// getAPI returns the API if present.
func (p *Platform) getAPI(c *apigateway.APIGateway) (api *apigateway.RestApi, err error) {
	name := p.config.Name

	res, err := c.GetRestApis(&apigateway.GetRestApisInput{
		Limit: aws.Int64(500),
	})

	if err != nil {
		return nil, errors.Wrap(err, "fetching apis")
	}

	for _, a := range res.Items {
		if *a.Name == name {
			api = a
		}
	}

	return
}

// injectProxy injects the Go proxy.
func (p *Platform) injectProxy() error {
	log.Debugf("injecting proxy")

	if err := ioutil.WriteFile("main", bin.MustAsset("up-proxy"), 0777); err != nil {
		return errors.Wrap(err, "writing up-proxy")
	}

	if err := ioutil.WriteFile("byline.js", shim.MustAsset("byline.js"), 0755); err != nil {
		return errors.Wrap(err, "writing byline.js")
	}

	if err := ioutil.WriteFile("_proxy.js", shim.MustAsset("index.js"), 0755); err != nil {
		return errors.Wrap(err, "writing _proxy.js")
	}

	return nil
}

// removeProxy removes the Go proxy.
func (p *Platform) removeProxy() error {
	log.Debugf("removing proxy")
	os.Remove("main")
	os.Remove("_proxy.js")
	os.Remove("byline.js")
	return nil
}

// isCreatingRole returns true if the role has not been created.
func isCreatingRole(err error) bool {
	return err != nil && strings.Contains(err.Error(), "role defined for the function cannot be assumed by Lambda")
}

// toEnv returns a lambda environment.
func toEnv(env config.Environment) *lambda.Environment {
	return &lambda.Environment{
		Variables: aws.StringMap(env),
	}
}

// hash returns a base64 encoded SHA256 hash of `b`.
func hash(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
