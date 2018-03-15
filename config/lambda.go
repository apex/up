package config

import (
	"errors"
	"time"
)

// Lambda configuration.
type Lambda struct {
	// Memory of the function.
	Memory int `json:"memory"`

	// Timeout of the function.
	Timeout int `json:"timeout"`

	// Role of the function.
	Role string `json:"role"`

	// Runtime of the function.
	Runtime string `json:"runtime"`

	// Accelerate enables S3 acceleration.
	Accelerate bool `json:"accelerate"`

	// Warm enables active warming.
	Warm *bool `json:"warm"`

	// WarmCount is the number of containers to keep active.
	WarmCount int `json:"warm_count"`

	// WarmRate is the schedule for performing worming.
	WarmRate Duration `json:"warm_rate"`
}

// Default implementation.
func (l *Lambda) Default() error {
	if l.Memory == 0 {
		l.Memory = 512
	}

	if l.Runtime == "" {
		l.Runtime = "nodejs8.10"
	}

	if l.WarmRate == 0 {
		l.WarmRate = Duration(15 * time.Minute)
	}

	if l.WarmCount == 0 {
		l.WarmCount = 15
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

// Override config.
func (l *Lambda) Override(c *Config) {
	if l.Memory != 0 {
		c.Lambda.Memory = l.Memory
	}

	if l.Timeout != 0 {
		c.Lambda.Timeout = l.Timeout
	}

	if l.Role != "" {
		c.Lambda.Role = l.Role
	}

	if l.Warm != nil {
		c.Lambda.Warm = l.Warm
	}

	if l.WarmCount > 0 {
		c.Lambda.WarmCount = l.WarmCount
	}

	if l.WarmRate != 0 {
		c.Lambda.WarmRate = l.WarmRate
	}
}
