package api

import "context"

// Fields represents structured log fields.
type Fields map[string]any

// Logger is the stable interface used by microservices.
type Logger interface {
	Debug(ctx context.Context, msg string, fields Fields)
	Info(ctx context.Context, msg string, fields Fields)
	Warn(ctx context.Context, msg string, fields Fields)
	Error(ctx context.Context, msg string, fields Fields)

	With(fields Fields) Logger
	Subsystem(name string) Logger
}
