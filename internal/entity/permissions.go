package entity

import (
	"time"
	"wn/internal/domain/enum"

	"github.com/google/uuid"
)

type Permission struct {
	Id         uuid.UUID
	ToUserId   uuid.UUID
	FromUserId uuid.UUID
	TargetId   uuid.UUID
	Kind       enum.PermissionsKind
	CanRead    bool
	CanWrite   bool
	CanEdit    bool
	CreatedAt  time.Time
}
