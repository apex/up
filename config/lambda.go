package config

// Lambda configuration.
type Lambda struct {
	Role string `json:"role"`

	// Memory of the function.
	Memory int `json:"memory"`

	// Timeout of the function.
	Timeout int `json:"timeout"`
}

// Default implementation.
func (l *Lambda) Default() error {
	if l.Memory == 0 {
		l.Memory = 512
	}

	if l.Timeout == 0 {
		l.Timeout = 15
	}

	return nil
}
