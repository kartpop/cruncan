package otel

import (
	"context"

	"github.com/kartpop/cruncan/backend/pkg/util"
)

func Setup(tracerName string, meterName string) (context.Context, func()) {
	ctx := context.Background()

	// initialize OTEL logger
	ctx, cancelLogger, err := InitLogger(ctx)
	if err != nil {
		util.Fatal("error initializing logger: %v", err)
	}

	// initialize OTEL tracer
	ctx, cancelTracer, err := InitTracer(ctx, tracerName)
	if err != nil {
		util.Fatal("error initializing tracer: %v", err)
	}

	// initialize OTEL meter
	ctx, cancelMeter, err := InitMeter(ctx, meterName)
	if err != nil {
		util.Fatal("error initializing meter: %v", err)
	}

	return ctx, func() {
		cancelLogger()
		cancelTracer()
		cancelMeter()
	}
}
