package config

import (
	"strings"
	"time"

	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/validate"
	"github.com/pkg/errors"
)

// namespace mappings.
var namespaceMap = map[string]string{
	"http": "AWS/ApiGateway",
}

// metrics mappings.
var metricsMap = map[string]string{
	"http.count":   "Count",
	"http.latency": "Latency",
	"http.4xx":     "4XXError",
	"http.5xx":     "5XXError",
}

// operator mappings.
var operatorMap = map[string]string{
	">":  "GreaterThanThreshold",
	"<":  "LessThanThreshold",
	">=": "GreaterThanOrEqualToThreshold",
	"<=": "LessThanOrEqualToThreshold",
}

// statistic mappings.
var statisticMap = map[string]string{
	"count":   "SampleCount",
	"sum":     "Sum",
	"avg":     "Average",
	"average": "Average",
	"min":     "Minimum",
	"minimum": "Minimum",
	"max":     "Maximum",
	"maximum": "Maximum",
}

// missingData options.
var missingData = []string{
	"breaching",
	"notBreaching",
	"ignore",
	"missing",
}

// AlertAction config.
type AlertAction struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Emails []string `json:"emails"`
}

// Validate implementation.
func (a *AlertAction) Validate() error {
	if err := validate.RequiredString(a.Name); err != nil {
		return errors.Wrap(err, ".name")
	}

	if err := validate.List(a.Type, []string{"email"}); err != nil {
		return errors.Wrap(err, ".type")
	}

	if a.Type == "email" {
		if err := validate.MinStrings(a.Emails, 1); err != nil {
			return errors.Wrap(err, ".emails")
		}
	}

	return nil
}

// Alert config.
type Alert struct {
	Description       string   `json:"description"`
	Disable           bool     `json:"disable"`
	Metric            string   `json:"metric"`
	Namespace         string   `json:"namespace"`
	Statistic         string   `json:"statistic"`
	Operator          string   `json:"operator"`
	Threshold         int      `json:"threshold"`
	Period            Duration `json:"period"` // TODO: must be multiple of 60?
	EvaluationPeriods int      `json:"evaluation_periods"`
	Stage             string   `json:"stage"`
	Action            string   `json:"action"`
	Missing           string   `json:"missing"`
}

// Default implementation.
func (a *Alert) Default() error {
	if a.Operator == "" {
		a.Operator = ">"
	}

	if a.Missing == "" {
		a.Missing = "missing"
	}

	if a.Period == 0 {
		a.Period = Duration(time.Minute)
	}

	if a.EvaluationPeriods == 0 {
		a.EvaluationPeriods = 1
	}

	if s := a.Metric; s != "" {
		parts := strings.Split(a.Metric, ".")

		if s, ok := namespaceMap[parts[0]]; ok {
			a.Namespace = s
		}

		if s, ok := metricsMap[a.Metric]; ok {
			a.Metric = s
		}
	}

	return nil
}

// Validate implementation.
func (a *Alert) Validate() error {
	if s, ok := operatorMap[a.Operator]; ok {
		a.Operator = s
	} else {
		if err := validate.List(a.Operator, util.StringMapKeys(operatorMap)); err != nil {
			return errors.Wrap(err, ".operator")
		}
	}

	if s, ok := statisticMap[a.Statistic]; ok {
		a.Statistic = s
	} else {
		if err := validate.List(a.Statistic, util.StringMapKeys(statisticMap)); err != nil {
			return errors.Wrap(err, ".statistic")
		}
	}

	if err := validate.List(a.Missing, missingData); err != nil {
		return errors.Wrap(err, ".missing")
	}

	if err := validate.RequiredString(a.Metric); err != nil {
		return errors.Wrap(err, ".metric")
	}

	if err := validate.RequiredString(a.Namespace); err != nil {
		return errors.Wrap(err, ".namespace")
	}

	return nil
}

// Alerting config.
type Alerting struct {
	Actions []*AlertAction `json:"actions"`
	Alerts  []*Alert       `json:"alerts"`
}

// Default implementation.
func (a *Alerting) Default() error {
	for i, a := range a.Alerts {
		if err := a.Default(); err != nil {
			return errors.Wrapf(err, ".actions %d", i)
		}
	}

	return nil
}

// Validate implementation.
func (a *Alerting) Validate() error {
	for i, a := range a.Actions {
		if err := a.Validate(); err != nil {
			return errors.Wrapf(err, ".actions %d", i)
		}
	}

	for i, alert := range a.Alerts {
		if err := alert.Validate(); err != nil {
			return errors.Wrapf(err, ".actions %d", i)
		}

		if a.GetActionByName(alert.Action) == nil {
			return errors.Errorf(".action %q is not defined", alert.Action)
		}
	}

	return nil
}

// GetActionByName returns the action by name or nil.
func (a *Alerting) GetActionByName(name string) *AlertAction {
	for _, action := range a.Actions {
		if action.Name == name {
			return action
		}
	}

	return nil
}
