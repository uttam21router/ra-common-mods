// Package ctxmeta standardizes correlation metadata in context.Context.
//
// It is dependency-free and is meant to be used by:
// - server middleware (to set request/tenant/user IDs)
// - clients (to read IDs for propagation)
// - logging modules (to enrich logs with correlation fields)
package ctxmeta
