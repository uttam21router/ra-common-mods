package logging

import (
    "path/filepath"
    "runtime"
    "strings"
)

type callerInfo struct {
    file string
    line int
}

func findCaller() callerInfo {
    const maxDepth = 32
    pcs := make([]uintptr, maxDepth)

    // Skip runtime.Callers + findCaller + emit
    n := runtime.Callers(3, pcs)
    frames := runtime.CallersFrames(pcs[:n])

    for {
        fr, more := frames.Next()
        if fr.File == "" {
            if !more {
                break
            }
            continue
        }
        if isInternalFrame(fr.Function, fr.File) {
            if !more {
                break
            }
            continue
        }
        return callerInfo{
            file: filepath.Base(fr.File),
            line: fr.Line,
        }
    }
    return callerInfo{}
}

func isInternalFrame(fn, file string) bool {
    lfn := strings.ToLower(fn)
    lf := strings.ToLower(file)

    // Skip frames from this module and its adapters.
    if strings.Contains(lfn, "ra-common-mods/logging") {
        return true
    }
    if strings.Contains(lf, string(filepath.Separator)+"adapters"+string(filepath.Separator)) &&
        strings.Contains(lf, string(filepath.Separator)+"logging"+string(filepath.Separator)) {
        return true
    }
    return false
}
