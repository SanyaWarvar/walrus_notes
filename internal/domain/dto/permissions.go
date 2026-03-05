package dto

import (
	"time"
	"wn/internal/domain/enum"
	"wn/internal/entity"

	"github.com/google/uuid"
)

type PermissionToken struct {
	FromUserId uuid.UUID            `json:"fromUserId"`
	TargetId   uuid.UUID            `json:"targetId"`
	Kind       enum.PermissionsKind `json:"kind"`
	CanRead    bool                 `json:"canRead"`
	CanWrite   bool                 `json:"canWrite"`
	CanEdit    bool                 `json:"canEdit"`
	ExpiredAt  time.Time            `json:"expiredAt"`
}

type Permission struct {
	FromUserId uuid.UUID            `json:"fromUserId"`
	TargetId   uuid.UUID            `json:"targetId"`
	Kind       enum.PermissionsKind `json:"kind"`
	CanRead    bool                 `json:"canRead"`
	CanWrite   bool                 `json:"canWrite"`
	CanEdit    bool                 `json:"canEdit"`
}

func PermissionFromEntity(e *entity.Permission) *Permission {
	return &Permission{
		FromUserId: e.FromUserId,
		TargetId:   e.TargetId,
		Kind:       e.Kind,
		CanRead:    e.CanRead,
		CanWrite:   e.CanWrite,
		CanEdit:    e.CanEdit,
	}
}

type GeneratePermissionLinkRequest struct {
	TargetId  uuid.UUID            `json:"targetId"`
	Kind      enum.PermissionsKind `json:"kind"`
	CanRead   bool                 `json:"canRead"`
	CanWrite  bool                 `json:"canWrite"`
	CanEdit   bool                 `json:"canEdit"`
	ExpiredAt time.Time            `json:"expiredAt"`
}

type GeneratePermissionsLinkResponse struct {
	LinkId uuid.UUID `json:"linkId"`
}

type ApplyPermissionsRequest struct {
	LinkId uuid.UUID `json:"linkId"`
}

type PermissionsDashbord struct {
	Recivied []Permission `json:"recivied"`
	Shared   []Permission `json:"shared"`
}
