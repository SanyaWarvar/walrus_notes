package dto

import (
	"testing"
	"time"
	"wn/internal/domain/entity"

	"github.com/google/uuid"
)

func TestPermissionFromEntity(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	fromUser := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	toUser := uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	target := uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")

	src := &entity.Permission{
		Id:         id,
		FromUserId: fromUser,
		ToUserId:   toUser,
		TargetId:   target,
		CanRead:    true,
		CanWrite:   false,
		CanEdit:    true,
		CreatedAt:  time.Now(),
	}

	got := PermissionFromEntity(src)

	if got.Id != id {
		t.Errorf("Id = %v, want %v", got.Id, id)
	}
	if got.FromUserId != fromUser {
		t.Errorf("FromUserId = %v, want %v", got.FromUserId, fromUser)
	}
	if got.ToUserId != toUser {
		t.Errorf("ToUserId = %v, want %v", got.ToUserId, toUser)
	}
	if got.TargetId != target {
		t.Errorf("TargetId = %v, want %v", got.TargetId, target)
	}
	if !got.CanRead || got.CanWrite || !got.CanEdit {
		t.Errorf("permissions mismatch: read=%v write=%v edit=%v", got.CanRead, got.CanWrite, got.CanEdit)
	}
}
