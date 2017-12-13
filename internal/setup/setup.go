// Package setup provides up.json initialization.
package setup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/up/internal/validate"
	"github.com/tj/survey"
)

// questions for the user.
var questions = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Name of your project:",
			Default: defaultName(),
		},
		Validate: validateName,
	},
}

// Create an up.json file for the user.
func Create() error {
	var in struct {
		Name string `json:"name"`
	}

	// confirm
	var ok bool
	err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("This directory has no up.json, create it?"),
	}, &ok, nil)

	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	if err := survey.Ask(questions, &in); err != nil {
		return err
	}

	b, _ := json.MarshalIndent(in, "", "  ")
	return ioutil.WriteFile("up.json", b, 0755)
}

// defaultName returns the default app name.
// The name is only inferred if it is valid.
func defaultName() string {
	dir, _ := os.Getwd()
	name := filepath.Base(dir)
	if validate.Name(name) != nil {
		return ""
	}
	return name
}

// validateName validates the name prompt.
func validateName(v interface{}) error {
	if err := validate.Name(v.(string)); err != nil {
		return err
	}

	return survey.Required(v)
}
