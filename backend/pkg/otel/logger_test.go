package otel

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestInitLogger_AddOtelAttributes(t *testing.T) {
	err := os.Setenv("OTEL_SERVICE_NAME", "test-service")
	if err != nil {
		t.Errorf("failed to set OTEL_SERVICE_NAME: %v", err)
	}
	err = os.Setenv("OTEL_EXTRA_ATTRIBUTES", "key1=value1,key2=value2")
	if err != nil {
		t.Errorf("failed to set OTEL_EXTRA_ATTRIBUTES: %v", err)
	}
	b := &bytes.Buffer{}
	ctx, cancel, err := InitLogger(context.Background(), WithTestBuffer(b), WithConsoleOnly())
	defer cancel()
	if err != nil {
		t.Errorf("failed to initialize logger: %v", err)
	}
	//{"time":"2024-03-15T17:03:29.517115+01:00","level":"INFO","msg":"test message","service.name":"test-service","key1":"value1","key2":"value2"}
	slog.Default().InfoContext(ctx, "test message")
	if !strings.Contains(b.String(), `"service.name":"test-service"`) {
		t.Errorf("expected service.name to be test-service, got %s", b.String())
	}
	if !strings.Contains(b.String(), `"key1":"value1"`) {
		t.Errorf("expected key1 to be value1, got %s", b.String())
	}
	if !strings.Contains(b.String(), `"key2":"value2"`) {
		t.Errorf("expected key2 to be value2, got %s", b.String())
	}
}
