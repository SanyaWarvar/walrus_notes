//go:build integration

package integration

import (
	"testing"
	"time"
	"wn/internal/domain/entity"
	tokenrepo "wn/internal/infrastructure/repository/tokens"
	"wn/pkg/util"

	"github.com/google/uuid"
)

func TestTokenRepository_CRUD(t *testing.T) {
	resetDB(t)
	user := seedUser(t)
	repo := tokenrepo.NewRepository(env.DB)

	token := &entity.RefreshToken{
		Id:       uuid.New(),
		UserId:   user.Id,
		AccessId: uuid.New(),
		ExpAt:    util.GetCurrentUTCTime().Add(time.Hour),
	}

	if err := repo.Create(env.Ctx, token); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, exists, err := repo.GetByID(env.Ctx, token.Id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if !exists {
		t.Fatal("token should exist")
	}
	if got.UserId != user.Id {
		t.Errorf("UserId = %v, want %v", got.UserId, user.Id)
	}

	if err := repo.Delete(env.Ctx, token.Id); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, exists, err = repo.GetByID(env.Ctx, token.Id)
	if err != nil {
		t.Fatalf("GetByID() after delete error = %v", err)
	}
	if exists {
		t.Fatal("token should be deleted")
	}
}

func TestTokenRepository_DeleteExpired(t *testing.T) {
	resetDB(t)
	user := seedUser(t)
	repo := tokenrepo.NewRepository(env.DB)

	expired := &entity.RefreshToken{
		Id:       uuid.New(),
		UserId:   user.Id,
		AccessId: uuid.New(),
		ExpAt:    util.GetCurrentUTCTime().Add(-time.Hour),
	}
	active := &entity.RefreshToken{
		Id:       uuid.New(),
		UserId:   user.Id,
		AccessId: uuid.New(),
		ExpAt:    util.GetCurrentUTCTime().Add(time.Hour),
	}

	for _, tok := range []*entity.RefreshToken{expired, active} {
		if err := repo.Create(env.Ctx, tok); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	if err := repo.DeleteExpired(env.Ctx, util.GetCurrentUTCTime()); err != nil {
		t.Fatalf("DeleteExpired() error = %v", err)
	}

	_, exists, _ := repo.GetByID(env.Ctx, expired.Id)
	if exists {
		t.Error("expired token should be deleted")
	}
	_, exists, _ = repo.GetByID(env.Ctx, active.Id)
	if !exists {
		t.Error("active token should remain")
	}
}
