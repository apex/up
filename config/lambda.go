package config

import "errors"

// Lambda configuration.
type Lambda struct {
	// Memory of the function.
	Memory int `json:"memory"`

	// Timeout of the function.
	Timeout int `json:"timeout"`

	// Role of the function.
	Role string `json:"role"`
}

// Default implementation.
func (l *Lambda) Default() error {
	if l.Memory == 0 {
		l.Memory = 512
	}

	return nil
}

// Validate implementation.
func (l *Lambda) Validate() error {
	if l.Timeout != 0 {
		return errors.New(".lambda.timeout is deprecated, use .proxy.timeout")
	}

	return nil
}
