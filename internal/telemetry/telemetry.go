package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// Config holds configuration for telemetry setup
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string // OTLP endpoint, e.g., "localhost:4317"
	// Add more configuration options as needed
}

// Provider holds the OpenTelemetry providers and other components
type Provider struct {
	TracerProvider *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
	Logger         *zap.Logger
	Propagator     propagation.TextMapPropagator
	// Add more providers as needed
}

// NewProvider creates and configures a new telemetry provider
func NewProvider(ctx context.Context, cfg Config) (*Provider, error) {
	// Create a resource describing the service
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Configure trace provider with OTLP exporter
	// Default to localhost:4317 if not specified
	otlpEndpoint := cfg.OTLPEndpoint
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4317"
	}
	
	// Create gRPC connection to the collector
	conn, err := grpc.Dial(otlpEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Create OTLP trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, 
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithTimeout(5 * time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(traceProvider)

	// Configure metric provider with OTLP exporter
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
		otlpmetricgrpc.WithTimeout(5 * time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(15*time.Second))),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	// Configure propagator
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	// Configure logger
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &Provider{
		TracerProvider: traceProvider,
		MeterProvider:  meterProvider,
		Logger:         logger,
		Propagator:     propagator,
	}, nil
}

// Shutdown gracefully shuts down the telemetry provider
func (p *Provider) Shutdown(ctx context.Context) error {
	var shutdownErrs []error

	// Shutdown trace provider
	if err := p.TracerProvider.Shutdown(ctx); err != nil {
		shutdownErrs = append(shutdownErrs, fmt.Errorf("failed to shutdown tracer provider: %w", err))
	}

	// Shutdown meter provider
	if err := p.MeterProvider.Shutdown(ctx); err != nil {
		shutdownErrs = append(shutdownErrs, fmt.Errorf("failed to shutdown meter provider: %w", err))
	}

	// Sync logger
	if err := p.Logger.Sync(); err != nil {
		shutdownErrs = append(shutdownErrs, fmt.Errorf("failed to sync logger: %w", err))
	}

	if len(shutdownErrs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", shutdownErrs)
	}
	return nil
}
