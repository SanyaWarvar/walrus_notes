package layout

import (
	"context"
	"testing"
	"wn/internal/domain/dto"
	"wn/internal/domain/events"
	"wn/pkg/applogger"

	"github.com/google/uuid"
)

type mockLayoutService struct {
	layouts []dto.Layout
	id      uuid.UUID
}

func (m *mockLayoutService) CreateLayout(_ context.Context, _, _ string, _ uuid.UUID, _ bool) (uuid.UUID, error) {
	return m.id, nil
}
func (m *mockLayoutService) DeleteLayoutById(_ context.Context, _, _ uuid.UUID) error { return nil }
func (m *mockLayoutService) GetAvailableLayouts(_ context.Context, _ uuid.UUID) ([]dto.Layout, error) {
	return m.layouts, nil
}
func (m *mockLayoutService) ExportLayouts(_ context.Context, _ uuid.UUID) (*dto.ExportInfo, error) {
	return nil, nil
}
func (m *mockLayoutService) UpdateLayout(_ context.Context, _ dto.UpdateLayout, _ uuid.UUID) error {
	return nil
}
func (m *mockLayoutService) ImportLayouts(_ context.Context, _ uuid.UUID, _ *dto.ExportInfo) error {
	return nil
}

type mockPermService struct{}

func (m *mockPermService) CheckPermissionByLayoutId(_ context.Context, _, _ uuid.UUID, _, _, _ bool) error {
	return nil
}
func (m *mockPermService) GetAssociatedUsersByLayout(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{uuid.New()}, nil
}

type mockProducer struct{}

func (m *mockProducer) SendToAssociatedUsers(_ context.Context, _ uuid.UUID, _ []uuid.UUID, _ events.Event) error {
	return nil
}

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

func TestCreateLayout(t *testing.T) {
	layoutID := uuid.New()
	srv := NewService(noopTx{}, noopLogger{}, &mockLayoutService{id: layoutID}, &mockPermService{}, &mockProducer{})

	id, err := srv.CreateLayout(context.Background(), dto.NewLayoutRequest{Title: "T", Color: "#000"}, uuid.New())
	if err != nil {
		t.Fatalf("CreateLayout() error = %v", err)
	}
	if id != layoutID {
		t.Errorf("id = %v, want %v", id, layoutID)
	}
}

func TestGetLayoutsByUserId(t *testing.T) {
	userID := uuid.New()
	layouts := []dto.Layout{{Id: uuid.New(), OwnerId: userID}}
	srv := NewService(noopTx{}, noopLogger{}, &mockLayoutService{layouts: layouts}, &mockPermService{}, &mockProducer{})

	got, err := srv.GetLayoutsByUserId(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetLayoutsByUserId() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}
