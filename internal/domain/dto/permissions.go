package dto

import (
	"time"
	"wn/internal/domain/enum"

	"github.com/google/uuid"
)

type PermissionsToken struct {
	FromUserId uuid.UUID            `json:"fromUserId"`
	TargetId   uuid.UUID            `json:"targetId"`
	Kind       enum.PermissionsKind `json:"kind"`
	CanRead    bool                 `json:"canRead"`
	CanWrite   bool                 `json:"canWrite"`
	CanEdit    bool                 `json:"canEdit"`
	ExpiredAt  time.Time            `json:"expiredAt"`
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
