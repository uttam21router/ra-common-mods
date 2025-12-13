package logging

import "runtime/debug"

const (
    _errKey          = "__logging_err"
    _captureStackKey = "__logging_capture_stack"
)

// Err marks an error for logging.
// The logger will output "error" and may output "stack_trace" based on config.
func Err(err error) Fields {
    if err == nil {
        return Fields{}
    }
    return Fields{
        _errKey:          err,
        _captureStackKey: true,
    }
}

func errorToFields(in Fields, mode StackTraceMode, level Level) Fields {
    if in == nil {
        return Fields{}
    }
    out := Fields{}
    for k, v := range in {
        if k == _errKey || k == _captureStackKey {
            continue
        }
        out[k] = v
    }

    if _, ok := out["error"]; ok {
        return out
    }

    raw, ok := in[_errKey]
    if !ok {
        return out
    }
    err, ok := raw.(error)
    if !ok || err == nil {
        return out
    }

    out["error"] = err.Error()

    if cap, _ := in[_captureStackKey].(bool); cap && stackEnabled(mode, level) {
        out["stack_trace"] = string(debug.Stack())
    }
    return out
}

func stackEnabled(mode StackTraceMode, level Level) bool {
    switch mode {
    case StackTraceAll:
        return true
    case StackTraceWarn:
        return levelRank(level) >= levelRank(LevelWarn)
    case StackTraceError:
        return levelRank(level) >= levelRank(LevelError)
    default:
        return false
    }
}
