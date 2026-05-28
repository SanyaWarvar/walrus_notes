package apperrors

import (
	"testing"
	"wn/pkg/apperror"
)

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *apperror.AppError
		errType  apperror.ErrType
		code     string
		contains string
	}{
		{"invalid auth header", InvalidAuthorizationHeader, apperror.UnauthorizedError, "invalid_authorization_header", "authorization"},
		{"invalid token", InvalidTokenError, apperror.UnauthorizedError, "invalid_token", "token"},
		{"user not found", UserNotFound, apperror.InvalidDataError, "user_not_found", "user"},
		{"incorrect password", IncorrectPassword, apperror.UnauthorizedError, "incorrect_password", "password"},
		{"note not found", NoteNotFound, apperror.InvalidDataError, "note_not_found", "note"},
		{"permissions not enough", PermissionsNotEnough, apperror.InvalidDataError, "premissions_not_enough", "permissions"},
		{"no new password", NoNewPassword, apperror.BadRequestError, "no_new_password", "password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Type != tt.errType {
				t.Errorf("Type = %q, want %q", tt.err.Type, tt.errType)
			}
			if tt.err.Code != tt.code {
				t.Errorf("Code = %q, want %q", tt.err.Code, tt.code)
			}
			if tt.err.Message == "" {
				t.Error("Message should not be empty")
			}
		})
	}
}
