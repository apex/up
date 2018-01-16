package lambda

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/tj/assert"
)

func TestGetCert(t *testing.T) {
	certs := []*acm.CertificateSummary{
		{
			DomainName:     aws.String("example.com"),
			CertificateArn: aws.String("arn:example.com"),
		},
		{
			DomainName:     aws.String("*.example.com"),
			CertificateArn: aws.String("arn:*.example.com"),
		},
	}

	certArn := getCert(certs, "example.com")
	assert.Equal(t, "arn:example.com", certArn)

	certArn = getCert(certs, "sub.example.com")
	assert.Equal(t, "arn:*.example.com", certArn)

	certArn = getCert(certs, "www.sub.example.com")
	assert.Equal(t, "", certArn)
}
