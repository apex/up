package up

import (
	"io"
	"time"
)

// TODO: finalize and finish documentation

// LogsConfig is configuration for viewing logs.
type LogsConfig struct {
	// Region is the target region.
	Region string

	// Query is the filter pattern.
	Query string

	// Since is used as the starting point when filtering
	// historical logs, no logs before this point are returned.
	Since time.Time

	// Follow is used to stream new logs.
	Follow bool

	// Expand is used to expand logs to a verbose format.
	Expand bool

	// OutputJSON is used to output raw json.
	OutputJSON bool
}

// Logs is the interface for viewing platform logs.
type Logs interface {
	io.Reader
}

// Domains is the interface for purchasing and
// managing domains names.
type Domains interface {
	Availability(domain string) (*Domain, error)
	Suggestions(domain string) ([]*Domain, error)
	Purchase(domain string, contact DomainContact) error
	List() ([]*Domain, error)
}

// Secret is an encrypted variable..
type Secret struct {
	App              string
	Name             string
	Stage            string
	Value            string
	Description      string
	LastModifiedUser string
	LastModified     time.Time
}

// Secrets is the interface for managing encrypted secrets.
type Secrets interface {
	Add(key, val, desc string, clear bool) error
	Remove(key string) error
	List() ([]*Secret, error)
	Load() ([]*Secret, error)
}

// Interface for platforms.
type Interface interface {
	// Build the project.
	Build() error

	// Deploy to the given stage, to the
	// region(s) configured by the user.
	Deploy(stage string) error

	// Logs returns an interface for working
	// with logging data.
	Logs(LogsConfig) Logs

	// Domains returns an interface for
	// managing domain names.
	Domains() Domains

	// Secrets returns an interface for
	// managing secret variables.
	Secrets(stage string) Secrets

	// URL returns the endpoitn for the given
	// region and stage combination, or an
	// empty string.
	URL(region, stage string) (string, error)

	CreateStack(region, version string) error
	DeleteStack(region string, wait bool) error
	ShowStack(region string) error
	PlanStack(region string) error
	ApplyStack(region string) error

	ShowMetrics(region, stage string, start time.Time) error
}

// Runtime is the interface used by a platform to support
// runtime operations such as initializing environment
// variables from remote storage.
type Runtime interface {
	Init(stage string) error
}

// Zipper is the interface used by platforms which
// utilize zips for delivery of deployments.
type Zipper interface {
	Zip() io.Reader
}

// Domain is a domain name and its availability.
type Domain struct {
	Name      string
	Available bool
	Expiry    time.Time
	AutoRenew bool
}

// DomainContact is the domain name contact
// information required for registration.
type DomainContact struct {
	Email            string
	FirstName        string
	LastName         string
	CountryCode      string
	City             string
	Address          string
	OrganizationName string
	PhoneNumber      string
	State            string
	ZipCode          string
}
