// Package setup provides up.json initialization.
package setup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/mitchellh/go-homedir"
	"github.com/tj/go/term"
	"github.com/tj/survey"

	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/validate"
)

// Errors.
var (
	// ErrNoCredentials is the error returned when no AWS credential profiles are available.
	ErrNoCredentials = errors.New("no credentials")
)

// questions for the user.
var questions = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Name of the project:",
			Default: defaultName(),
		},
		Validate: validateName,
	},
	{
		Name: "profile",
		Prompt: &survey.Select{
			Message:  "AWS credentials profile:",
			Options:  awsProfiles(),
			Default:  os.Getenv("AWS_PROFILE"),
			PageSize: 10,
		},
		Validate: survey.Required,
	},
}

// Create an up.json file for the user.
func Create() error {
	var in struct {
		Name    string `json:"name"`
		Profile string `json:"profile"`
	}

	if len(awsProfiles()) == 0 {
		return ErrNoCredentials
	}

	println()

	// confirm
	var ok bool
	err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("This directory has no up.json, create it?"),
	}, &ok, nil)

	if err != nil {
		return err
	}

	if !ok {
		return errors.New("aborted")
	}

	// prompt
	term.MoveUp(1)
	term.ClearLine()
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

// awsProfiles returns the AWS profiles found.
func awsProfiles() []string {
	path, err := homedir.Expand("~/.aws/credentials")
	if err != nil {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	s, err := util.ParseSections(f)
	if err != nil {
		return nil
	}

	sort.Strings(s)
	return s
}
