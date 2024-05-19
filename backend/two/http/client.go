package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/kartpop/cruncan/backend/pkg/model"
	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	DefaultTimeout               = 10 * time.Second
	ContentType           string = "Content-Type"
	ApplicationJson       string = "application/json"
	ApplicationUrlEncoded string = "application/x-www-form-urlencoded"
	threeRequestPath      string = "/three"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
	logger     *slog.Logger
}

func NewClient(httpClient *http.Client, baseUrl string, logger *slog.Logger) *Client {
	httpClient.Timeout = DefaultTimeout
	return &Client{
		httpClient: httpClient,
		baseUrl:    baseUrl,
		logger:     logger,
	}
}

func (c *Client) PostThreeRequest(ctx context.Context, threeReq *model.ThreeRequest) (*http.Response, error) {
	tracer, _ := otelContext.Tracer(ctx)
	var span trace.Span
	ctx, span = tracer.Start(ctx, "http.Client.PostThreeRequest", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	reqJson, err := json.Marshal(threeReq)
	if err != nil {
		errMssg := fmt.Sprintf("failed to marshal three request: %v", err)
		span.SetStatus(codes.Error, errMssg)
		span.RecordError(err)
		return nil, fmt.Errorf(errMssg)
	}

	span.SetAttributes(
		attribute.KeyValue{
			Key:   "reqbody",
			Value: attribute.StringValue(string(reqJson)),
		})

	req, err := http.NewRequest(http.MethodPost, c.baseUrl, bytes.NewReader(reqJson))
	if err != nil {
		errMssg := fmt.Sprintf("failed to create request: %v", err)
		span.SetStatus(codes.Error, errMssg)
		span.RecordError(err)
		return nil, fmt.Errorf(errMssg)
	}

	req.Header.Add(ContentType, ApplicationJson)
	req.Header.Add("language", "en")
	req.URL.Path += threeRequestPath
	req = req.WithContext(ctx)

	// Dump the request for debugging
	dump, _ := httputil.DumpRequestOut(req, true)
	c.logger.Debug("three_req", "dump", string(dump))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		errMssg := fmt.Sprintf("failed to send request: %v", err)
		span.SetStatus(codes.Error, errMssg)
		span.RecordError(err)
		return nil, fmt.Errorf(errMssg)
	}

	span.SetStatus(codes.Ok, "request sent")
	return resp, nil
}
