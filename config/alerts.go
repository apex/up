package config

import (
	"strings"
	"time"

	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/validate"
	"github.com/pkg/errors"
)

var namespaceMap = map[string]string{
	"http": "AWS/ApiGateway",
}

var metricsMap = map[string]string{
	"http.count":   "Count",
	"http.latency": "Latency",
	"http.4xx":     "4XXError",
	"http.5xx":     "5XXError",
}

var operatorMap = map[string]string{
	">":  "GreaterThanThreshold",
	"<":  "LessThanThreshold",
	">=": "GreaterThanOrEqualToThreshold",
	"<=": "LessThanOrEqualToThreshold",
}

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

// Action config.
type Action struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Email string `json:"email"` // TODO: decide
}

// Validate implementation.
func (a *Action) Validate() error {
	// TODO: implement
	return nil
}

// Alert config.
type Alert struct {
	Description string   `json:"description"`
	Disable     bool     `json:"disable"`
	Metric      string   `json:"metric"`
	Namespace   string   `json:"namespace"`
	Statistic   string   `json:"statistic"`
	Operator    string   `json:"operator"`
	Threshold   int      `json:"threshold"`
	Period      Duration `json:"period"`
	Stage       string   `json:"stage"`
	Action      string   `json:"action"`
}

// Default implementation.
func (a *Alert) Default() error {
	if a.Operator == "" {
		a.Operator = ">"
	}

	if a.Period == 0 {
		a.Period = Duration(time.Minute)
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
	// operator
	if s, ok := operatorMap[a.Operator]; ok {
		a.Operator = s
	} else {
		return errors.Wrap(validate.List(a.Operator, util.StringMapKeys(operatorMap)), ".operator")
	}

	// statistic
	if s, ok := statisticMap[a.Statistic]; ok {
		a.Statistic = s
	} else {
		return errors.Wrap(validate.List(a.Statistic, util.StringMapKeys(statisticMap)), ".statistic")
	}

	// metric
	if err := validate.RequiredString(a.Metric); err != nil {
		return errors.Wrap(err, ".metric")
	}

	// namespace
	if err := validate.RequiredString(a.Namespace); err != nil {
		return errors.Wrap(err, ".namespace")
	}

	// action
	if err := validate.RequiredString(a.Action); err != nil {
		return errors.Wrap(err, ".action")
	}

	return nil
}

// Alerting config.
type Alerting struct {
	Actions []*Action `json:"actions"`
	Alerts  []*Alert  `json:"alerts"`
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

	for i, a := range a.Alerts {
		if err := a.Validate(); err != nil {
			return errors.Wrapf(err, ".actions %d", i)
		}
	}

	return nil
}
