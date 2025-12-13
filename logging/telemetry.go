package logging

import (
    "context"
    "fmt"

    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

func extractOTelFields(ctx context.Context) Fields {
    if ctx == nil {
        return Fields{}
    }
    sc := trace.SpanContextFromContext(ctx)
    if !sc.IsValid() {
        return Fields{}
    }
    return Fields{
        "trace_id": sc.TraceID().String(),
        "span_id":  sc.SpanID().String(),
    }
}

func addSpanEvent(ctx context.Context, eventName string, attrs Fields) {
    if ctx == nil {
        return
    }
    span := trace.SpanFromContext(ctx)
    if span == nil {
        return
    }
    kvs := make([]attribute.KeyValue, 0, len(attrs))
    for k, v := range attrs {
        kvs = append(kvs, attribute.String(k, fmt.Sprint(v)))
    }
    span.AddEvent(eventName, trace.WithAttributes(kvs...))
}
