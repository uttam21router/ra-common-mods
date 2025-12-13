package ctxmeta

import "context"

func RequestID(ctx context.Context) (string, bool) { return getString(ctx, keyRequestID) }
func TenantID(ctx context.Context) (string, bool)  { return getString(ctx, keyTenantID) }
func UserID(ctx context.Context) (string, bool)    { return getString(ctx, keyUserID) }

func getString(ctx context.Context, k key) (string, bool) {
	if ctx == nil {
		return "", false
	}
	v := ctx.Value(k)
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

func Fields(ctx context.Context) map[string]any {
	out := map[string]any{}
	if v, ok := RequestID(ctx); ok {
		out["request_id"] = v
	}
	if v, ok := TenantID(ctx); ok {
		out["tenant_id"] = v
	}
	if v, ok := UserID(ctx); ok {
		out["user_id"] = v
	}
	return out
}
