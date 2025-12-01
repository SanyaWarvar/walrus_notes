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

// UploadFileRequest
// @Schema
type UploadFileRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// NoteRequest
// @Schema
type NoteRequest struct {
	Title    string    `json:"title"`
	Payload  string    `json:"payload"`
	LayoutId uuid.UUID `json:"layoutId"`
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

// NewLayoutRequest
// @Schema
type NewLayoutRequest struct {
	Title string `json:"title" binding:"required"`
	Color string `json:"color" binding:"required"`
}

// UpdateLayout
// @Schema
type UpdateLayout struct {
	LayoutId uuid.UUID `json:"layoutId" binding:"required"`
	Title    string    `json:"title"`
	Color    string    `json:"color"`
}

// LayoutIdRequest
// @Schema
type LayoutIdRequest struct {
	LayoutId uuid.UUID `json:"layoutId"`
}

// GetNotesFromLayoutRequest
// @Schema
type GetNotesFromLayoutRequest struct {
	LayoutId uuid.UUID `json:"layoutId"`
	Page     int       `json:"page"`
}

// GetNotesFromLayoutWithoutPagRequest
// @Schema
type GetNotesFromLayoutWithoutPagRequest struct {
	LayoutId uuid.UUID `json:"layoutId" binding:"required"`
}

// UpdateNotePositionRequest
// @Schema
type UpdateNotePositionRequest struct {
	LayoutId uuid.UUID `json:"layoutId" binding:"required"`
	NoteId   uuid.UUID `json:"noteId" binding:"required"`
	XPos     *float64  `json:"xPos"`
	YPos     *float64  `json:"yPos"`
}

type LinkBetweenNotesRequest struct {
	LayoutId     uuid.UUID `json:"layoutId" binding:"required"`
	FirstNoteId  uuid.UUID `json:"firstNoteId" binding:"required"`
	SecondNoteId uuid.UUID `json:"secondNoteId" binding:"required"`
}
