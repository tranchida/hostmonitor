package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Tracer returns a named tracer from the global tracer provider
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// StartSpan starts a new span with the given name and returns the context and span
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer("github.com/tranchida/hostmonitor").Start(ctx, name, opts...)
}

// AddAttributes adds attributes to the span in the current context
func AddAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordError records an error in the span from the current context
func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err, opts...)
	span.SetStatus(codes.Error, err.Error())
}

// EndSpan ends the span in the current context
func EndSpan(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	span.End()
}

// WithSpan wraps a function with a new span
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	ctx, span := StartSpan(ctx, name)
	defer span.End()
	
	err := fn(ctx)
	if err != nil {
		RecordError(ctx, err)
	}
	
	return err
}
