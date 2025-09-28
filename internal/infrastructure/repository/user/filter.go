package user

import "github.com/google/uuid"

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
