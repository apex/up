package platform

import (
	"io"
	"time"
)

// TODO: these interfaces suck, don't mind them for now :D

// Logs is the interface for viewing logs.
type Logs interface {
	Follow()
	Since(time.Time)
	io.Reader
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
	Logs(region, query string) Logs

	// URL returns the endpoitn for the given
	// region and stage combination, or an
	// empty string.
	URL(region, stage string) (string, error)

	// TODO: finalize and document
	CreateStack(region, version string) error
	DeleteStack(region string, wait bool) error
	ShowStack(region string) error
	PlanStack(region string) error
	ApplyStack(region string) error

	ShowMetrics(region, stage string, start time.Time) error
}

// Zipper is the interface used by platforms which
// utilize zips for delivery of deployments.
type Zipper interface {
	Zip() io.Reader
}
