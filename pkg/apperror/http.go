package apperror

import (
	"net/http"

	"github.com/pkg/errors"
)

func GetErrorByHttpStatus(status int, message, code string) error {
	switch status {
	case http.StatusBadRequest:
		return NewBadRequestError(message, code)
	case http.StatusUnauthorized:
		return NewUnauthorizedError(message, code)
	case http.StatusForbidden:
		return NewAccessDeniedError(message, code)
	case http.StatusConflict:
		return NewConflictError(message, code)
	case http.StatusGone:
		return NewNotFoundError(message, code)
	case http.StatusUnprocessableEntity:
		return NewInvalidDataError(message, code)
	default:
		return NewInternalError(errors.New(message))
	}
}

func GetHttpStatusByErrorType(errType ErrType) int {
	status := http.StatusInternalServerError
	switch errType {
	case AccessDeniedError:
		status = http.StatusForbidden
	case InvalidDataError:
		status = http.StatusUnprocessableEntity
	case ConflictError:
		status = http.StatusConflict
	case NotFoundError:
		status = http.StatusGone
	case UnauthorizedError:
		status = http.StatusUnauthorized
	case BadRequestError:
		status = http.StatusBadRequest
	}
	return status
}
