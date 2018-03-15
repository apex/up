// Package logs provides logging utilities.
package logs

import (
	"os"

	"github.com/apex/log"
)

// Fields returns the global log fields.
func Fields() log.Fields {
	f := log.Fields{
		"app":     os.Getenv("AWS_LAMBDA_FUNCTION_NAME"),
		"region":  os.Getenv("AWS_REGION"),
		"version": os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"),
		"stage":   os.Getenv("UP_STAGE"),
	}

	if s := os.Getenv("UP_COMMIT"); s != "" {
		f["commit"] = s
	}

	return f
}

// Plugin returns a log context for the given plugin name.
func Plugin(name string) log.Interface {
	f := Fields()
	f["plugin"] = name
	return log.WithFields(f)
}
