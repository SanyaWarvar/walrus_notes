package dto

import (
	"wn/internal/domain/enum"

	"github.com/google/uuid"
)

type GetPermissionsFilter struct {
	Id         *uuid.UUID
	Kind       *enum.PermissionsKind
	FromUserId *uuid.UUID
	ToUserId   *uuid.UUID
	TargetId   *uuid.UUID

	Limit uint64
}
