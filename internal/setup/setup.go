// Package setup provides up.json initialization.
package setup

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tj/survey"
)

// questions for the user.
var questions = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Name:",
			Default: defaultName(),
		},
		Validate: survey.Required,
	},
}

// Create an up.json file for the user.
func Create() error {
	var in struct {
		Name string `json:"name"`
	}

	if err := survey.Ask(questions, &in); err != nil {
		return err
	}

	b, _ := json.MarshalIndent(in, "", "  ")
	return ioutil.WriteFile("up.json", b, 0755)
}

// defaultName returns the default app name.
func defaultName() string {
	dir, _ := os.Getwd()
	return filepath.Base(dir)
}
