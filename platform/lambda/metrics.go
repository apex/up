package lambda

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/apex/up/internal/metrics"
	"github.com/apex/up/platform/event"
)

// TODO: write a higher level pkg in tj/aws
// TODO: move the metrics pkg to tj/aws

type stat struct {
	Name   string
	Metric string
	Stat   string
	point  *cloudwatch.Datapoint
}

// Value returns the metric value.
func (s *stat) Value() int {
	if s.point == nil {
		return 0
	}

	switch s.Stat {
	case "Sum":
		return int(*s.point.Sum)
	case "Average":
		return int(*s.point.Average)
	case "Minimum":
		return int(*s.point.Minimum)
	case "Maximum":
		return int(*s.point.Maximum)
	default:
		return 0
	}
}

// stats to fetch.
var stats = []*stat{
	{"Requests", "Count", "Sum", nil},
	{"Avg Latency", "Latency", "Average", nil},
	{"Min Latency", "Latency", "Minimum", nil},
	{"Max Latency", "Latency", "Maximum", nil},
	{"Client Errors", "4XXError", "Sum", nil},
	{"Server Errors", "5XXError", "Sum", nil},
}

// ShowMetrics implementation.
func (p *Platform) ShowMetrics(region, stage string, start time.Time) error {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := cloudwatch.New(s)

	errc := make(chan error, len(stats))
	var wg sync.WaitGroup

	d := time.Now().Sub(start)
	period := int(d / time.Second)

	for _, s := range stats {
		wg.Add(1)
		go func(s *stat) {
			defer wg.Done()

			m := metrics.New().
				Namespace("AWS/ApiGateway").
				Dimension("ApiName", p.config.Name).
				Dimension("Stage", stage).
				TimeRange(time.Now().Add(-d), time.Now()).
				Period(period).
				Stat(s.Stat).
				Metric(s.Metric)

			res, err := c.GetMetricStatistics(m.Params())
			if err != nil {
				errc <- err
				return
			}

			if len(res.Datapoints) > 0 {
				s.point = res.Datapoints[0]
			}
		}(s)
	}

	wg.Wait()

	select {
	case err := <-errc:
		return err
	default:
	}

	for _, s := range stats {
		p.events.Emit("metrics.value", event.Fields{
			"name":  s.Name,
			"value": s.Value(),
		})
	}

	return nil
}
