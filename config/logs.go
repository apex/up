package config

// Logs configuration.
type Logs struct {
	// Disable json log output.
	Disable bool `json:"disable"`

	// Stdout default log level.
	Stdout string `json:"stdout"`

	// Stderr default log level.
	Stderr string `json:"stderr"`
}

// Default implementation.
func (l *Logs) Default() error {
	if l.Stdout == "" {
		l.Stdout = "info"
	}

	if l.Stderr == "" {
		l.Stderr = "error"
	}

	return nil
}
