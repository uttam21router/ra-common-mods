package logging

func mergeFields(a, b Fields) Fields {
    if a == nil && b == nil {
        return Fields{}
    }
    out := Fields{}
    for k, v := range a {
        out[k] = v
    }
    for k, v := range b {
        out[k] = v
    }
    return out
}

func levelRank(l Level) int {
    switch l {
    case LevelDebug:
        return 10
    case LevelInfo:
        return 20
    case LevelWarn:
        return 30
    case LevelError:
        return 40
    case LevelDisabled:
        return 100
    default:
        return 20
    }
}
