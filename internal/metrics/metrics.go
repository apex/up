// Package metrics provides higher level CloudWatch metrics operations.
package metrics

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics helper.
type Metrics struct {
	in cloudwatch.GetMetricStatisticsInput
}

// New metrics.
func New() *Metrics {
	return &Metrics{}
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

// Period sets the period in seconds.
func (m *Metrics) Period(seconds int) *Metrics {
	m.in.Period = aws.Int64(int64(seconds))
	return m
}

// TimeRange sets the start and time times.
func (m *Metrics) TimeRange(start, end time.Time) *Metrics {
	m.in.StartTime = &start
	m.in.EndTime = &end
	return m
}

// Params returns the API input.
func (m *Metrics) Params() *cloudwatch.GetMetricStatisticsInput {
	return &m.in
}
