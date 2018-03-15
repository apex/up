package resources

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/apex/up/config"
	"github.com/apex/up/internal/util"
)

// alertActionID returns the alert action id.
func alertActionID(name string) string {
	return util.Camelcase("alert_action_%s", name)
}

// alert resources.
func alert(c *Config, a *config.Alert, m Map) {
	period := a.Period.Seconds()
	alertAction := ref(alertActionID(a.Action))
	id := util.Camelcase("alert_%s_%s_%s_period_%d_threshold_%d", a.Namespace, a.Metric, a.Statistic, int(period), a.Threshold)

	m[id] = Map{
		"Type": "AWS::CloudWatch::Alarm",
		"Properties": Map{
			"ActionsEnabled":     !a.Disable,
			"AlarmDescription":   a.Description,
			"MetricName":         a.Metric,
			"Namespace":          a.Namespace,
			"Statistic":          a.Statistic,
			"TreatMissingData":   a.Missing,
			"Period":             period,
			"EvaluationPeriods":  strconv.Itoa(a.EvaluationPeriods),
			"Threshold":          strconv.Itoa(a.Threshold),
			"ComparisonOperator": a.Operator,
			"OKActions":          []Map{alertAction},
			"AlarmActions":       []Map{alertAction},
			"Dimensions": []Map{
				{
					"Name":  "ApiName",
					"Value": c.Name,
				},
				{
					"Name":  "Stage",
					"Value": "production",
				},
			},
		},
	}
}

// action resources.
func action(c *Config, a *config.AlertAction, m Map) {
	id := alertActionID(a.Name)

	m[id] = Map{
		"Type": "AWS::SNS::Topic",
		"Properties": Map{
			"DisplayName": a.Name,
		},
	}

	v := url.Values{}
	v.Add("type", a.Type)
	v.Add("emails", strings.Join(a.Emails, ","))
	v.Add("numbers", strings.Join(a.Numbers, ","))
	v.Add("url", a.URL)
	v.Add("channel", a.Channel)

	if a.Gifs {
		v.Add("gifs", "1")
	}

	sub := util.Camelcase("alert_action_%s_subscription", a.Name)
	url := fmt.Sprintf("https://up.apex.sh/alert?%s", v.Encode())

	m[sub] = Map{
		"Type": "AWS::SNS::Subscription",
		"Properties": Map{
			"Endpoint": url,
			"Protocol": "https",
			"TopicArn": ref(id),
		},
	}
}

// actionEmail resource.
func actionEmail(c *Config, a *config.AlertAction, m Map, id string) {
	emails := strings.Join(a.Emails, ",")
	sub := util.Camelcase("alert_action_%s_subscription", a.Name)
	url := fmt.Sprintf("https://up.apex.sh/alert?emails=%s", emails)

	m[sub] = Map{
		"Type": "AWS::SNS::Subscription",
		"Properties": Map{
			"Endpoint": url,
			"Protocol": "https",
			"TopicArn": ref(id),
		},
	}
}

// alerting resources.
func alerting(c *Config, m Map) {
	for _, a := range c.Alerting.Alerts {
		alert(c, a, m)
	}

	for _, a := range c.Alerting.Actions {
		action(c, a, m)
	}
}
