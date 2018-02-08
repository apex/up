package runtime

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"

	"github.com/apex/up"
	"github.com/apex/up/internal/secret"
	"github.com/apex/up/internal/util"
)

// TODO: secret pagination

// Secrets implementation.
type Secrets struct {
	client *ssm.SSM
	name   string
	stage  string
}

// NewSecrets returns a new secrets manager.
func NewSecrets(name, stage, region string) *Secrets {
	return &Secrets{
		client: ssm.New(session.New(aws.NewConfig().WithRegion(region))),
		name:   name,
		stage:  stage,
	}
}

// List implementation.
func (s *Secrets) List() (v []*up.Secret, err error) {
	res, err := s.client.DescribeParameters(&ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
		Filters: []*ssm.ParametersFilter{
			{
				Key:    aws.String("Name"),
				Values: aws.StringSlice([]string{"/up/" + s.name + "/"}),
			},
		},
	})

	if err != nil {
		return
	}

	for _, p := range res.Parameters {
		var value string

		if *p.Type == "String" {
			p, err := s.client.GetParameter(&ssm.GetParameterInput{
				Name: p.Name,
			})

			if err != nil {
				return nil, errors.Wrap(err, "getting parameter")
			}

			value = *p.Parameter.Value
		}

		app, stage, name := secret.Parse(*p.Name)
		v = append(v, &up.Secret{
			App:              app,
			Name:             name,
			Stage:            stage,
			Type:             *p.Type,
			Description:      util.DefaultString(p.Description, ""),
			LastModifiedUser: userFromARN(p.LastModifiedUser),
			LastModified:     *p.LastModifiedDate,
			Value:            value,
		})
	}

	return
}

// Load implementation.
func (s *Secrets) Load() (v []*up.Secret, err error) {
	var token *string

	for {
		res, err := s.client.GetParametersByPath(&ssm.GetParametersByPathInput{
			MaxResults:     aws.Int64(10),
			Path:           aws.String("/up/" + s.name + "/"),
			WithDecryption: aws.Bool(true),
			Recursive:      aws.Bool(true),
			NextToken:      token,
		})

		if err != nil {
			return nil, err
		}

		for _, p := range res.Parameters {
			app, stage, name := secret.Parse(*p.Name)
			v = append(v, &up.Secret{
				App:   app,
				Name:  name,
				Stage: stage,
				Type:  *p.Type,
				Value: *p.Value,
			})
		}

		token = res.NextToken

		if token == nil {
			break
		}
	}

	return
}

// Add implementation.
func (s *Secrets) Add(key, val, desc string, clear bool) error {
	key = s.secretName(key)

	kind := "SecureString"
	if clear {
		kind = "String"
	}

	params := &ssm.PutParameterInput{
		Type:      &kind,
		Name:      &key,
		Value:     &val,
		Overwrite: aws.Bool(true),
	}

	if desc != "" {
		params.Description = &desc
	}

	_, err := s.client.PutParameter(params)
	return err
}

// Remove implementation.
func (s *Secrets) Remove(key string) error {
	key = s.secretName(key)

	_, err := s.client.DeleteParameter(&ssm.DeleteParameterInput{
		Name: &key,
	})

	return err
}

// secretName returns the secret name normalized.
func (s *Secrets) secretName(name string) string {
	return secret.Format(s.name, s.stage, name)
}

// userFromARN returns the username from ARN if present.
func userFromARN(arn *string) string {
	if arn == nil {
		return ""
	}

	p := strings.Split(*arn, "/")

	if len(p) >= 2 {
		return p[1]
	}

	return ""
}
