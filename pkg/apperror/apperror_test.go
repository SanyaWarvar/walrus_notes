package apperror

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	inner := errors.New("inner error")
	err := NewAppError(inner, "user message")

	if err.Error() != "inner error" {
		t.Errorf("Error() = %q, want %q", err.Error(), "inner error")
	}
}

func TestAppError_Unwrap(t *testing.T) {
	inner := errors.New("inner error")
	err := NewAppError(inner, "user message")

	if !errors.Is(err, inner) {
		t.Error("Unwrap() should return inner error")
	}
}

func TestAppError_WithCode(t *testing.T) {
	err := NewBadRequestError("bad input", "bad_code").WithCode("override")
	if err.Code != "override" {
		t.Errorf("WithCode() = %q, want %q", err.Code, "override")
	}
}

func TestAppError_Marshal(t *testing.T) {
	err := NewNotFoundError("not found", "not_found")
	data := err.Marshal()

	var decoded map[string]any
	if unmarshalErr := json.Unmarshal(data, &decoded); unmarshalErr != nil {
		t.Fatalf("json.Unmarshal() error = %v", unmarshalErr)
	}

	if decoded["type"] != string(NotFoundError) {
		t.Errorf("type = %v, want %q", decoded["type"], NotFoundError)
	}
	if decoded["message"] != "not found" {
		t.Errorf("message = %v, want %q", decoded["message"], "not found")
	}
	if decoded["code"] != "not_found" {
		t.Errorf("code = %v, want %q", decoded["code"], "not_found")
	}
}

func TestNewInternalError(t *testing.T) {
	inner := errors.New("db failure")
	err := NewInternalError(inner)

	if err.Type != InternalServerError {
		t.Errorf("Type = %q, want %q", err.Type, InternalServerError)
	}
	if err.Message != unexpectedErrorMessage {
		t.Errorf("Message = %q, want %q", err.Message, unexpectedErrorMessage)
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name    string
		err     *AppError
		errType ErrType
		message string
		code    string
	}{
		{"bad request", NewBadRequestError("bad", "bad_code"), BadRequestError, "bad", "bad_code"},
		{"not found", NewNotFoundError("missing", "nf"), NotFoundError, "missing", "nf"},
		{"access denied", NewAccessDeniedError("denied", "ad"), AccessDeniedError, "denied", "ad"},
		{"conflict", NewConflictError("dup", "cf"), ConflictError, "dup", "cf"},
		{"invalid data", NewInvalidDataError("invalid", "iv"), InvalidDataError, "invalid", "iv"},
		{"unauthorized", NewUnauthorizedError("no auth", "ua"), UnauthorizedError, "no auth", "ua"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Type != tt.errType {
				t.Errorf("Type = %q, want %q", tt.err.Type, tt.errType)
			}
			if tt.err.Message != tt.message {
				t.Errorf("Message = %q, want %q", tt.err.Message, tt.message)
			}
			if tt.err.Code != tt.code {
				t.Errorf("Code = %q, want %q", tt.err.Code, tt.code)
			}
		})
	}
}
