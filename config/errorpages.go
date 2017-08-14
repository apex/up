package config

// ErrorPages configuration.
type ErrorPages struct {
	// Enable default error pages.
	Enable bool `json:"enable"`

	// Dir containing error pages.
	Dir string `json:"dir"`

	// Variables are passed to the template for use.
	Variables map[string]interface{} `json:"variables"`
}

// Default implementation.
func (e *ErrorPages) Default() error {
	if e.Dir == "" {
		e.Dir = "."
	}

	return nil
}
