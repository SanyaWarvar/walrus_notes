package layout

import (
	"context"
	"testing"
	"wn/internal/domain/dto"
	"wn/internal/domain/entity"
	apperrors "wn/internal/errors"
	"wn/pkg/applogger"

	"github.com/google/uuid"
)

type noopLogger struct{}

func (noopLogger) IsDebugLevel() bool                         { return false }
func (noopLogger) IsInfoLevel() bool                          { return false }
func (noopLogger) Debug(string)                               {}
func (noopLogger) Info(string)                                {}
func (noopLogger) Warn(string)                                {}
func (noopLogger) Error(string)                               {}
func (noopLogger) Warnf(string, ...any)                       {}
func (noopLogger) Errorf(string, ...any)                      {}
func (noopLogger) Debugf(string, ...any)                      {}
func (noopLogger) Infof(string, ...any)                       {}
func (l noopLogger) WithCtx(context.Context) applogger.Logger { return l }

type noopTx struct{}

func (noopTx) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockLayoutRepo struct {
	layouts []entity.Layout
	created *entity.Layout
}

func (m *mockLayoutRepo) CreateLayout(_ context.Context, item *entity.Layout) (uuid.UUID, error) {
	m.created = item
	return item.Id, nil
}

func (m *mockLayoutRepo) DeleteLayoutById(_ context.Context, _ uuid.UUID) error { return nil }

func (m *mockLayoutRepo) GetAvailableLayouts(_ context.Context, _ uuid.UUID) ([]entity.Layout, error) {
	return m.layouts, nil
}

func (m *mockLayoutRepo) UpdateLayout(_ context.Context, _, _ uuid.UUID, _, _ string) (int, error) {
	return 0, nil
}

type mockPermissionsRepo struct {
	perms []entity.Permission
}

func (m *mockPermissionsRepo) GetPermission(_ context.Context, _ *dto.GetPermissionsFilter) (*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionsRepo) GetPermissions(_ context.Context, _ *dto.GetPermissionsFilter) ([]entity.Permission, error) {
	return m.perms, nil
}

func (m *mockPermissionsRepo) DeletePermissions(_ context.Context, _ ...uuid.UUID) error { return nil }
func (m *mockPermissionsRepo) UpdatePermissions(_ context.Context, _ *entity.Permission) error {
	return nil
}
func (m *mockPermissionsRepo) CreatePermissions(_ context.Context, _ *entity.Permission) error {
	return nil
}

func TestCreateLayout(t *testing.T) {
	repo := &mockLayoutRepo{}
	srv := NewService(noopTx{}, noopLogger{}, repo, nil, nil, nil, nil, &mockPermissionsRepo{})

	ownerID := uuid.New()
	id, err := srv.CreateLayout(context.Background(), "My Layout", "#FF0000", ownerID, true)
	if err != nil {
		t.Fatalf("CreateLayout() error = %v", err)
	}
	if id == uuid.Nil {
		t.Fatal("expected non-nil layout id")
	}
	if repo.created.Title != "My Layout" {
		t.Errorf("Title = %q, want My Layout", repo.created.Title)
	}
	if !repo.created.IsMain {
		t.Error("IsMain should be true")
	}
}

func TestGetAvailableLayouts(t *testing.T) {
	layoutID := uuid.New()
	userID := uuid.New()
	permID := uuid.New()

	repo := &mockLayoutRepo{
		layouts: []entity.Layout{{Id: layoutID, Title: "L1", OwnerId: userID, Color: "#fff"}},
	}
	perms := &mockPermissionsRepo{
		perms: []entity.Permission{
			{Id: permID, TargetId: layoutID, ToUserId: userID, CanRead: true},
		},
	}
	srv := NewService(noopTx{}, noopLogger{}, repo, nil, nil, nil, nil, perms)

	got, err := srv.GetAvailableLayouts(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetAvailableLayouts() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Permission == nil {
		t.Fatal("Permission should be attached")
	}
	if !got[0].Permission.CanRead {
		t.Error("CanRead should be true")
	}
}

func TestUpdateLayout_NotFound(t *testing.T) {
	repo := &mockLayoutRepo{}
	srv := NewService(noopTx{}, noopLogger{}, repo, nil, nil, nil, nil, &mockPermissionsRepo{})

	err := srv.UpdateLayout(context.Background(), dto.UpdateLayout{LayoutId: uuid.New()}, uuid.New())
	if err != apperrors.LayoutNotFound {
		t.Errorf("error = %v, want LayoutNotFound", err)
	}
}
