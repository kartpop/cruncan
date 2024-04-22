package context

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	noopTracer "go.opentelemetry.io/otel/trace/noop"
	"log/slog"
)

type meterKey struct{}
type tracerKey struct{}
type loggerKey struct{}
type slogAttributesHandlerKey struct{}

// WithMeterFromGlobalProvider returns a new context with the metric.Meter
func WithMeterFromGlobalProvider(ctx context.Context, meterName string, opts ...metric.MeterOption) context.Context {
	return WithMeter(ctx, otel.GetMeterProvider().Meter(meterName, opts...))
}

// WithTracerFromGlobalProvider returns a new context with the trace.Tracer
func WithTracerFromGlobalProvider(ctx context.Context, tracerName string, opts ...trace.TracerOption) context.Context {
	return WithTracer(ctx, otel.GetTracerProvider().Tracer(tracerName, opts...))
}

// WithLoggerFromGlobalProvider returns a new context with the *slog.Logger
func WithLoggerFromGlobalProvider(ctx context.Context) context.Context {
	return WithLogger(ctx, slog.Default())
}

// WithMeter returns a new context with the given metric.Meter associated.
// The returned context carries a value with a key of type meterKey. This
// context can then be used with Meter to retrieve the stored metric.Meter.
func WithMeter(ctx context.Context, meter metric.Meter) context.Context {
	return context.WithValue(ctx, meterKey{}, meter)
}

// Meter extracts the metric.Meter from the context if it exists.
// The function returns the metric.Meter and a boolean indicating whether
// the metric.Meter was found in the context.
func Meter(ctx context.Context) (metric.Meter, bool) {
	val := ctx.Value(meterKey{})
	if val == nil {
		return &noop.Meter{}, false
	}
	return val.(metric.Meter), true
}

// WithTracer returns a new context with the given trace.Tracer associated.
// The returned context carries a value with a key of type tracerKey. This
// context can then be used with Tracer to retrieve the stored trace.Tracer.
func WithTracer(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

// Tracer extracts the trace.Tracer from the context if it exists.
// The function returns the trace.Tracer and a boolean indicating whether
// the trace.Tracer was found in the context.
func Tracer(ctx context.Context) (trace.Tracer, bool) {
	val := ctx.Value(tracerKey{})
	if val == nil {
		return noopTracer.Tracer{}, false
	}
	return val.(trace.Tracer), true
}

// WithLogger returns a new context with the given *slog.Logger associated.
// The returned context carries a value with a key of type loggerKey. This
// context can then be used with Logger to retrieve the stored *slog.Logger.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// WithSlogAttributes returns a new context with the given slog.Attributes associated.
func WithSlogAttributes(ctx context.Context, attr ...slog.Attr) context.Context {
	return context.WithValue(ctx, slogAttributesHandlerKey{}, attr)
}

// AddSlogAttributes returns a new context with the given slog.Attributes added to the existing attributes.
func AddSlogAttributes(ctx context.Context, attr ...slog.Attr) context.Context {
	val := ctx.Value(slogAttributesHandlerKey{})
	if val == nil {
		return context.WithValue(ctx, slogAttributesHandlerKey{}, attr)
	}
	attr = append(val.([]slog.Attr), attr...)
	return context.WithValue(ctx, slogAttributesHandlerKey{}, attr)
}

// SlogAttributes extracts the slog.Attributes from the context if it exists.
func SlogAttributes(ctx context.Context) ([]slog.Attr, bool) {
	val := ctx.Value(slogAttributesHandlerKey{})
	if val == nil {
		return nil, false
	}
	return val.([]slog.Attr), true
}

// Logger extracts the *slog.Logger from the context if it exists.
// The function returns the *slog.Logger and a boolean indicating whether
// the *slog.Logger was found in the context.
func Logger(ctx context.Context) (*slog.Logger, bool) {
	val := ctx.Value(loggerKey{})
	if val == nil {
		return nil, false
	}
	return val.(*slog.Logger), true
}
