//go:build integration

package integration

import (
	"testing"
	"time"
	"wn/internal/domain/dto"
	"wn/internal/domain/entity"
	layoutrepo "wn/internal/infrastructure/repository/layout"
	permrepo "wn/internal/infrastructure/repository/permissions"
	"wn/pkg/util"

	"github.com/google/uuid"
)

func TestPermissionsRepository_CRUD(t *testing.T) {
	resetDB(t)
	owner := seedUser(t)
	grantee := seedUser(t)
	layoutRepo := layoutrepo.NewRepository(env.DB)
	repo := permrepo.NewRepository(env.DB)

	layout := &entity.Layout{
		Id:         uuid.New(),
		Title:      "Shared",
		OwnerId:    owner.Id,
		HaveAccess: []uuid.UUID{owner.Id},
		Color:      "#ABCDEF",
	}
	if _, err := layoutRepo.CreateLayout(env.Ctx, layout); err != nil {
		t.Fatalf("CreateLayout() error = %v", err)
	}

	perm := &entity.Permission{
		Id:         uuid.New(),
		ToUserId:   grantee.Id,
		FromUserId: owner.Id,
		TargetId:   layout.Id,
		CanRead:    true,
		CanWrite:   false,
		CanEdit:    false,
		CreatedAt:  util.GetCurrentUTCTime(),
	}
	if err := repo.CreatePermissions(env.Ctx, perm); err != nil {
		t.Fatalf("CreatePermissions() error = %v", err)
	}

	got, err := repo.GetPermission(env.Ctx, &dto.GetPermissionsFilter{
		ToUserId: &grantee.Id,
		TargetId: &layout.Id,
	})
	if err != nil {
		t.Fatalf("GetPermission() error = %v", err)
	}
	if !got.CanRead || got.CanWrite {
		t.Errorf("permission = %+v", got)
	}

	perm.CanWrite = true
	if err := repo.UpdatePermissions(env.Ctx, perm); err != nil {
		t.Fatalf("UpdatePermissions() error = %v", err)
	}

	got, err = repo.GetPermission(env.Ctx, &dto.GetPermissionsFilter{Id: &perm.Id})
	if err != nil {
		t.Fatalf("GetPermission() error = %v", err)
	}
	if !got.CanWrite {
		t.Error("CanWrite should be true after update")
	}

	perms, err := repo.GetPermissions(env.Ctx, &dto.GetPermissionsFilter{TargetId: &layout.Id})
	if err != nil {
		t.Fatalf("GetPermissions() error = %v", err)
	}
	if len(perms) != 1 {
		t.Fatalf("len = %d, want 1", len(perms))
	}

	if err := repo.DeletePermissions(env.Ctx, perm.Id); err != nil {
		t.Fatalf("DeletePermissions() error = %v", err)
	}

	perms, err = repo.GetPermissions(env.Ctx, &dto.GetPermissionsFilter{TargetId: &layout.Id})
	if err != nil {
		t.Fatalf("GetPermissions() error = %v", err)
	}
	if len(perms) != 0 {
		t.Errorf("len = %d, want 0 after delete", len(perms))
	}
}

func TestPermissionsRepository_GrantedUserSeesLayout(t *testing.T) {
	resetDB(t)
	owner := seedUser(t)
	grantee := seedUser(t)
	layoutRepo := layoutrepo.NewRepository(env.DB)
	repo := permrepo.NewRepository(env.DB)

	layout := &entity.Layout{
		Id:         uuid.New(),
		Title:      "Shared Board",
		OwnerId:    owner.Id,
		HaveAccess: []uuid.UUID{owner.Id},
		Color:      "#FFFFFF",
	}
	if _, err := layoutRepo.CreateLayout(env.Ctx, layout); err != nil {
		t.Fatalf("CreateLayout() error = %v", err)
	}

	if err := repo.CreatePermissions(env.Ctx, &entity.Permission{
		Id:         uuid.New(),
		ToUserId:   grantee.Id,
		FromUserId: owner.Id,
		TargetId:   layout.Id,
		CanRead:    true,
		CreatedAt:  time.Now().UTC(),
	}); err != nil {
		t.Fatalf("CreatePermissions() error = %v", err)
	}

	layouts, err := layoutRepo.GetAvailableLayouts(env.Ctx, grantee.Id)
	if err != nil {
		t.Fatalf("GetAvailableLayouts() error = %v", err)
	}
	if len(layouts) != 1 {
		t.Fatalf("grantee should see 1 layout, got %d", len(layouts))
	}
}
