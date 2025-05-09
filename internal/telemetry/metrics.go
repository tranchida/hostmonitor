package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Meter returns a named meter from the global meter provider
func Meter(name string) metric.Meter {
	return otel.Meter(name)
}

// Counter creates a new counter instrument
func Counter(name, description string) (metric.Int64Counter, error) {
	return Meter("github.com/tranchida/hostmonitor").Int64Counter(
		name,
		metric.WithDescription(description),
		metric.WithUnit("1"),
	)
}

// Gauge creates a new gauge instrument
func Gauge(name, description string) (metric.Int64ObservableGauge, error) {
	return Meter("github.com/tranchida/hostmonitor").Int64ObservableGauge(
		name,
		metric.WithDescription(description),
		metric.WithUnit("1"),
	)
}

// Histogram creates a new histogram instrument
func Histogram(name, description string) (metric.Int64Histogram, error) {
	return Meter("github.com/tranchida/hostmonitor").Int64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithUnit("1"),
	)
}

// RecordCounter records a count with the given counter
func RecordCounter(ctx context.Context, counter metric.Int64Counter, value int64, attrs ...attribute.KeyValue) {
	counter.Add(ctx, value, metric.WithAttributes(attrs...))
}

// RecordHistogram records a value with the given histogram
func RecordHistogram(ctx context.Context, histogram metric.Int64Histogram, value int64, attrs ...attribute.KeyValue) {
	histogram.Record(ctx, value, metric.WithAttributes(attrs...))
}

// RegisterCallback registers a callback function for observable instruments
func RegisterCallback(meter metric.Meter, callback metric.Callback, instruments ...metric.Observable) (metric.Registration, error) {
	return meter.RegisterCallback(callback, instruments...)
}
