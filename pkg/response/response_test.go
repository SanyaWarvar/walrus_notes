package response

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"wn/pkg/constants"
)

func ctxWithRequestID(id string) context.Context {
	return context.WithValue(context.Background(), constants.RequestIdCtx, id)
}

func TestBuildSuccessResponse(t *testing.T) {
	builder := NewResponseBuilder(false)
	status, resp := builder.BuildSuccessResponse(ctxWithRequestID("req-1"))

	if status != http.StatusOK {
		t.Errorf("status = %d, want %d", status, http.StatusOK)
	}
	if resp.Meta.Message != "OK" {
		t.Errorf("Meta.Message = %q, want %q", resp.Meta.Message, "OK")
	}
	if resp.Meta.RequestId != "req-1" {
		t.Errorf("Meta.RequestId = %q, want %q", resp.Meta.RequestId, "req-1")
	}
	if resp.Data != nil {
		t.Errorf("Data = %v, want nil", resp.Data)
	}
}

func TestBuildSuccessResponseBody(t *testing.T) {
	builder := NewResponseBuilder(false)
	data := map[string]string{"key": "value"}

	status, resp := builder.BuildSuccessResponseBody(ctxWithRequestID("req-2"), data)

	if status != http.StatusOK {
		t.Errorf("status = %d, want %d", status, http.StatusOK)
	}
	if resp.Data == nil {
		t.Fatal("Data should not be nil")
	}
	payload, ok := resp.Data.(map[string]string)
	if !ok {
		t.Fatalf("Data type = %T, want map[string]string", resp.Data)
	}
	if payload["key"] != "value" {
		t.Errorf("Data[key] = %q, want %q", payload["key"], "value")
	}
}

func TestBuildSuccessPaginationResponse(t *testing.T) {
	builder := NewResponseBuilder(false)
	items := []int{1, 2, 3}

	resp := builder.BuildSuccessPaginationResponse(ctxWithRequestID("req-3"), 2, 50, 5, items)

	if resp.Pagination.Page != 2 {
		t.Errorf("Pagination.Page = %d, want 2", resp.Pagination.Page)
	}
	if resp.Pagination.PerPage != 50 {
		t.Errorf("Pagination.PerPage = %d, want 50", resp.Pagination.PerPage)
	}
	if resp.Pagination.Pages != 5 {
		t.Errorf("Pagination.Pages = %d, want 5", resp.Pagination.Pages)
	}
}

func TestBuildErrorResponse(t *testing.T) {
	t.Run("without error export", func(t *testing.T) {
		builder := NewResponseBuilder(false)
		err := errors.New("internal details")

		resp := builder.BuildErrorResponse(ctxWithRequestID("req-4"), "something failed", "err_code", err)

		if resp.Meta.Message != "something failed" {
			t.Errorf("Meta.Message = %q, want %q", resp.Meta.Message, "something failed")
		}
		if resp.Meta.Code != "err_code" {
			t.Errorf("Meta.Code = %q, want %q", resp.Meta.Code, "err_code")
		}
		if resp.Meta.Error != "" {
			t.Errorf("Meta.Error = %q, want empty when ErrorExport is false", resp.Meta.Error)
		}
	})

	t.Run("with error export", func(t *testing.T) {
		builder := NewResponseBuilder(true)
		err := errors.New("internal details")

		resp := builder.BuildErrorResponse(context.Background(), "fail", "code", err)

		if resp.Meta.Error != "internal details" {
			t.Errorf("Meta.Error = %q, want %q", resp.Meta.Error, "internal details")
		}
	})
}

func TestBuildResponse_WithoutRequestID(t *testing.T) {
	builder := NewResponseBuilder(false)
	_, resp := builder.BuildSuccessResponse(context.Background())

	if resp.Meta.RequestId != "" {
		t.Errorf("Meta.RequestId = %q, want empty", resp.Meta.RequestId)
	}
}
