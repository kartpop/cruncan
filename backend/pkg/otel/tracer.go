package otel

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes the tracer and returns a context with the tracer and a function to shut down the tracer.
// The tracer is configured to send traces to the otel-collector.
func InitTracer(ctx context.Context, name string, opts ...trace.TracerOption) (newCtx context.Context, cancel func(), err error) {

	// Initialize the OTLP exporter using environment variables for configuration.
	exporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())
	if err != nil {
		return ctx, noop, fmt.Errorf("failed to create trace exporter: %v", err)
	}

	res, err := resource.New(ctx,
		// The service name is now picked up from the OTEL_SERVICE_NAME environment variable.
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
		resource.WithHost(),      // This option configures a set of Detectors that discover host information
	)

	if err != nil {
		return ctx, noop, fmt.Errorf("failed to create trace resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	ctx = otelContext.WithTracer(ctx, otel.Tracer(name, opts...))

	return ctx, func() {
		ctx, cancelDeadline := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
		defer cancelDeadline()
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down tracer provider: %v", err)
		}
	}, nil
}
