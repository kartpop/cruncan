package otel

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/agoda-com/opentelemetry-go/otelslog"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	logsSdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/otel/sdk/resource"
)

type LoggerConfig struct {
	extraHandler []slog.Handler
	level        slog.Level
	testBuffer   *bytes.Buffer
}

type LoggerOption interface {
	apply(LoggerConfig) LoggerConfig
	prio() int
}

type loggerOptionImpl struct {
	f func(LoggerConfig) LoggerConfig
	p int
}

func (o loggerOptionImpl) apply(c LoggerConfig) LoggerConfig {
	return o.f(c)
}

func (o loggerOptionImpl) prio() int {
	return o.p
}

func newLoggerOption(f func(LoggerConfig) LoggerConfig, p int) LoggerOption {
	return &loggerOptionImpl{
		f: f,
		p: p,
	}
}

// WithExtraHandler adds an extra handler to the logger
func WithExtraHandler(handler slog.Handler) LoggerOption {
	return newLoggerOption(func(cfg LoggerConfig) LoggerConfig {
		cfg.extraHandler = append(cfg.extraHandler, handler)
		return cfg
	}, 0)
}

// WithConsoleHandler adds a console handler to the logger
// Add this after WithLogLevel
func WithConsoleHandler() LoggerOption {
	return newLoggerOption(func(cfg LoggerConfig) LoggerConfig {
		var w io.Writer
		if cfg.testBuffer != nil {
			w = cfg.testBuffer
		} else {
			w = os.Stdout
		}
		cfg.extraHandler = append(cfg.extraHandler, slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: cfg.level,
		}))
		return cfg
	}, 100)
}

// WithTestBuffer configures the logger to send logs to the provided buffer
func WithTestBuffer(buffer *bytes.Buffer) LoggerOption {
	return newLoggerOption(func(cfg LoggerConfig) LoggerConfig {
		cfg.testBuffer = buffer
		return cfg
	}, -999998)
}

// WithLevel sets the log level for the logger
func WithLevel(level slog.Level) LoggerOption {
	return newLoggerOption(func(cfg LoggerConfig) LoggerConfig {
		cfg.level = level
		return cfg
	}, -999999)
}

// WithEnvLevel sets the log level for the logger from the LOG_LEVEL environment variable
// or from the first environment variable that is set from the list of environment variables.
func WithEnvLevel(env ...string) LoggerOption {
	var ls string
	for _, e := range env {
		if ls == "" {
			ls = os.Getenv(e)
		}
	}
	if ls == "" {
		ls = os.Getenv("LOG_LEVEL")
	}
	switch strings.ToUpper(ls) {
	case "DEBUG":
		return WithLevel(slog.LevelDebug)
	case "INFO":
		return WithLevel(slog.LevelInfo)
	case "WARN":
		return WithLevel(slog.LevelWarn)
	case "ERROR":
		return WithLevel(slog.LevelError)
	default:
		if ls != "" {
			v, err := strconv.Atoi(ls)
			if err == nil {
				return WithLevel(slog.Level(v))
			}
		}
		return WithLevel(slog.LevelInfo)
	}
}

// InitLogger initializes the logger and returns a context with the logger and a function to shut down the logger.
// The logger is configured to send logs to the otel-collector.
// Use WithConsoleHandler to add a console handler to the logger.
// Use WithExtraHandler to add extra custom handlers to the logger.
func InitLogger(ctx context.Context, loggerOptions ...LoggerOption) (context.Context, func(), error) {

	loggerOptions, cfg := renderConfig(loggerOptions)

	// configure opentelemetry logger provider
	logExporter, err := otlplogs.NewExporter(ctx)
	if err != nil {
		return ctx, noop, fmt.Errorf("failed to create log exporter: %v", err)
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
		return ctx, noop, fmt.Errorf("failed to create log resource: %v", err)
	}

	loggerProvider := logsSdk.NewLoggerProvider(
		logsSdk.WithBatcher(logExporter),
		logsSdk.WithResource(res),
	)

	var handler slog.Handler

	handler = otelslog.NewOtelHandler(loggerProvider, &otelslog.HandlerOptions{
		Level: cfg.level,
	})

	if len(cfg.extraHandler) > 0 {
		allHandlers := append(cfg.extraHandler, handler)
		handler = slogmulti.Fanout(allHandlers...)
	}

	otelLogger := slog.New(NewSlogOTELAttributesHandler(handler))

	//configure default logger
	slog.SetDefault(otelLogger)

	ctx = otelContext.WithLogger(ctx, otelLogger)

	return ctx, func() {
		ctx, cancelDeadline := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
		defer cancelDeadline()
		if err := loggerProvider.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down tracer provider: %v", err)
		}
	}, nil
}

func renderConfig(loggerOptions []LoggerOption) ([]LoggerOption, LoggerConfig) {
	if len(loggerOptions) > 1 {
		sort.Slice(loggerOptions, func(i, j int) bool {
			return loggerOptions[i].prio() < loggerOptions[j].prio()
		})
	}

	// set log level from env if not provided by any option
	fakeCfg := LoggerConfig{
		level: -992459912,
	}
	for _, opt := range loggerOptions {
		fakeCfg = opt.apply(fakeCfg)
	}
	if fakeCfg.level == -992459912 {
		l := []LoggerOption{WithEnvLevel()}
		loggerOptions = append(l, loggerOptions...)
	}

	var cfg LoggerConfig
	for _, opt := range loggerOptions {
		cfg = opt.apply(cfg)
	}
	return loggerOptions, cfg
}

func noop() {
	// no-op
}
