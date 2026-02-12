package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	promexporter "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Init initializes the OpenTelemetry MeterProvider with a Prometheus exporter.
// Returns a shutdown function that should be deferred by the caller.
func Init(ctx context.Context) (func(context.Context) error, error) {
	exporter, err := promexporter.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)

	otel.SetMeterProvider(mp)

	return mp.Shutdown, nil
}

// Handler returns an http.Handler that serves Prometheus metrics.
func Handler() http.Handler {
	return promhttp.Handler()
}

// ServeMetrics starts a lightweight HTTP server that exposes a /metrics endpoint.
// This is intended for gRPC services that don't have an HTTP server.
// It blocks, so it should be called in a goroutine.
func ServeMetrics(port string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":"+port, mux)
}
