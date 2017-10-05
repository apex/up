// Package validate provides config validation functions.
package validate

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// RequiredString validation.
func RequiredString(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("is required")
	}

	return nil
}

// RequiredStrings validation.
func RequiredStrings(s []string) error {
	for i, v := range s {
		if err := RequiredString(v); err != nil {
			return errors.Wrapf(err, "at index %d", i)
		}
	}

	return nil
}

// MinStrings validation.
func MinStrings(s []string, n int) error {
	if len(s) < n {
		if n == 1 {
			return errors.Errorf("must have at least %d value", n)
		}

		return errors.Errorf("must have at least %d values", n)
	}

	return nil
}

// name regexp.
var name = regexp.MustCompile(`^[a-z][-a-z0-9]*$`)

// Name validation.
func Name(s string) error {
	if !name.MatchString(s) {
		return errors.Errorf("must contain only lowercase alphanumeric characters and '-'")
	}

	return nil
}

// stage regexp.
var stage = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// Stage name validation.
func Stage(s string) error {
	if !stage.MatchString(s) {
		return errors.Errorf("must contain only alphanumeric characters and '_'")
	}

	return nil
}

// List validation.
func List(s string, list []string) error {
	for _, v := range list {
		if s == v {
			return nil
		}
	}

	return errors.Errorf("%q is invalid, must be one of:\n\n  • %s", s, strings.Join(list, "\n  • "))
}

// Lists validation.
func Lists(vals, list []string) error {
	for _, v := range vals {
		if err := List(v, list); err != nil {
			return err
		}
	}

	return nil
}

// OptionalStage validation.
func OptionalStage(stage string) error {
	if stage == "" {
		return nil
	}

	return Stage(stage)
}
