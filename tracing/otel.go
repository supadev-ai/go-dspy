package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer provides tracing capabilities for DSPy operations.
type Tracer struct {
	tracer trace.Tracer
}

// NewTracer creates a new tracer with the given name.
func NewTracer(name string) *Tracer {
	return &Tracer{
		tracer: otel.Tracer(name),
	}
}

// StartSpan starts a new span with the given name.
func (t *Tracer) StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name)
}

// StartSpanWithOptions starts a new span with additional options.
func (t *Tracer) StartSpanWithOptions(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}
