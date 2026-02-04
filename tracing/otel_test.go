package tracing

import (
	"context"
	"testing"
)

func TestNewTracer(t *testing.T) {
	tracer := NewTracer("test-tracer")
	if tracer == nil {
		t.Fatal("Expected tracer, got nil")
	}

	if tracer.tracer == nil {
		t.Error("Expected internal tracer to be set, got nil")
	}
}

func TestTracer_StartSpan(t *testing.T) {
	tracer := NewTracer("test-tracer")
	ctx := context.Background()

	newCtx, span := tracer.StartSpan(ctx, "test-span")
	if newCtx == nil {
		t.Error("Expected context, got nil")
	}

	if span == nil {
		t.Error("Expected span, got nil")
	}

	span.End()
}

func TestTracer_StartSpanWithOptions(t *testing.T) {
	tracer := NewTracer("test-tracer")
	ctx := context.Background()

	newCtx, span := tracer.StartSpanWithOptions(ctx, "test-span")
	if newCtx == nil {
		t.Error("Expected context, got nil")
	}

	if span == nil {
		t.Error("Expected span, got nil")
	}

	span.End()
}
