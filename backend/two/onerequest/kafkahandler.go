package onerequest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http/httputil"

	"github.com/kartpop/cruncan/backend/pkg/model"
	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
	httpInternal "github.com/kartpop/cruncan/backend/two/http"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type KafkaHandler struct {
	logger *slog.Logger
	tracer trace.Tracer
	client *httpInternal.Client
}

func NewKafkaHandler(ctx context.Context, client *httpInternal.Client) *KafkaHandler {
	tracer, _ := otelContext.Tracer(ctx)

	return &KafkaHandler{
		logger: slog.Default(),
		tracer: tracer,
		client: client,
	}
}

func (h *KafkaHandler) Handle(ctx context.Context, message []byte, topic string) error {
	ctx, span := h.tracer.Start(ctx, "onerequest.kafkahandler.Handle", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	var oneRequest model.OneRequest
	err := json.Unmarshal(message, &oneRequest)
	if err != nil {
		return h.logAndMonitorError(ctx, "failed to unmarshal message", span, err)
	}

	threeRequest := &model.ThreeRequest{
		OneRequest: oneRequest,
		Metadata:   "hardcoded metadata for now",
	}
	resp, err := h.client.PostThreeRequest(ctx, threeRequest)
	if err != nil {
		return h.logAndMonitorError(ctx, "failed to post three request", span, err)
	}

	defer resp.Body.Close()
	respBody, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return h.logAndMonitorError(ctx, "failed to dump response", span, err)
	}

	h.logger.Info(fmt.Sprintf("response for ThreeRequest: %v", string(respBody)))

	span.SetStatus(codes.Ok, "request handled successfully by kafka consumer")
	return nil
}

func (h *KafkaHandler) logAndMonitorError(ctx context.Context, errPrefix string, span trace.Span, err error) error {
	errMssg := fmt.Sprintf("%s: %v", errPrefix, err)
	h.logger.ErrorContext(ctx, errMssg)
	span.SetStatus(codes.Error, errMssg)
	span.RecordError(err)
	return err
}
