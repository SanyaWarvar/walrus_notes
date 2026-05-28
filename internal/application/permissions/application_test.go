package permissions

import (
	"context"
	"testing"
	"time"
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

type mockPermService struct {
	checkErr error
}

func (m *mockPermService) CheckPermissionByLayoutId(_ context.Context, _, _ uuid.UUID, _, _, _ bool) error {
	return m.checkErr
}

func (m *mockPermService) ApplyUpdateRequest(req *dto.UpdatePermissionRequest, e *entity.Permission) *entity.Permission {
	e.CanRead = req.CanRead
	e.CanWrite = req.CanWrite
	e.CanEdit = req.CanEdit
	return e
}

type mockPermRepo struct {
	permission *entity.Permission
	received   []entity.Permission
	shared     []entity.Permission
}

func (m *mockPermRepo) GetPermission(_ context.Context, filter *dto.GetPermissionsFilter) (*entity.Permission, error) {
	if m.permission == nil {
		return nil, apperrors.RecordNotFound
	}
	if filter.Id != nil && m.permission.Id == *filter.Id {
		return m.permission, nil
	}
	if filter.ToUserId != nil && filter.TargetId != nil {
		if m.permission.ToUserId == *filter.ToUserId && m.permission.TargetId == *filter.TargetId {
			return m.permission, nil
		}
		return nil, apperrors.RecordNotFound
	}
	return m.permission, nil
}

func (m *mockPermRepo) GetPermissions(_ context.Context, filter *dto.GetPermissionsFilter) ([]entity.Permission, error) {
	if filter.ToUserId != nil {
		return m.received, nil
	}
	if filter.FromUserId != nil {
		return m.shared, nil
	}
	return nil, nil
}

func (m *mockPermRepo) DeletePermissions(_ context.Context, _ ...uuid.UUID) error { return nil }
func (m *mockPermRepo) UpdatePermissions(_ context.Context, item *entity.Permission) error {
	m.permission = item
	return nil
}
func (m *mockPermRepo) CreatePermissions(_ context.Context, _ *entity.Permission) error { return nil }

type mockLinkRepo struct {
	token *dto.PermissionToken
	exists bool
}

func (m *mockLinkRepo) SavePermissionsLink(_ context.Context, item *dto.PermissionToken, _ uuid.UUID, _ *time.Duration) error {
	m.token = item
	m.exists = true
	return nil
}

func (m *mockLinkRepo) GetPermissionsLink(_ context.Context, _ uuid.UUID) (*dto.PermissionToken, bool, error) {
	return m.token, m.exists, nil
}

func newTestApp(permSvc *mockPermService, permRepo *mockPermRepo, linkRepo *mockLinkRepo) *Application {
	return NewApplication(noopTx{}, noopLogger{}, permSvc, permRepo, linkRepo, nil, nil)
}

func TestGetPermissionsDashboard(t *testing.T) {
	userID := uuid.New()
	repo := &mockPermRepo{
		received: []entity.Permission{{Id: uuid.New(), ToUserId: userID, CanRead: true}},
		shared:   []entity.Permission{{Id: uuid.New(), FromUserId: userID, CanWrite: true}},
	}
	app := newTestApp(&mockPermService{}, repo, &mockLinkRepo{})

	dash, err := app.GetPermissionsDashboard(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetPermissionsDashboard() error = %v", err)
	}
	if len(dash.Received) != 1 || len(dash.Shared) != 1 {
		t.Fatalf("dashboard = %+v", dash)
	}
}

func TestApplyPermissionsLink_CantApplySelf(t *testing.T) {
	userID := uuid.New()
	linkRepo := &mockLinkRepo{
		exists: true,
		token:  &dto.PermissionToken{FromUserId: userID, TargetId: uuid.New()},
	}
	app := newTestApp(&mockPermService{}, &mockPermRepo{}, linkRepo)

	err := app.ApplyPermissionsLink(context.Background(), userID, &dto.ApplyPermissionsRequest{LinkId: uuid.New()})
	if err != apperrors.CantApply {
		t.Errorf("error = %v, want CantApply", err)
	}
}

func TestApplyPermissionsLink_AlreadyExist(t *testing.T) {
	userID := uuid.New()
	targetID := uuid.New()
	linkRepo := &mockLinkRepo{
		exists: true,
		token:  &dto.PermissionToken{FromUserId: uuid.New(), TargetId: targetID},
	}
	repo := &mockPermRepo{
		permission: &entity.Permission{ToUserId: userID, TargetId: targetID},
	}
	app := newTestApp(&mockPermService{}, repo, linkRepo)

	err := app.ApplyPermissionsLink(context.Background(), userID, &dto.ApplyPermissionsRequest{LinkId: uuid.New()})
	if err != apperrors.AlreadyExist {
		t.Errorf("error = %v, want AlreadyExist", err)
	}
}

func TestDeletePermission_NotEnough(t *testing.T) {
	permID := uuid.New()
	repo := &mockPermRepo{
		permission: &entity.Permission{
			Id:         permID,
			FromUserId: uuid.New(),
			ToUserId:   uuid.New(),
		},
	}
	app := newTestApp(&mockPermService{}, repo, &mockLinkRepo{})

	err := app.DeletePermission(context.Background(), uuid.New(), &dto.DeletePermissionsRequest{PermissionId: permID})
	if err != apperrors.PermissionsNotEnough {
		t.Errorf("error = %v, want PermissionsNotEnough", err)
	}
}

func TestGeneratePermissionsLink_BadExpiry(t *testing.T) {
	app := newTestApp(&mockPermService{}, &mockPermRepo{}, &mockLinkRepo{})

	_, err := app.GeneratePermissionsLink(context.Background(), uuid.New(), &dto.GeneratePermissionLinkRequest{
		TargetId:  uuid.New(),
		ExpiredAt: time.Now().Add(-time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for expired link")
	}
}
