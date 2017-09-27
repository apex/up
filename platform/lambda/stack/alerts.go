package stack

import (
	"strconv"

	"github.com/apex/up"
	"github.com/apex/up/config"
	"github.com/apex/up/internal/util"
)

// alert sets up alarms.
func alert(c *up.Config, a *config.Alert, m Map) {
	id := util.Camelcase("alert_alarm_%s_%s_%s", a.Metric, a.Statistic, a.Action)
	topicID := util.Camelcase("alert_topic_%s", a.Action)

	m[id] = Map{
		"Type": "AWS::CloudWatch::Alarm",
		"Properties": Map{
			"ActionsEnabled":     !a.Disable,
			"AlarmDescription":   a.Description,
			"MetricName":         a.Metric,
			"Namespace":          a.Namespace,
			"Statistic":          a.Statistic,
			"Period":             strconv.Itoa(int(a.Period.Seconds())),
			"EvaluationPeriods":  "1",
			"Threshold":          strconv.Itoa(a.Threshold),
			"ComparisonOperator": a.Operator,
			"OKActions": []Map{
				ref(topicID),
			},
			"AlarmActions": []Map{
				ref(topicID),
			},
			"Dimensions": []Map{ // TODO: allow passing other dimensions
				{
					"Name":  "ApiName",
					"Value": ref("Name"),
				},
				{
					"Name":  "Stage",
					"Value": "production",
				},
			},
		},
	}
}

// action sets up SNS action triggers.
func action(c *up.Config, a *config.Action, m Map) {
	id := util.Camelcase("alert_topic_%s", a.Name)
	name := util.Camelcase("alert_%s_%s", c.Name, a.Name)

	m[id] = Map{
		"Type": "AWS::SNS::Topic",
		"Properties": Map{
			"DisplayName": a.Name,
			"TopicName":   name,
		},
	}

	sid := util.Camelcase("alert_topic_%s_subscription", a.Name)
	m[sid] = Map{
		"Type": "AWS::SNS::Subscription",
		"Properties": Map{
			"Endpoint": a.Email,
			"Protocol": "email",
			"TopicArn": ref(id),
		},
	}
}

// alerting sets up alerts and actions.
func alerting(c *up.Config, m Map) {
	for _, a := range c.Alerting.Alerts {
		alert(c, a, m)
	}

	for _, a := range c.Alerting.Actions {
		action(c, a, m)
	}
}
