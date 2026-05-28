package dto

import (
	"testing"
	"time"
	"wn/internal/domain/entity"
	"wn/pkg/constants"

	"github.com/google/uuid"
)

func TestUserDtoFromEntity(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	createdAt := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	src := &entity.User{
		Id:             id,
		Username:       "alice",
		Email:          "alice@example.com",
		Password:       "hashed",
		Role:           constants.ClientRole,
		ImgUrl:         "avatar.png",
		ConfirmedEmail: true,
		CreatedAt:      createdAt,
	}

	got := UserDtoFromEntity(src)

	if got.Id != id {
		t.Errorf("Id = %v, want %v", got.Id, id)
	}
	if got.Username != "alice" {
		t.Errorf("Username = %q, want %q", got.Username, "alice")
	}
	if got.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", got.Email, "alice@example.com")
	}
	if got.Role != constants.ClientRole {
		t.Errorf("Role = %q, want %q", got.Role, constants.ClientRole)
	}
	if got.ImgUrl != "avatar.png" {
		t.Errorf("ImgUrl = %q, want %q", got.ImgUrl, "avatar.png")
	}
	if !got.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, createdAt)
	}
}
