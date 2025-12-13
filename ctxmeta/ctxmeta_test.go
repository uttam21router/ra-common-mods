package ctxmeta

import (
	"context"
	"reflect"
	"testing"
)

func TestSetGet(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "r1")
	ctx = WithTenantID(ctx, "t1")
	ctx = WithUserID(ctx, "u1")

	if v, ok := RequestID(ctx); !ok || v != "r1" {
		t.Fatalf("request_id mismatch: %v %v", v, ok)
	}
	if v, ok := TenantID(ctx); !ok || v != "t1" {
		t.Fatalf("tenant_id mismatch: %v %v", v, ok)
	}
	if v, ok := UserID(ctx); !ok || v != "u1" {
		t.Fatalf("user_id mismatch: %v %v", v, ok)
	}
}

func TestFields(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "r1")
	ctx = WithTenantID(ctx, "t1")
	ctx = WithUserID(ctx, "u1")

	got := Fields(ctx)
	want := map[string]any{
		"request_id": "r1",
		"tenant_id":  "t1",
		"user_id":    "u1",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fields mismatch  got: %#v want: %#v", got, want)
	}
}

func TestClone(t *testing.T) {
	src := context.Background()
	src = WithRequestID(src, "r1")
	src = WithTenantID(src, "t1")

	dst := Clone(context.Background(), src)
	if v, ok := RequestID(dst); !ok || v != "r1" {
		t.Fatalf("clone request_id mismatch: %v %v", v, ok)
	}
	if v, ok := TenantID(dst); !ok || v != "t1" {
		t.Fatalf("clone tenant_id mismatch: %v %v", v, ok)
	}
}
