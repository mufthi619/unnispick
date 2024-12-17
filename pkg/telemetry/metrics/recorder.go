package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	meterName = "ecommerce"
)

type Recorder struct {
	meter metric.Meter
}

func NewRecorder() *Recorder {
	return &Recorder{
		meter: otel.Meter(meterName),
	}
}

func (r *Recorder) Counter(name, description string, opts ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return r.meter.Int64Counter(name, metric.WithDescription(description))
}

func (r *Recorder) Histogram(name, description string, opts ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return r.meter.Float64Histogram(name, metric.WithDescription(description))
}
