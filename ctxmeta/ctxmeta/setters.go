package ctxmeta

import "context"

// WithRequestID returns a derived context with request_id set.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, keyRequestID, requestID)
}

// WithTenantID returns a derived context with tenant_id set.
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, keyTenantID, tenantID)
}

// WithUserID returns a derived context with user_id set.
func WithUserID(ctx context.Context, userID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, keyUserID, userID)
}
