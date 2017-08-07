package platform

import (
	"io"
	"time"
)

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
	Logs(query string) Logs

	// URL returns the endpoitn for the given
	// region and stage combination, or an
	// empty string.
	URL(region, stage string) (string, error)

	// TODO: finalize and document
	CreateStack(region, version string) error
	DeleteStack(region string, wait bool) error
	ShowStack(region string) error
}
