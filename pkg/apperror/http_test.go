package apperror

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetErrorByHttpStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		wantType   ErrType
		wantCode   string
		wantMsg    string
	}{
		{"bad request", http.StatusBadRequest, BadRequestError, "code1", "msg1"},
		{"unauthorized", http.StatusUnauthorized, UnauthorizedError, "code2", "msg2"},
		{"forbidden", http.StatusForbidden, AccessDeniedError, "code3", "msg3"},
		{"conflict", http.StatusConflict, ConflictError, "code4", "msg4"},
		{"gone", http.StatusGone, NotFoundError, "code5", "msg5"},
		{"unprocessable", http.StatusUnprocessableEntity, InvalidDataError, "code6", "msg6"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetErrorByHttpStatus(tt.status, tt.wantMsg, tt.wantCode)
			appErr, ok := err.(*AppError)
			if !ok {
				t.Fatalf("expected *AppError, got %T", err)
			}
			if appErr.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", appErr.Type, tt.wantType)
			}
			if appErr.Message != tt.wantMsg {
				t.Errorf("Message = %q, want %q", appErr.Message, tt.wantMsg)
			}
			if appErr.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", appErr.Code, tt.wantCode)
			}
		})
	}

	t.Run("default internal error", func(t *testing.T) {
		err := GetErrorByHttpStatus(http.StatusTeapot, "oops", "code")
		appErr, ok := err.(*AppError)
		if !ok {
			t.Fatalf("expected *AppError, got %T", err)
		}
		if appErr.Type != InternalServerError {
			t.Errorf("Type = %q, want %q", appErr.Type, InternalServerError)
		}
		if !errors.Is(appErr.Err, errors.New("oops")) {
			// errors.New creates new instance each time; compare message instead
			if appErr.Err.Error() != "oops" {
				t.Errorf("Err message = %q, want %q", appErr.Err.Error(), "oops")
			}
		}
	})
}

func TestGetHttpStatusByErrorType(t *testing.T) {
	tests := []struct {
		errType ErrType
		want    int
	}{
		{AccessDeniedError, http.StatusForbidden},
		{InvalidDataError, http.StatusUnprocessableEntity},
		{ConflictError, http.StatusConflict},
		{NotFoundError, http.StatusGone},
		{UnauthorizedError, http.StatusUnauthorized},
		{BadRequestError, http.StatusBadRequest},
		{InternalServerError, http.StatusInternalServerError},
		{ErrType("Unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(string(tt.errType), func(t *testing.T) {
			if got := GetHttpStatusByErrorType(tt.errType); got != tt.want {
				t.Errorf("GetHttpStatusByErrorType(%q) = %d, want %d", tt.errType, got, tt.want)
			}
		})
	}
}
