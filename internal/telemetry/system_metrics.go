package telemetry

import (
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

// SystemMetricsCollector collects system metrics using OpenTelemetry
type SystemMetricsCollector struct {
	meter        metric.Meter
	logger       *zap.Logger
	registration metric.Registration
	interval     time.Duration
}

// NewSystemMetricsCollector creates a new system metrics collector
func NewSystemMetricsCollector(logger *zap.Logger, interval time.Duration) *SystemMetricsCollector {
	return &SystemMetricsCollector{
		meter:    otel.Meter("github.com/tranchida/hostmonitor/system"),
		logger:   logger,
		interval: interval,
	}
}

// Start begins collecting system metrics at the specified interval
func (c *SystemMetricsCollector) Start(ctx context.Context) error {
	// Create observable gauges for system metrics
	cpuUsage, err := c.meter.Float64ObservableGauge(
		"system.cpu.usage",
		metric.WithDescription("CPU usage percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return err
	}

	memUsage, err := c.meter.Float64ObservableGauge(
		"system.memory.usage",
		metric.WithDescription("Memory usage percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return err
	}

	diskUsage, err := c.meter.Float64ObservableGauge(
		"system.disk.usage",
		metric.WithDescription("Disk usage percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return err
	}

	goroutines, err := c.meter.Int64ObservableGauge(
		"system.goroutines",
		metric.WithDescription("Number of goroutines"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	// Register the callback function
	c.registration, err = c.meter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			// Collect CPU metrics
			cpuPercent, err := cpu.Percent(0, false)
			if err != nil {
				c.logger.Error("Failed to collect CPU metrics", zap.Error(err))
			} else if len(cpuPercent) > 0 {
				observer.ObserveFloat64(cpuUsage, cpuPercent[0], metric.WithAttributes(attribute.String("host", "local")))
			}

			// Collect memory metrics
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				c.logger.Error("Failed to collect memory metrics", zap.Error(err))
			} else {
				observer.ObserveFloat64(memUsage, memInfo.UsedPercent, metric.WithAttributes(attribute.String("host", "local")))
			}

			// Collect disk metrics
			diskInfo, err := disk.Usage("/")
			if err != nil {
				c.logger.Error("Failed to collect disk metrics", zap.Error(err))
			} else {
				observer.ObserveFloat64(diskUsage, diskInfo.UsedPercent, metric.WithAttributes(attribute.String("host", "local")))
			}

			// Collect Go runtime metrics
			observer.ObserveInt64(goroutines, int64(runtime.NumGoroutine()), metric.WithAttributes(attribute.String("host", "local")))

			return nil
		},
		cpuUsage, memUsage, diskUsage, goroutines,
	)

	if err != nil {
		return err
	}

	c.logger.Info("System metrics collector started", zap.Duration("interval", c.interval))
	return nil
}

// Stop stops the metrics collection
func (c *SystemMetricsCollector) Stop(ctx context.Context) error {
	if c.registration != nil {
		if err := c.registration.Unregister(); err != nil {
			return err
		}
		c.logger.Info("System metrics collector stopped")
	}
	return nil
}
