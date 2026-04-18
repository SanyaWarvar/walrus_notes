package entity

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	Id         uuid.UUID
	ToUserId   uuid.UUID
	FromUserId uuid.UUID
	TargetId   uuid.UUID
	CanRead    bool
	CanWrite   bool
	CanEdit    bool
	CreatedAt  time.Time
}
