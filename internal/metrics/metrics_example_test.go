package metrics_test

import (
	"time"

	"github.com/apex/up/internal/metrics"
)

func Example() {
	m := metrics.New().
		Namespace("AWS/ApiGateway").
		Metric("Count").
		Stat("Sum").
		Dimension("ApiName", "app").
		Period(5).
		TimeRange(time.Now().Add(-time.Hour), time.Now())

	res, err := metrics.Get(m)
	_ = res
	_ = err
}
