package util

import (
	"context"
	"testing"
	"wn/pkg/constants"

	"github.com/google/uuid"
)

func TestGetStringFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.UserRoleCtx, "ADMIN")

	role, err := GetUserRole(ctx)
	if err != nil {
		t.Fatalf("GetUserRole() error = %v", err)
	}
	if role != "ADMIN" {
		t.Errorf("GetUserRole() = %q, want %q", role, "ADMIN")
	}
}

func TestGetStringFromContext_Missing(t *testing.T) {
	_, err := GetUserRole(context.Background())
	if err == nil {
		t.Fatal("expected error when context value is missing")
	}
}

func TestGetUserId(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	ctx := context.WithValue(context.Background(), constants.UserIdCtx, id.String())

	got, err := GetUserId(ctx)
	if err != nil {
		t.Fatalf("GetUserId() error = %v", err)
	}
	if got != id {
		t.Errorf("GetUserId() = %v, want %v", got, id)
	}
}

func TestGetUserId_InvalidUUID(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.UserIdCtx, "not-a-uuid")

	_, err := GetUserId(ctx)
	if err == nil {
		t.Fatal("expected error for invalid UUID in context")
	}
}

func TestGetRequestId(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.RequestIdCtx, "req-123")

	id, err := GetRequestId(ctx)
	if err != nil {
		t.Fatalf("GetRequestId() error = %v", err)
	}
	if id != "req-123" {
		t.Errorf("GetRequestId() = %q, want %q", id, "req-123")
	}
}

func TestGetTrace(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.TraceIdCtx, "trace-abc")

	id, err := GetTrace(ctx)
	if err != nil {
		t.Fatalf("GetTrace() error = %v", err)
	}
	if id != "trace-abc" {
		t.Errorf("GetTrace() = %q, want %q", id, "trace-abc")
	}
}

func TestGetSpan(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.SpanIdCtx, "span-xyz")

	id, err := GetSpan(ctx)
	if err != nil {
		t.Fatalf("GetSpan() error = %v", err)
	}
	if id != "span-xyz" {
		t.Errorf("GetSpan() = %q, want %q", id, "span-xyz")
	}
}

func TestCopyContextValues(t *testing.T) {
	parent := context.WithValue(context.Background(), constants.UserRoleCtx, "CLIENT")
	parent = context.WithValue(parent, constants.RequestIdCtx, "req-1")

	copied := CopyContextValues(parent, constants.UserRoleCtx, constants.RequestIdCtx, constants.TraceIdCtx)

	role, err := GetUserRole(copied)
	if err != nil {
		t.Fatalf("GetUserRole() on copied ctx error = %v", err)
	}
	if role != "CLIENT" {
		t.Errorf("GetUserRole() = %q, want %q", role, "CLIENT")
	}

	reqID, err := GetRequestId(copied)
	if err != nil {
		t.Fatalf("GetRequestId() on copied ctx error = %v", err)
	}
	if reqID != "req-1" {
		t.Errorf("GetRequestId() = %q, want %q", reqID, "req-1")
	}

	_, err = GetTrace(copied)
	if err == nil {
		t.Error("expected error for key not present in parent context")
	}
}

func TestGetUUIDFromContext(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	ctx := context.WithValue(context.Background(), constants.UserIdCtx, id.String())

	got, err := GetUUIDFromContext(ctx, constants.UserIdCtx)
	if err != nil {
		t.Fatalf("GetUUIDFromContext() error = %v", err)
	}
	if *got != id {
		t.Errorf("GetUUIDFromContext() = %v, want %v", *got, id)
	}
}
