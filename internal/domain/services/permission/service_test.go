package permission

import (
	"context"
	"testing"
	"wn/internal/domain/dto"
	"wn/internal/domain/entity"
	apperrors "wn/internal/errors"

	"github.com/google/uuid"
)

type mockPermissionsRepo struct {
	permission *entity.Permission
	perms      []entity.Permission
}

func (m *mockPermissionsRepo) GetPermission(_ context.Context, _ *dto.GetPermissionsFilter) (*entity.Permission, error) {
	return m.permission, nil
}

func (m *mockPermissionsRepo) GetPermissions(_ context.Context, _ *dto.GetPermissionsFilter) ([]entity.Permission, error) {
	return m.perms, nil
}

func (m *mockPermissionsRepo) DeletePermissions(_ context.Context, _ ...uuid.UUID) error {
	return nil
}

func (m *mockPermissionsRepo) UpdatePermissions(_ context.Context, _ *entity.Permission) error {
	return nil
}

func (m *mockPermissionsRepo) CreatePermissions(_ context.Context, _ *entity.Permission) error {
	return nil
}

type mockLayoutRepo struct {
	layout *entity.Layout
}

func (m *mockLayoutRepo) GetById(_ context.Context, _ uuid.UUID) (*entity.Layout, error) {
	return m.layout, nil
}

type mockNoteRepo struct {
	note *entity.Note
}

func (m *mockNoteRepo) GetById(_ context.Context, _ uuid.UUID) (*entity.Note, error) {
	return m.note, nil
}

func TestApplyUpdateRequest(t *testing.T) {
	srv := NewPermissionsService(nil, nil, nil)
	perm := &entity.Permission{CanRead: false, CanWrite: false, CanEdit: false}
	req := &dto.UpdatePermissionRequest{CanRead: true, CanWrite: true, CanEdit: false}

	got := srv.ApplyUpdateRequest(req, perm)

	if !got.CanRead || !got.CanWrite || got.CanEdit {
		t.Errorf("permissions = read:%v write:%v edit:%v", got.CanRead, got.CanWrite, got.CanEdit)
	}
}

func TestCheckPermissionByLayoutId_Owner(t *testing.T) {
	ownerID := uuid.New()
	layoutID := uuid.New()

	srv := NewPermissionsService(
		&mockPermissionsRepo{},
		&mockLayoutRepo{layout: &entity.Layout{Id: layoutID, OwnerId: ownerID}},
		nil,
	)

	err := srv.CheckPermissionByLayoutId(context.Background(), layoutID, ownerID, true, true, true)
	if err != nil {
		t.Errorf("owner should have access, got error: %v", err)
	}
}

func TestCheckPermissionByLayoutId_InsufficientPermissions(t *testing.T) {
	userID := uuid.New()
	layoutID := uuid.New()

	srv := NewPermissionsService(
		&mockPermissionsRepo{permission: &entity.Permission{CanRead: true, CanWrite: false, CanEdit: false}},
		&mockLayoutRepo{layout: &entity.Layout{Id: layoutID, OwnerId: uuid.New()}},
		nil,
	)

	err := srv.CheckPermissionByLayoutId(context.Background(), layoutID, userID, false, true, false)
	if err != apperrors.PermissionsNotEnough {
		t.Errorf("error = %v, want PermissionsNotEnough", err)
	}
}

func TestCheckPermissionByNoteId_Owner(t *testing.T) {
	ownerID := uuid.New()
	noteID := uuid.New()
	layoutID := uuid.New()

	srv := NewPermissionsService(nil, nil, &mockNoteRepo{
		note: &entity.Note{Id: noteID, OwnerId: ownerID, LayoutId: layoutID},
	})

	err := srv.CheckPermissionByNoteId(context.Background(), noteID, ownerID, true, true, true)
	if err != nil {
		t.Errorf("note owner should have access, got error: %v", err)
	}
}

func TestCheckPermissionByNoteId_LayoutOwner(t *testing.T) {
	ownerID := uuid.New()
	noteID := uuid.New()
	layoutID := uuid.New()

	srv := NewPermissionsService(
		&mockPermissionsRepo{},
		&mockLayoutRepo{layout: &entity.Layout{Id: layoutID, OwnerId: ownerID}},
		&mockNoteRepo{note: &entity.Note{Id: noteID, OwnerId: uuid.New(), LayoutId: layoutID}},
	)

	err := srv.CheckPermissionByNoteId(context.Background(), noteID, ownerID, true, true, true)
	if err != nil {
		t.Errorf("layout owner should have access, got error: %v", err)
	}
}

func TestGetAssociatedUsersByNote(t *testing.T) {
	ownerID := uuid.New()
	user1 := uuid.New()
	noteID := uuid.New()

	srv := NewPermissionsService(
		&mockPermissionsRepo{perms: []entity.Permission{{ToUserId: user1}}},
		nil,
		&mockNoteRepo{note: &entity.Note{Id: noteID, OwnerId: ownerID}},
	)

	users, err := srv.GetAssociatedUsersByNote(context.Background(), noteID)
	if err != nil {
		t.Fatalf("GetAssociatedUsersByNote() error = %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("len = %d, want 2", len(users))
	}
}

func TestGetAssociatedUsersByLayout(t *testing.T) {
	ownerID := uuid.New()
	user1 := uuid.New()
	user2 := uuid.New()
	layoutID := uuid.New()

	srv := NewPermissionsService(
		&mockPermissionsRepo{perms: []entity.Permission{
			{ToUserId: user1},
			{ToUserId: user2},
		}},
		&mockLayoutRepo{layout: &entity.Layout{Id: layoutID, OwnerId: ownerID}},
		nil,
	)

	users, err := srv.GetAssociatedUsersByLayout(context.Background(), layoutID)
	if err != nil {
		t.Fatalf("GetAssociatedUsersByLayout() error = %v", err)
	}
	if len(users) != 3 {
		t.Fatalf("len = %d, want 3", len(users))
	}
	if users[2] != ownerID {
		t.Errorf("last user = %v, want owner %v", users[2], ownerID)
	}
}
