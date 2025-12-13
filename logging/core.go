package logging

import (
	"context"
	"os"
	"sync/atomic"

	logrusadapter "github.com/routerarchitects/ra-common-mods/logging/adapters/logrus"
)

type coreLogger struct {
	backend   Backend
	static    Fields
	subsystem string
}

func (l *coreLogger) Debug(ctx context.Context, msg string, fields Fields) {
	emit(ctx, l, LevelDebug, msg, fields)
}
func (l *coreLogger) Info(ctx context.Context, msg string, fields Fields) {
	emit(ctx, l, LevelInfo, msg, fields)
}
func (l *coreLogger) Warn(ctx context.Context, msg string, fields Fields) {
	emit(ctx, l, LevelWarn, msg, fields)
}
func (l *coreLogger) Error(ctx context.Context, msg string, fields Fields) {
	emit(ctx, l, LevelError, msg, fields)
}

func (l *coreLogger) With(fields Fields) Logger {
	return &coreLogger{
		backend:   l.backend.With(fields),
		static:    mergeFields(l.static, fields),
		subsystem: l.subsystem,
	}
}

func (l *coreLogger) Subsystem(name string) Logger {
	sf := Fields{"subsystem": name}
	return &coreLogger{
		backend:   l.backend.With(sf),
		static:    mergeFields(l.static, sf),
		subsystem: name,
	}
}

var (
	gLogger atomic.Value // Logger
	gCfg    atomic.Value // Config
	gReg    atomic.Value // *subsystemRegistry
)

func Init(cfg Config) {
	cfg = applyDefaults(cfg)
	gCfg.Store(cfg)
	gReg.Store(newSubsystemRegistry(cfg.SubsystemDefault, cfg.SubsystemLevels))

	// Default backend: logrus
	b := logrusadapter.New(logrusadapter.Options{
		Format: cfg.Format,
		Level:  cfg.Level,
		Output: os.Stdout,
	})

	base := Fields{
		"service": cfg.ServiceName,
		"env":     cfg.Environment,
	}

	gl := &coreLogger{
		backend:   b.With(base),
		static:    base,
		subsystem: "",
	}
	gLogger.Store(Logger(gl))
}

func InitService(serviceName string) {
	Init(ConfigFromEnv().WithService(serviceName))
}

func L() Logger {
	if v := gLogger.Load(); v != nil {
		return v.(Logger)
	}
	Init(ConfigFromEnv())
	return gLogger.Load().(Logger)
}

func Subsystem(name string) Logger { return L().Subsystem(name) }

func SetSubsystemLevel(name string, lvl Level) {
	reg := getRegistry()
	reg.setLevel(name, lvl)
}

func applyDefaults(cfg Config) Config {
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown-service"
	}
	if cfg.Environment == "" {
		cfg.Environment = "unknown"
	}
	if cfg.Level == "" {
		cfg.Level = LevelInfo
	}
	if cfg.Format == "" {
		cfg.Format = FormatJSON
	}
	if cfg.StackTrace == "" {
		cfg.StackTrace = StackTraceNone
	}
	if cfg.SubsystemDefault == "" {
		cfg.SubsystemDefault = LevelError
	}
	if cfg.SubsystemLevels == nil {
		cfg.SubsystemLevels = map[string]Level{}
	}
	return cfg
}

func getCfg() Config {
	if v := gCfg.Load(); v != nil {
		return v.(Config)
	}
	cfg := applyDefaults(ConfigFromEnv())
	gCfg.Store(cfg)
	return cfg
}

func getRegistry() *subsystemRegistry {
	if v := gReg.Load(); v != nil {
		return v.(*subsystemRegistry)
	}
	cfg := getCfg()
	r := newSubsystemRegistry(cfg.SubsystemDefault, cfg.SubsystemLevels)
	gReg.Store(r)
	return r
}

func emit(ctx context.Context, l *coreLogger, level Level, msg string, fields Fields) {
	cfg := getCfg()

	thr := getRegistry().levelFor(l.subsystem, cfg.Level)
	if thr == LevelDisabled {
		return
	}
	if levelRank(level) < levelRank(thr) {
		return
	}

	out := mergeFields(l.static, extractContextFields(ctx))

	if cfg.EnableOTel {
		out = mergeFields(out, extractOTelFields(ctx))
	}

	c := findCaller()
	if c.file != "" {
		out["caller_file"] = c.file
		out["caller_line"] = c.line
	}

	out = mergeFields(out, fields)
	out = errorToFields(out, cfg.StackTrace, level)

	l.backend.Log(level, msg, out)

	if cfg.EnableOTel && (level == LevelWarn || level == LevelError) {
		attrs := Fields{
			"message":   msg,
			"level":     string(level),
			"subsystem": l.subsystem,
		}
		if v, ok := out["error"]; ok {
			attrs["error"] = v
		}
		addSpanEvent(ctx, "log."+string(level), attrs)
	}
}

var _ Logger = (*coreLogger)(nil)
