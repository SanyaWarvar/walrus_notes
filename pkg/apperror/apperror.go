package apperror

import (
	"encoding/json"
	"errors"
)

type ErrType string

const (
	NotFoundError       ErrType = "NotFound"
	ConflictError       ErrType = "Conflict"
	InternalServerError ErrType = "InternalServer"
	BadRequestError     ErrType = "BadRequest"
	InvalidDataError    ErrType = "InvalidData"
	AccessDeniedError   ErrType = "AccessDenied"
	UnauthorizedError   ErrType = "Unauthorized"

	unexpectedErrorMessage = "something went wrong"
)

type AppError struct {
	Err     error   `json:"error"`
	Type    ErrType `json:"type"`
	Message string  `json:"message,omitempty"`
	Code    string  `json:"code"`
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	marshal, _ := json.Marshal(e)
	return marshal
}

func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

func NewAppError(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{Err: err, Type: InternalServerError, Message: unexpectedErrorMessage}
}

func NewBadRequestError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: BadRequestError, Message: message, Code: code}
}

func NewNotFoundError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: NotFoundError, Message: message, Code: code}
}

func NewAccessDeniedError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: AccessDeniedError, Message: message, Code: code}
}

func NewConflictError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: ConflictError, Message: message, Code: code}
}

func NewInvalidDataError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: InvalidDataError, Message: message, Code: code}
}

func NewUnauthorizedError(message, code string) *AppError {
	return &AppError{Err: errors.New(message), Type: UnauthorizedError, Message: message, Code: code}
}
