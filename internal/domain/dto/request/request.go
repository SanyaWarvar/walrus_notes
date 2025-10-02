package request

import (
	"mime/multipart"

	"github.com/google/uuid"
)

// RegisterCredentials
// @Schema
type RegisterCredentials struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest
// @Schema
type LoginRequest struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"`
}

// ConfimationCodeRequest
// @Schema
type ConfimationCodeRequest struct {
	Code        string `json:"code"  binding:"required"`
	Email       string `json:"email" binding:"required"`
	NewPassword string `json:"newPassword"`
}

// ForgotPasswordRequest
// @Schema
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

// ChangeProfilePicture
// @Schema
type ChangeProfilePicture struct {
	File   *multipart.FileHeader `form:"file" binding:"required"`
	UserId uuid.UUID
}

// NoteRequest
// @Schema
type NoteRequest struct {
	Title   string `json:"title"`
	Payload string `json:"payload"`
}

// NoteWithIdRequest
// @Schema
type NoteWithIdRequest struct {
	NoteId  uuid.UUID `json:"noteId" binding:"required"`
	Title   string    `json:"title"`
	Payload string    `json:"payload"`
}

// NoteId
// @Schema
type NoteId struct {
	NoteId uuid.UUID `json:"noteId" binding:"required"`
}
