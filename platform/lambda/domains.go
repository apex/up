package lambda

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	r "github.com/aws/aws-sdk-go/service/route53domains"

	"github.com/apex/up"
	"github.com/apex/up/platform"
)

// Domains implementation.
type Domains struct {
	client *r.Route53Domains
}

// NewDomains returns a new domain manager.
func NewDomains(c *up.Config) *Domains {
	return &Domains{
		client: r.New(session.New(aws.NewConfig().WithRegion("us-east-1").WithEndpoint(c.GetEndpoint(up.Route53)))),
	}
}

// List implementation.
func (d *Domains) List() (v []*platform.Domain, err error) {
	res, err := d.client.ListDomains(&r.ListDomainsInput{
		MaxItems: aws.Int64(100),
	})

	if err != nil {
		return
	}

	for _, d := range res.Domains {
		v = append(v, &platform.Domain{
			Name:      *d.DomainName,
			Expiry:    *d.Expiry,
			AutoRenew: *d.AutoRenew,
		})
	}

	return
}

// Availability implementation.
func (d *Domains) Availability(domain string) (*platform.Domain, error) {
	res, err := d.client.CheckDomainAvailability(&r.CheckDomainAvailabilityInput{
		DomainName: &domain,
	})

	if err != nil {
		return nil, err
	}

	if *res.Availability == "AVAILABLE" {
		return &platform.Domain{
			Name:      domain,
			Available: true,
		}, nil
	}

	return &platform.Domain{
		Name:      domain,
		Available: false,
	}, nil
}

// Suggestions implementation.
func (d *Domains) Suggestions(domain string) (domains []*platform.Domain, err error) {
	res, err := d.client.GetDomainSuggestions(&r.GetDomainSuggestionsInput{
		DomainName:      &domain,
		OnlyAvailable:   aws.Bool(true),
		SuggestionCount: aws.Int64(15),
	})

	if err != nil {
		return
	}

	for _, s := range res.SuggestionsList {
		domains = append(domains, &platform.Domain{
			Name:      *s.DomainName,
			Available: true,
		})
	}

	return
}

// Purchase implementation.
func (d *Domains) Purchase(domain string, contact platform.DomainContact) error {
	_, err := d.client.RegisterDomain(&r.RegisterDomainInput{
		DomainName:        &domain,
		AutoRenew:         aws.Bool(true),
		DurationInYears:   aws.Int64(1),
		RegistrantContact: contactDetails(contact),
		AdminContact:      contactDetails(contact),
		TechContact:       contactDetails(contact),
	})

	return err
}

// contactDetails returns route53 contact details.
func contactDetails(c platform.DomainContact) *r.ContactDetail {
	return &r.ContactDetail{
		AddressLine1: aws.String(c.Address),
		City:         aws.String(c.City),
		State:        aws.String(c.State),
		ZipCode:      aws.String(c.ZipCode),
		CountryCode:  aws.String(c.CountryCode),
		Email:        aws.String(c.Email),
		PhoneNumber:  aws.String(c.PhoneNumber),
		FirstName:    aws.String(c.FirstName),
		LastName:     aws.String(c.LastName),
		ContactType:  aws.String("PERSON"),
	}
}
