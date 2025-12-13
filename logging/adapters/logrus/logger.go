package logrus

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/routerarchitects/ra-common-mods/logging/api"
)

type Options struct {
	Format api.Format
	Level  api.Level
	Output io.Writer
}

type adapter struct {
	base *logrus.Entry
}

func New(opt Options) api.Backend {
	l := logrus.New()
	l.SetOutput(opt.Output)
	l.SetLevel(mapLevel(opt.Level))

	switch opt.Format {
	case api.FormatText:
		l.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	default:
		l.SetFormatter(&logrus.JSONFormatter{})
	}

	return &adapter{base: logrus.NewEntry(l)}
}

func (a *adapter) Log(level api.Level, msg string, fields api.Fields) {
	e := a.base.WithFields(logrus.Fields(fields))
	switch strings.ToLower(string(level)) {
	case "debug":
		e.Debug(msg)
	case "warn":
		e.Warn(msg)
	case "error":
		e.Error(msg)
	default:
		e.Info(msg)
	}
}

func (a *adapter) With(fields api.Fields) api.Backend {
	return &adapter{base: a.base.WithFields(logrus.Fields(fields))}
}

func mapLevel(l api.Level) logrus.Level {
	switch strings.ToLower(string(l)) {
	case "debug":
		return logrus.DebugLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

var _ api.Backend = (*adapter)(nil)
