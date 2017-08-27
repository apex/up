package config

import (
	"encoding/json"
	"errors"
)

// Hooks for the project.
type Hooks struct {
	Build Hook `json:"build"`
	Clean Hook `json:"clean"`
}

// Hook is one or more commands.
type Hook []string

// UnmarshalJSON implementation.
func (h *Hook) UnmarshalJSON(b []byte) error {
	switch b[0] {
	case '"':
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*h = append(*h, s)
		return nil
	case '[':
		return json.Unmarshal(b, (*[]string)(h))
	default:
		return errors.New("hook must be a string or array of strings")
	}
}

// IsEmpty returns true if the hook is empty.
func (h *Hook) IsEmpty() bool {
	return h == nil || len(*h) == 0
}
