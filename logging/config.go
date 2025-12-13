package logging

import (
	"os"
	"strings"
)

// ConfigFromEnv reads logger configuration from environment variables.
func ConfigFromEnv() Config {
	return Config{
		ServiceName:      getenv("SERVICE_NAME", "unknown-service"),
		Environment:      getenv("LOG_ENV", "unknown"),
		Level:            parseLevel(getenv("LOG_LEVEL", "info")),
		Format:           parseFormat(getenv("LOG_FORMAT", "json")),
		EnableOTel:       parseBool(getenv("LOG_ENABLE_OTEL", "false")),
		StackTrace:       parseStackTraceMode(getenv("LOG_STACKTRACE", "none")),
		SubsystemDefault: parseLevel(getenv("LOG_SUBSYSTEM_DEFAULT", "error")),
		SubsystemLevels:  parseSubsystemLevels(getenv("LOG_SUBSYSTEM_LEVELS", "")),
	}
}

func getenv(k, d string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return d
}

func parseBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "t", "true", "y", "yes", "on":
		return true
	default:
		return false
	}
}

func parseLevel(v string) Level {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return LevelDebug
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	case "disabled", "off", "none":
		return LevelDisabled
	default:
		return LevelInfo
	}
}

func parseFormat(v string) Format {
	if strings.ToLower(strings.TrimSpace(v)) == "text" {
		return FormatText
	}
	return FormatJSON
}

func parseStackTraceMode(v string) StackTraceMode {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "all":
		return StackTraceAll
	case "warn":
		return StackTraceWarn
	case "error":
		return StackTraceError
	default:
		return StackTraceNone
	}
}

func parseSubsystemLevels(v string) map[string]Level {
	out := map[string]Level{}
	v = strings.TrimSpace(v)
	if v == "" {
		return out
	}
	for _, part := range strings.Split(v, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		name := strings.TrimSpace(kv[0])
		lvl := parseLevel(kv[1])
		if name != "" {
			out[name] = lvl
		}
	}
	return out
}
