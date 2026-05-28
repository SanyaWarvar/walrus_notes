package auth

import (
	"context"
	"mime/multipart"
	"testing"
	"wn/internal/domain/dto"
	"wn/pkg/applogger"

	"github.com/google/uuid"
)

type mockDomainUserService struct {
	user *dto.User
	err  error
}

func (m *mockDomainUserService) CreateUserFromAuthCredentials(_ context.Context, _ dto.RegisterCredentials) (*dto.User, error) {
	return m.user, m.err
}

func (m *mockDomainUserService) UpdateUser(_ context.Context, _ uuid.UUID, _ *dto.UserUpdateParams) error {
	return m.err
}

func (m *mockDomainUserService) GetUserById(_ context.Context, _ uuid.UUID, _ string) (*dto.User, error) {
	return m.user, m.err
}

type mockFileService struct {
	filename string
	err      error
}

func (m *mockFileService) NewFile(_ context.Context, _ *multipart.FileHeader) (string, error) {
	return m.filename, m.err
}

type mockLayoutService struct {
	layoutID uuid.UUID
	err      error
}

func (m *mockLayoutService) CreateLayout(_ context.Context, _, _ string, _ uuid.UUID, _ bool) (uuid.UUID, error) {
	return m.layoutID, m.err
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

func TestRegisterUser(t *testing.T) {
	userID := uuid.New()
	layoutID := uuid.New()

	srv := NewService(
		noopTx{},
		noopLogger{},
		&mockDomainUserService{user: &dto.User{Id: userID}},
		nil,
		&mockLayoutService{layoutID: layoutID},
	)

	resp, err := srv.RegisterUser(context.Background(), dto.RegisterCredentials{
		Username: "alice",
		Email:    "alice@test.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}
	if resp.UserId != userID {
		t.Errorf("UserId = %v, want %v", resp.UserId, userID)
	}
}

func TestGetUserById(t *testing.T) {
	userID := uuid.New()
	srv := NewService(
		noopTx{},
		noopLogger{},
		&mockDomainUserService{user: &dto.User{Id: userID, ImgUrl: "pic.png"}},
		nil,
		nil,
	)

	user, err := srv.GetUserById(context.Background(), userID, "http://host")
	if err != nil {
		t.Fatalf("GetUserById() error = %v", err)
	}
	if user.ImgUrl != "http://host/statics/images/pic.png" {
		t.Errorf("ImgUrl = %q, want prefixed url", user.ImgUrl)
	}
}
