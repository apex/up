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
	"github.com/apex/up/platform/aws/regions"
)

// ErrNoCredentials is the error returned when no AWS credential profiles are available.
var ErrNoCredentials = errors.New("no credentials")

// config saved to up.json
type config struct {
	Name    string   `json:"name"`
	Profile string   `json:"profile"`
	Regions []string `json:"regions"`
}

// questions for the user.
var questions = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Project name:",
			Default: defaultName(),
		},
		Validate: validateName,
	},
	{
		Name: "profile",
		Prompt: &survey.Select{
			Message:  "AWS profile:",
			Options:  awsProfiles(),
			Default:  os.Getenv("AWS_PROFILE"),
			PageSize: 10,
		},
		Validate: survey.Required,
	},
	{
		Name: "region",
		Prompt: &survey.Select{
			Message:  "AWS region:",
			Options:  regions.Names,
			Default:  defaultRegion(),
			PageSize: 15,
		},
		Validate: survey.Required,
	},
}

// Create an up.json file for the user.
func Create() error {
	var in struct {
		Name    string `json:"name"`
		Profile string `json:"profile"`
		Region  string `json:"region"`
	}

	if len(awsProfiles()) == 0 {
		return ErrNoCredentials
	}

	println()

	// confirm create new project
	var ok bool
	err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("No up.json found, create a new project?"),
		Default: true,
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

	c := config{
		Name:    in.Name,
		Profile: in.Profile,
		Regions: []string{
			regions.GetIdByName(in.Region),
		},
	}

	b, _ := json.MarshalIndent(c, "", "  ")
	err = ioutil.WriteFile("up.json", b, 0644)
	if err != nil {
		return err
	}

	// confirm create .upignore
	err = survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Would you like to add an .upignore?"),
		Default: true,
	}, &ok, nil)

	if err != nil {
		return nil
	}

	if !ok {
		return errors.New("aborted")
	}

	defaultIgnore := ".*\n"
	b = []byte(defaultIgnore)
	return ioutil.WriteFile(".upignore", b, 0644)
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

// defaultRegion returns the default aws region.
func defaultRegion() string {
	if s := os.Getenv("AWS_DEFAULT_REGION"); s != "" {
		return s
	}

	if s := os.Getenv("AWS_REGION"); s != "" {
		return s
	}

	return ""
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
