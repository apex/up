package config

import (
	"encoding/json"
	"os"
)

// Environment variables.
type Environment map[string]string

// UnmarshalJSON implementation.
func (e *Environment) UnmarshalJSON(b []byte) error {
	valueFromConfig := map[string]string{}

	if err := json.Unmarshal(b, &valueFromConfig); err != nil {
		return err
	}

	for k, v := range valueFromConfig {
		valueFromConfig[k] = os.ExpandEnv(v)
	}

	*e = valueFromConfig

	return nil
}
