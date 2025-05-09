package telemetry

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TracingMiddleware returns a gin middleware that adds OpenTelemetry tracing to requests
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		// Extract context from the request headers
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start a new span
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// Set span attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.host", c.Request.Host),
			attribute.String("http.user_agent", c.Request.UserAgent()),
		)

		// Store the context in the request
		c.Request = c.Request.WithContext(ctx)

		// Process the request
		c.Next()

		// Update span with response information
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		// Record errors if any
		if len(c.Errors) > 0 {
			span.SetStatus(codes.Error, c.Errors.String())
			span.RecordError(fmt.Errorf("%s", c.Errors.String()))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	}
}

// LoggingMiddleware returns a gin middleware that logs requests using zap
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the current span from context
		span := trace.SpanFromContext(c.Request.Context())
		traceID := span.SpanContext().TraceID().String()

		// Process the request
		c.Next()

		// Log the request with trace ID
		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("trace_id", traceID),
			zap.Int("size", c.Writer.Size()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// MetricsMiddleware returns a gin middleware that records metrics for requests
func MetricsMiddleware(serviceName string) gin.HandlerFunc {
	meter := otel.Meter(serviceName)
	
	// Create instruments
	requestCounter, _ := meter.Int64Counter(
		"http.server.request_count",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("1"),
	)

	return func(c *gin.Context) {
		// Process the request
		c.Next()

		// Record metrics
		requestCounter.Add(c.Request.Context(), 1,
			metric.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.route", c.FullPath()),
				attribute.Int("http.status_code", c.Writer.Status()),
			),
		)
	}
}
