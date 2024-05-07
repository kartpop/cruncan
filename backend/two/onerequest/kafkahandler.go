package onerequest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/kartpop/cruncan/backend/pkg/model"
	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type KafkaHandler struct {
	logger *slog.Logger
	tracer trace.Tracer
}

func NewKafkaHandler(ctx context.Context) *KafkaHandler {
	tracer, _ := otelContext.Tracer(ctx)

	return &KafkaHandler{
		logger: slog.Default(),
		tracer: tracer,
	}
}

func (h *KafkaHandler) Handle(ctx context.Context, message []byte, topic string) error {
	ctx, span := h.tracer.Start(ctx, "onerequest.kafkahandler.Handle", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	var oneRequest model.OneRequest
	err := json.Unmarshal(message, &oneRequest)
	if err != nil {
		errMsg := fmt.Sprintf("failed to unmarshal message: %v", err)
		h.logAndMonitorError(ctx, errMsg, span, err)
		return err
	}

	h.logger.Info(fmt.Sprintf("Unmarshaled message: %v", oneRequest))
	span.SetStatus(codes.Ok, "request handled successfully by kafka consumer")

	return nil
}

func (h *KafkaHandler) logAndMonitorError(ctx context.Context, errMsg string, span trace.Span, err error) {
	h.logger.ErrorContext(ctx, errMsg)
	span.SetStatus(codes.Error, errMsg)
	span.RecordError(err)
}
