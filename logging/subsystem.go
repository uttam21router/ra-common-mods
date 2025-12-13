package logging

import "sync"

type subsystemRegistry struct {
	mu           sync.RWMutex
	defaultLevel Level
	levels       map[string]Level
}

func newSubsystemRegistry(defaultLevel Level, overrides map[string]Level) *subsystemRegistry {
	cp := map[string]Level{}
	for k, v := range overrides {
		cp[k] = v
	}
	return &subsystemRegistry{defaultLevel: defaultLevel, levels: cp}
}

func (r *subsystemRegistry) setLevel(name string, lvl Level) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.levels[name] = lvl
}

func (r *subsystemRegistry) levelFor(subsystem string, global Level) Level {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if subsystem == "" {
		return global
	}
	if lvl, ok := r.levels[subsystem]; ok {
		return lvl
	}
	if r.defaultLevel != "" {
		return r.defaultLevel
	}
	return global
}
