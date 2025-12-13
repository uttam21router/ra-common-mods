package api

type Level string

const (
	LevelDebug    Level = "debug"
	LevelInfo     Level = "info"
	LevelWarn     Level = "warn"
	LevelError    Level = "error"
	LevelDisabled Level = "disabled"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

type StackTraceMode string

const (
	StackTraceNone  StackTraceMode = "none"
	StackTraceError StackTraceMode = "error"
	StackTraceWarn  StackTraceMode = "warn"
	StackTraceAll   StackTraceMode = "all"
)

type Config struct {
	ServiceName string
	Environment string
	Level       Level
	Format      Format
	EnableOTel  bool
	StackTrace  StackTraceMode

	SubsystemDefault Level
	SubsystemLevels  map[string]Level
}

func (c Config) WithService(name string) Config {
	c.ServiceName = name
	return c
}
