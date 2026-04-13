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
