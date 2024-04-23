package otel

import (
	"context"
	"log/slog"
	"sync"

	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
)

// SlogOTELAttributesHandler is a slog.Handler that adds OTEL attributes to the log record.
type SlogOTELAttributesHandler struct {
	child slog.Handler
	mu    sync.Mutex
}

// NewSlogOTELAttributesHandler creates a new SlogOTELAttributesHandler.
func NewSlogOTELAttributesHandler(child slog.Handler) *SlogOTELAttributesHandler {
	return &SlogOTELAttributesHandler{
		child: child,
	}
}

// Enabled returns true if the log level is enabled.
func (r *SlogOTELAttributesHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return r.child.Enabled(ctx, level)
}

// Handle adds OTEL attributes to the log record and passes it to the child handler.
func (r *SlogOTELAttributesHandler) Handle(ctx context.Context, rec slog.Record) error {
	attr, _ := otelContext.SlogAttributes(ctx)
	if len(attr) > 0 {
		rec.AddAttrs(attr...)
	}
	return r.child.Handle(ctx, rec)
}

// WithAttrs returns a new handler with the given attributes.
func (r *SlogOTELAttributesHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.child = r.child.WithAttrs(attrs)
	return r
}

// WithGroup returns a new handler with the given group name.
func (r *SlogOTELAttributesHandler) WithGroup(name string) slog.Handler {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.child = r.child.WithGroup(name)
	return r
}
