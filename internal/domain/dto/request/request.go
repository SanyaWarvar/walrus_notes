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
