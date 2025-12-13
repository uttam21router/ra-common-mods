package api

// Backend is the swappable backend interface used by the logging core.
type Backend interface {
    Log(level Level, msg string, fields Fields)
    With(fields Fields) Backend
}
