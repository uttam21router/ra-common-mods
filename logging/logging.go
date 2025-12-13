package logging

import "github.com/routerarchitects/ra-common-mods/logging/api"

// Re-exports (SemVer-stable public API)
type Logger = api.Logger
type Fields = api.Fields
type Level = api.Level
type Format = api.Format
type StackTraceMode = api.StackTraceMode
type Config = api.Config
type Backend = api.Backend

const (
	LevelDebug    = api.LevelDebug
	LevelInfo     = api.LevelInfo
	LevelWarn     = api.LevelWarn
	LevelError    = api.LevelError
	LevelDisabled = api.LevelDisabled

	FormatJSON = api.FormatJSON
	FormatText = api.FormatText

	StackTraceNone  = api.StackTraceNone
	StackTraceError = api.StackTraceError
	StackTraceWarn  = api.StackTraceWarn
	StackTraceAll   = api.StackTraceAll
)
