package account

import (
	"github.com/stripe/stripe-go"
	"github.com/tj/survey"
)

// Questions.
var questions = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{Message: "Name:"},
		Validate: survey.Required,
	},
	{
		Name:     "number",
		Prompt:   &survey.Input{Message: "Number:"},
		Validate: survey.Required,
	},
	{
		Name:     "cvc",
		Prompt:   &survey.Input{Message: "CVC:"},
		Validate: survey.Required,
	},
	{
		Name:     "month",
		Prompt:   &survey.Input{Message: "Expiration month:"},
		Validate: survey.Required,
	},
	{
		Name:     "year",
		Prompt:   &survey.Input{Message: "Expiration year:"},
		Validate: survey.Required,
	},
	{
		Name:     "address1",
		Prompt:   &survey.Input{Message: "Street Address:"},
		Validate: survey.Required,
	},
	{
		Name:     "city",
		Prompt:   &survey.Input{Message: "City:"},
		Validate: survey.Required,
	},
	{
		Name:     "state",
		Prompt:   &survey.Input{Message: "State:"},
		Validate: survey.Required,
	},
	{
		Name:     "country",
		Prompt:   &survey.Input{Message: "Country:"},
		Validate: survey.Required,
	},
	{
		Name:     "zip",
		Prompt:   &survey.Input{Message: "Zip:"},
		Validate: survey.Required,
	},
}

// PromptForCard displays an interactive form for the user to provide CC details.
func PromptForCard() (card stripe.CardParams, err error) {
	err = survey.Ask(questions, &card)
	return
}
