package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	onerequest "github.com/kartpop/cruncan/backend/one/database/onerequest"
	"github.com/kartpop/cruncan/backend/pkg/id"
	kafkaUtil "github.com/kartpop/cruncan/backend/pkg/kafka"
	"github.com/kartpop/cruncan/backend/pkg/model"
	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
	"github.com/kartpop/cruncan/backend/pkg/util"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	// FailedPostMeterName is the name of the failed post meter
	FailedPostMeterName = "http.handlers.Post.failed"
	// SuccessPostMeterName is the name of the success post meter
	SuccessPostMeterName = "http.handlers.Post.success"
	// HandledPostMeterName is the name of the handled post meter
	HandledPostMeterName = "http.handlers.Post.handled"
)

type Response struct {
	ReqID string `json:"request_id"`
}

type Handler struct {
	repo             onerequest.Repository
	idService        id.Service
	logger           *slog.Logger
	kafkaProd        *kafkaUtil.Producer
	tracer           trace.Tracer
	successPostMeter metric.Int64Counter
	failedPostMeter  metric.Int64Counter
	handledPostMeter metric.Int64Counter
}

func NewHandler(ctx context.Context, repo onerequest.Repository, idService id.Service, kafkaProd *kafkaUtil.Producer) *Handler {
	tracer, _ := otelContext.Tracer(ctx)
	meter, _ := otelContext.Meter(ctx)

	validInt64Counter := func(name string) metric.Int64Counter {
		c, err := meter.Int64Counter(name)
		if err != nil {
			util.Fatal("failed to create counter \"%v\": %v", name, err)
		}
		return c
	}

	return &Handler{
		repo:             repo,
		idService:        idService,
		logger:           slog.Default(),
		kafkaProd:        kafkaProd,
		tracer:           tracer,
		successPostMeter: validInt64Counter(SuccessPostMeterName),
		failedPostMeter:  validInt64Counter(FailedPostMeterName),
		handledPostMeter: validInt64Counter(HandledPostMeterName),
	}
}

// Post is a handler for POST /one
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.handledPostMeter.Add(ctx, 1)
	ctx, span := h.tracer.Start(ctx, "http.handlers.Post", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read request body: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		h.logAndMonitorError(ctx, errMsg, span, err)
		return
	}

	var req model.OneRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse OneRequest json: %s, error: %v", body, err)
		http.Error(w, errMsg, http.StatusBadRequest)
		h.logAndMonitorError(ctx, errMsg, span, err)
		return
	}

	// TODO: kafka message should be sent after the request is saved to the database; should also include
	// the request ID in the message
	err = h.kafkaProd.SendMessage(ctx, body)
	if err != nil {
		errMsg := fmt.Sprintf("failed to send message to kafka, error: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		h.logAndMonitorError(ctx, errMsg, span, err)
		return
	}

	reqID := h.idService.GenerateID()

	err = h.repo.Create(ctx, &onerequest.OneRequest{
		ReqID:  reqID,
		UserID: req.UserID,
		Req:    body,
	})
	if err != nil {
		errMsg := fmt.Sprintf("failed to save request to database, error: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		h.logAndMonitorError(ctx, errMsg, span, err)
		return
	}

	res := Response{
		ReqID: reqID,
	}
	resBody, err := json.Marshal(res)
	if err != nil {
		errMsg := fmt.Sprintf("failed to marshal response, error: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		h.logAndMonitorError(ctx, errMsg, span, err)
		return
	}

	h.successPostMeter.Add(ctx, 1)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resBody)
	span.SetStatus(codes.Ok, "request processed successfully")
}

func (h *Handler) logAndMonitorError(ctx context.Context, errMsg string, span trace.Span, err error) {
	h.logger.ErrorContext(ctx, errMsg)
	h.failedPostMeter.Add(ctx, 1)
	span.SetStatus(codes.Error, errMsg)
	span.RecordError(err)
}
