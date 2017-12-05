// Package logs provides logging utilities.
package logs

import (
	"os"

	"github.com/apex/log"
)

// Plugin returns a log context for the given plugin name.
func Plugin(name string) log.Interface {
	return log.WithFields(log.Fields{
		"app":     os.Getenv("AWS_LAMBDA_FUNCTION_NAME"),
		"region":  os.Getenv("AWS_REGION"),
		"version": os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"),
		"stage":   os.Getenv("UP_STAGE"),
		"plugin":  name,
	})
}
