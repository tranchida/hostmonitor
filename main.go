package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tranchida/hostmonitor/internal/handler"
	"github.com/tranchida/hostmonitor/internal/telemetry"
	"go.uber.org/zap"
)

//go:embed template static
var contentFS embed.FS

func newEngine(logger *zap.Logger) *gin.Engine {
	e := gin.New()
	
	// Add OpenTelemetry middleware
	e.Use(telemetry.TracingMiddleware("hostmonitor"))
	e.Use(telemetry.LoggingMiddleware(logger))
	e.Use(telemetry.MetricsMiddleware("hostmonitor"))
	e.Use(gin.Recovery())

	staticfs, _ := fs.Sub(contentFS, "static")
	e.StaticFS("/static", http.FS(staticfs))

	e.GET("/", handler.IndexHandler)
	e.GET("/host", handler.HostInfoHandler)

	return e
}

func main() {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize OpenTelemetry with OTLP exporter
	telemetryProvider, err := telemetry.NewProvider(ctx, telemetry.Config{
		ServiceName:    "hostmonitor",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		OTLPEndpoint:   "localhost:4317", // Default OTLP gRPC endpoint
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize telemetry: %v", err))
	}
	defer telemetryProvider.Shutdown(ctx)

	// Use the logger from the telemetry provider
	logger := telemetryProvider.Logger
	defer logger.Sync()
	
	// Initialize system metrics collector
	metricsCollector := telemetry.NewSystemMetricsCollector(logger, 15*time.Second)
	if err := metricsCollector.Start(ctx); err != nil {
		logger.Error("Failed to start system metrics collector", zap.Error(err))
	}
	defer metricsCollector.Stop(ctx)

	// Create a new engine with the logger
	engine := newEngine(logger)

	// Create a server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		fmt.Println("open browser on : http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	
	// Create a deadline for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	// Shutdown the server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}

	logger.Info("Server gracefully stopped")
}

// Application version and build information
var (
	AppVersion = "1.0.0"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)
