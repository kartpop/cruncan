package otel

import (
	"context"
	"fmt"
	"log/slog"

	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetAttributesOnSpanAndContext sets the attributes on the otel span and context for slog
func SetAttributesOnSpanAndContext(ctx context.Context, span trace.Span, attributes ...attribute.KeyValue) context.Context {
	if len(attributes) == 0 {
		return ctx
	}
	if span != nil {
		span.SetAttributes(attributes...)
	}
	if ctx != nil {
		slAttr := make([]slog.Attr, 0, len(attributes))
		for _, attr := range attributes {
			switch attr.Value.Type() {
			case attribute.STRING:
				slAttr = append(slAttr, slog.String(string(attr.Key), attr.Value.AsString()))
			case attribute.BOOL:
				slAttr = append(slAttr, slog.Bool(string(attr.Key), attr.Value.AsBool()))
			case attribute.INT64:
				slAttr = append(slAttr, slog.Int64(string(attr.Key), attr.Value.AsInt64()))
			case attribute.FLOAT64:
				slAttr = append(slAttr, slog.Float64(string(attr.Key), attr.Value.AsFloat64()))
			default:
				slAttr = append(slAttr, slog.Any(string(attr.Key), attr.Value.AsInterface()))
			}
		}
		ctx = otelContext.AddSlogAttributes(ctx, slAttr...)
	}
	return ctx
}

// SetSpanOk sets the span status to Ok
func SetSpanOk(span trace.Span) {
	if span == nil {
		return
	}
	span.SetStatus(codes.Ok, "success")
}

// SetSpanOkWithMessage sets the span status to Ok with a message
func SetSpanOkWithMessage(span trace.Span, message string, args ...any) {
	if span == nil {
		return
	}
	if len(args) == 0 {
		span.SetStatus(codes.Ok, message)
	} else {
		span.SetStatus(codes.Ok, fmt.Sprintf(message, args...))
	}
}

// SetSpanErrorMessage sets the span status to Error with a message
func SetSpanErrorMessage(span trace.Span, message string, args ...any) {
	if span == nil {
		return
	}
	if len(args) == 0 {
		span.SetStatus(codes.Error, message)
	} else {
		span.SetStatus(codes.Error, fmt.Sprintf(message, args...))
	}
}

// SetSpanErrorWithMessage sets the span status to Error with a message and an error
func SetSpanErrorWithMessage(span trace.Span, err error, message string, args ...any) {
	if span == nil || err == nil {
		return
	}
	if len(args) == 0 {
		span.SetStatus(codes.Error, message)
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Error, fmt.Sprintf(message, args...))
		span.RecordError(err)
	}

}

// SetAutoSpanStatus sets the span status to Error if it is not nil, otherwise Ok
func SetAutoSpanStatus(span trace.Span, err error) {
	if span == nil {
		return
	}
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return
	}
	span.SetStatus(codes.Ok, "success")
}

// SetSpanError sets the span status to Error with an error
func SetSpanError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

// UnsetSpan sets the span status to Unset
func UnsetSpan(span trace.Span) {
	if span == nil {
		return
	}
	span.SetStatus(codes.Unset, "")
}
