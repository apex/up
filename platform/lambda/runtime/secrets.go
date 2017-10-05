package runtime

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/apex/up/internal/secret"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/platform"
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
func (s *Secrets) List() (v []*platform.Secret, err error) {
	res, err := s.client.DescribeParameters(&ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
		Filters: []*ssm.ParametersFilter{
			{
				Key:    aws.String("Name"),
				Values: aws.StringSlice([]string{"/up/"}),
			},
		},
	})

	if err != nil {
		return
	}

	for _, p := range res.Parameters {
		user := strings.Split(*p.LastModifiedUser, "/")[1]
		app, stage, name := secret.Parse(*p.Name)
		v = append(v, &platform.Secret{
			App:              app,
			Name:             name,
			Stage:            stage,
			Description:      util.DefaultString(p.Description, ""),
			LastModifiedUser: user,
			LastModified:     *p.LastModifiedDate,
		})
	}

	return
}

// Load implementation.
func (s *Secrets) Load() (v []*platform.Secret, err error) {
	res, err := s.client.GetParametersByPath(&ssm.GetParametersByPathInput{
		MaxResults:     aws.Int64(10),
		Path:           aws.String("/up/"),
		WithDecryption: aws.Bool(true),
		Recursive:      aws.Bool(true),
	})

	if err != nil {
		return
	}

	for _, p := range res.Parameters {
		app, stage, name := secret.Parse(*p.Name)
		v = append(v, &platform.Secret{
			App:   app,
			Name:  name,
			Stage: stage,
			Value: *p.Value,
		})
	}

	return
}

// Add implementation.
func (s *Secrets) Add(key, val, desc string) error {
	key = s.secretName(key)

	params := &ssm.PutParameterInput{
		Type:      aws.String("SecureString"),
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
