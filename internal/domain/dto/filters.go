package dto

import (
	"github.com/google/uuid"
)

type GetPermissionsFilter struct {
	Id         *uuid.UUID
	FromUserId *uuid.UUID
	ToUserId   *uuid.UUID
	TargetId   *uuid.UUID
	TargetIdIn []uuid.UUID

	Limit uint64
}

type UserFilter struct {
	Id    *uuid.UUID
	Email *string

	Limit uint64
}

// пароль передавать незахешированным. На уровне сервиса произойдет хеш
type UserUpdateParams struct {
	Username       *string
	Email          *string
	Password       *string
	ImgUrl         *string
	ConfirmedEmail *bool
}
