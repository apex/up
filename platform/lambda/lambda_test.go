package lambda

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/tj/assert"
	
	"github.com/apex/up/config"
	"github.com/apex/up/platform/event"
)

func TestGetCert(t *testing.T) {
	certs := []*acm.CertificateDetail{
		{
			DomainName:     aws.String("example.com"),
			CertificateArn: aws.String("arn:example.com"),
			SubjectAlternativeNames: aws.StringSlice([]string{
				"*.example.com",
			}),
		},
		{
			DomainName:     aws.String("*.apex.sh"),
			CertificateArn: aws.String("arn:*.apex.sh"),
		},
		{
			DomainName:     aws.String("api.example.com"),
			CertificateArn: aws.String("arn:api.example.com"),
			SubjectAlternativeNames: aws.StringSlice([]string{
				"*.api.example.com",
				"something.example.com",
			}),
		},
	}

	arn := getCert(certs, "example.com")
	assert.Equal(t, "arn:example.com", arn)

	arn = getCert(certs, "www.example.com")
	assert.Equal(t, "arn:example.com", arn)

	arn = getCert(certs, "api.example.com")
	assert.Equal(t, "arn:api.example.com", arn)

	arn = getCert(certs, "apex.sh")
	assert.Empty(t, arn)

	arn = getCert(certs, "api.apex.sh")
	assert.Equal(t, "arn:*.apex.sh", arn)

	arn = getCert(certs, "v1.api.example.com")
	assert.Equal(t, "arn:api.example.com", arn)

	arn = getCert(certs, "something.example.com")
	assert.Equal(t, "arn:api.example.com", arn)

	arn = getCert(certs, "staging.v1.api.example.com")
	assert.Empty(t, arn)
}

func TestCreateRole(t *testing.T) {
	t.Run("doesn't attempt to create configured role", func(t *testing.T) {
		c := &config.Config{
			Lambda: config.Lambda{
				Role: "custom-role-name",
			},
		}
		events := make(event.Events)
		p := New(c, events)
		assert.NoError(t, p.createRole(), "createRole")
	})
}

func TestDeleteRole(t *testing.T) {
	t.Run("doesn't attempt to delete configured role", func(t *testing.T) {
		c := &config.Config{
			Lambda: config.Lambda{
				Role: "custom-role-name",
			},
		}
		events := make(event.Events)
		p := New(c, events)
		assert.NoError(t, p.deleteRole("us-west-2"), "deleteRole")
	})
}
