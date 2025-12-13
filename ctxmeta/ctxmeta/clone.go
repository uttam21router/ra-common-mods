package ctxmeta

import "context"

// Clone copies all known ctxmeta values from src into dst.
// Useful when spawning goroutines with derived contexts.
func Clone(dst, src context.Context) context.Context {
	if dst == nil {
		dst = context.Background()
	}
	if v, ok := RequestID(src); ok {
		dst = WithRequestID(dst, v)
	}
	if v, ok := TenantID(src); ok {
		dst = WithTenantID(dst, v)
	}
	if v, ok := UserID(src); ok {
		dst = WithUserID(dst, v)
	}
	return dst
}
