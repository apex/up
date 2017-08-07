// Package logs provides logging utilities.
package logs

import (
	"os"

	"github.com/apex/log"
)

// Plugin returns a log context for the given plugin name.
func Plugin(name string) log.Interface {
	return log.WithFields(log.Fields{
		"app_name":    os.Getenv("AWS_LAMBDA_FUNCTION_NAME"),
		"app_region":  os.Getenv("AWS_REGION"),
		"app_version": os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"),
		"plugin":      name,
	})
}
