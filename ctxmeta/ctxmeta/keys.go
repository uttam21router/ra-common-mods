package ctxmeta

// Unexported key type to avoid collisions across packages.
type key int

const (
	keyRequestID key = iota
	keyTenantID
	keyUserID
)
