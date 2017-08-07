// Package metrics provides higher level CloudWatch metrics operations.
package metrics

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

// DefaultClient is the default client used for cloudwatch.
var DefaultClient = cloudwatch.New(session.New(aws.NewConfig()))

// Metrics helper.
type Metrics struct {
	client cloudwatchiface.CloudWatchAPI
	in     cloudwatch.GetMetricStatisticsInput
}

// New metrics.
func New() *Metrics {
	return &Metrics{
		client: DefaultClient,
	}
}

// Client sets the client.
func (m *Metrics) Client(c cloudwatchiface.CloudWatchAPI) *Metrics {
	m.client = c
	return m
}

// Namespace sets the namespace.
func (m *Metrics) Namespace(name string) *Metrics {
	m.in.Namespace = &name
	return m
}

// Metric sets the metric name.
func (m *Metrics) Metric(name string) *Metrics {
	m.in.MetricName = &name
	return m
}

// Stats sets the stats.
func (m *Metrics) Stats(names []string) *Metrics {
	m.in.Statistics = aws.StringSlice(names)
	return m
}

// Stat adds the stat.
func (m *Metrics) Stat(name string) *Metrics {
	m.in.Statistics = append(m.in.Statistics, &name)
	return m
}

// Dimension adds a dimension.
func (m *Metrics) Dimension(name, value string) *Metrics {
	m.in.Dimensions = append(m.in.Dimensions, &cloudwatch.Dimension{
		Name:  &name,
		Value: &value,
	})

	return m
}

// Period sets the period.
func (m *Metrics) Period(minutes int64) *Metrics {
	m.in.Period = &minutes
	return m
}

// TimeRange sets the start and time times.
func (m *Metrics) TimeRange(start, end time.Time) *Metrics {
	m.in.StartTime = &start
	m.in.EndTime = &end
	return m
}

// Input returns the API input.
func (m *Metrics) Input() *cloudwatch.GetMetricStatisticsInput {
	return &m.in
}

// Get metrics.
func Get(m *Metrics) (*cloudwatch.GetMetricStatisticsOutput, error) {
	return m.client.GetMetricStatistics(m.Input())
}
