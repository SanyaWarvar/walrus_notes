//go:build integration

package integration

import (
	"testing"
	"wn/internal/domain/dto"
	"wn/internal/domain/entity"
	apperrors "wn/internal/errors"
	userrepo "wn/internal/infrastructure/repository/user"
	"wn/pkg/util"

	"github.com/google/uuid"
)

func TestUserRepository_CreateAndGet(t *testing.T) {
	resetDB(t)
	repo := userrepo.NewRepository(env.DB)

	user := newTestUser()
	if err := repo.CreateUser(env.Ctx, user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	got, exists, err := repo.GetUser(env.Ctx, dto.UserFilter{Email: &user.Email})
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if !exists {
		t.Fatal("user should exist")
	}
	if got.Username != user.Username {
		t.Errorf("Username = %q, want %q", got.Username, user.Username)
	}
}

func TestUserRepository_CreateUser_DuplicateEmail(t *testing.T) {
	resetDB(t)
	repo := userrepo.NewRepository(env.DB)
	user := newTestUser()

	if err := repo.CreateUser(env.Ctx, user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	dup := newTestUser()
	dup.Email = user.Email
	err := repo.CreateUser(env.Ctx, dup)
	if err != apperrors.NotUnique {
		t.Errorf("error = %v, want NotUnique", err)
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	resetDB(t)
	repo := userrepo.NewRepository(env.DB)
	user := newTestUser()
	if err := repo.CreateUser(env.Ctx, user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	newEmail := uuid.New().String() + "@updated.test"
	confirmed := true
	if err := repo.UpdateUser(env.Ctx, user.Id, &dto.UserUpdateParams{
		Email:          &newEmail,
		ConfirmedEmail: &confirmed,
	}); err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}

	got, exists, err := repo.GetUser(env.Ctx, dto.UserFilter{Id: &user.Id})
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if !exists {
		t.Fatal("user should exist")
	}
	if got.Email != newEmail {
		t.Errorf("Email = %q, want %q", got.Email, newEmail)
	}
	if !got.ConfirmedEmail {
		t.Error("ConfirmedEmail should be true")
	}
}

func TestUserRepository_GetUser_NotFound(t *testing.T) {
	resetDB(t)
	repo := userrepo.NewRepository(env.DB)

	id := uuid.New()
	_, exists, err := repo.GetUser(env.Ctx, dto.UserFilter{Id: &id})
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if exists {
		t.Fatal("user should not exist")
	}
}

func seedUser(t *testing.T) *entity.User {
	t.Helper()
	repo := userrepo.NewRepository(env.DB)
	user := newTestUser()
	user.CreatedAt = util.GetCurrentUTCTime()
	if err := repo.CreateUser(env.Ctx, user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	return user
}
