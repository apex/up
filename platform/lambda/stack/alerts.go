package stack

import (
	"strconv"

	"github.com/apex/up"
	"github.com/apex/up/config"
	"github.com/apex/up/internal/util"
)

// alertActionID returns the alert action id.
func alertActionID(name string) string {
	return util.Camelcase("alert_action_%s", name)
}

// alert resources.
func alert(c *up.Config, a *config.Alert, m Map) {
	period := a.Period.Seconds()
	alertAction := ref(alertActionID(a.Action))
	id := util.Camelcase("alert_%s_%s_%s_period_%d", a.Namespace, a.Metric, a.Statistic, int(period))

	m[id] = Map{
		"Type": "AWS::CloudWatch::Alarm",
		"Properties": Map{
			"ActionsEnabled":     !a.Disable,
			"AlarmDescription":   a.Description,
			"MetricName":         a.Metric,
			"Namespace":          a.Namespace,
			"Statistic":          a.Statistic,
			"Period":             period,
			"EvaluationPeriods":  "1",
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
func action(c *up.Config, a *config.AlertAction, m Map) {
	id := alertActionID(a.Name)

	m[id] = Map{
		"Type": "AWS::SNS::Topic",
		"Properties": Map{
			"DisplayName": a.Name,
		},
	}

	for _, email := range a.Emails {
		sub := util.Camelcase("alert_action_%s_subscription_%s", a.Name, email)

		m[sub] = Map{
			"Type": "AWS::SNS::Subscription",
			"Properties": Map{
				"Endpoint": email,
				"Protocol": "email",
				"TopicArn": ref(id),
			},
		}
	}
}

// alerting resources.
func alerting(c *up.Config, m Map) {
	for _, a := range c.Alerting.Alerts {
		alert(c, a, m)
	}

	for _, a := range c.Alerting.Actions {
		action(c, a, m)
	}
}
