package config

import (
	"encoding/json"
	"errors"
)

// Hook is one or more commands.
type Hook []string

// Hooks for the project.
type Hooks struct {
	Build      Hook `json:"build"`
	Clean      Hook `json:"clean"`
	PreBuild   Hook `json:"prebuild"`
	PostBuild  Hook `json:"postbuild"`
	PreDeploy  Hook `json:"predeploy"`
	PostDeploy Hook `json:"postdeploy"`
}

// Override config.
func (h *Hooks) Override(c *Config) {
	if v := h.Build; v != nil {
		c.Hooks.Build = v
	}

	if v := h.Clean; v != nil {
		c.Hooks.Clean = v
	}

	if v := h.PreBuild; v != nil {
		c.Hooks.PreBuild = v
	}

	if v := h.PostBuild; v != nil {
		c.Hooks.PostBuild = v
	}

	if v := h.PreDeploy; v != nil {
		c.Hooks.PreDeploy = v
	}

	if v := h.PostDeploy; v != nil {
		c.Hooks.PostDeploy = v
	}
}

// Get returns the hook by name or nil.
func (h *Hooks) Get(s string) Hook {
	switch s {
	case "build":
		return h.Build
	case "clean":
		return h.Clean
	case "prebuild":
		return h.PreBuild
	case "postbuild":
		return h.PostBuild
	case "predeploy":
		return h.PreDeploy
	case "postdeploy":
		return h.PostDeploy
	default:
		return nil
	}
}

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
