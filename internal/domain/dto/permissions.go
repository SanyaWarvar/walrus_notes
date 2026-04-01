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
	Id         uuid.UUID            `json:"id"`
	FromUserId uuid.UUID            `json:"fromUserId"`
	ToUserId   uuid.UUID            `json:"toUserId"`
	TargetId   uuid.UUID            `json:"targetId"`
	Kind       enum.PermissionsKind `json:"kind"`
	CanRead    bool                 `json:"canRead"`
	CanWrite   bool                 `json:"canWrite"`
	CanEdit    bool                 `json:"canEdit"`
}

func PermissionFromEntity(e *entity.Permission) *Permission {
	return &Permission{
		Id:         e.Id,
		FromUserId: e.FromUserId,
		ToUserId:   e.ToUserId,
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

type PermissionsDashboard struct {
	Received []Permission `json:"received"`
	Shared   []Permission `json:"shared"`
}

type DeletePermissionsRequest struct {
	PermissionId uuid.UUID `json:"permissionId"`
}

type UpdatePermissionRequest struct {
	PermissionId uuid.UUID `json:"permissionsId"`
	CanWrite     bool      `json:"canWrite"`
	CanEdit      bool      `json:"canEdit"`
	CanRead      bool      `json:"canRead"`
}
