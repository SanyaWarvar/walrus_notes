package auth

import (
	"context"
	"testing"
	"wn/internal/domain/dto"
	"wn/internal/domain/entity"
	apperrors "wn/internal/errors"
	"wn/pkg/applogger"
	"wn/pkg/constants"

	"github.com/google/uuid"
)

type mockUserRepo struct {
	users map[uuid.UUID]*entity.User
	byEmail map[string]*entity.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[uuid.UUID]*entity.User),
		byEmail: make(map[string]*entity.User),
	}
}

func (m *mockUserRepo) CreateUser(_ context.Context, item *entity.User) error {
	m.users[item.Id] = item
	m.byEmail[item.Email] = item
	return nil
}

func (m *mockUserRepo) GetUser(_ context.Context, filter dto.UserFilter) (*entity.User, bool, error) {
	if filter.Id != nil {
		u, ok := m.users[*filter.Id]
		return u, ok, nil
	}
	if filter.Email != nil {
		u, ok := m.byEmail[*filter.Email]
		return u, ok, nil
	}
	return nil, false, nil
}

func (m *mockUserRepo) UpdateUser(_ context.Context, userId uuid.UUID, updateParams *dto.UserUpdateParams) error {
	u, ok := m.users[userId]
	if !ok {
		return apperrors.UserNotFound
	}
	if updateParams.Password != nil {
		u.Password = *updateParams.Password
	}
	if updateParams.ConfirmedEmail != nil {
		u.ConfirmedEmail = *updateParams.ConfirmedEmail
	}
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

func newTestUserService() *Service {
	return NewService("salt", noopTx{}, noopLogger{}, newMockUserRepo())
}

func TestCreateUserFromAuthCredentials(t *testing.T) {
	srv := newTestUserService()

	user, err := srv.CreateUserFromAuthCredentials(context.Background(), dto.RegisterCredentials{
		Username: "bob",
		Email:    "bob@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("CreateUserFromAuthCredentials() error = %v", err)
	}
	if user.Username != "bob" {
		t.Errorf("Username = %q, want %q", user.Username, "bob")
	}

	repo := srv.userRepo.(*mockUserRepo)
	stored := repo.byEmail["bob@example.com"]
	if stored == nil {
		t.Fatal("user was not stored")
	}
	if stored.Password == "password123" {
		t.Error("password should be hashed")
	}
	if stored.Role != constants.ClientRole {
		t.Errorf("Role = %q, want %q", stored.Role, constants.ClientRole)
	}
}

func TestGetUserByEmail_WrongPassword(t *testing.T) {
	srv := newTestUserService()
	repo := srv.userRepo.(*mockUserRepo)

	id := uuid.New()
	repo.byEmail["alice@example.com"] = &entity.User{
		Id:       id,
		Email:    "alice@example.com",
		Password: srv.generatePasswordHash("correct"),
	}

	_, err := srv.GetUserByEmail(context.Background(), "alice@example.com", "wrong")
	if err != apperrors.IncorrectPassword {
		t.Errorf("error = %v, want IncorrectPassword", err)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	srv := newTestUserService()
	repo := srv.userRepo.(*mockUserRepo)

	id := uuid.New()
	repo.byEmail["alice@example.com"] = &entity.User{
		Id:       id,
		Email:    "alice@example.com",
		Username: "alice",
		Password: srv.generatePasswordHash("secret"),
		Role:     constants.ClientRole,
	}

	user, err := srv.GetUserByEmail(context.Background(), "alice@example.com", "secret")
	if err != nil {
		t.Fatalf("GetUserByEmail() error = %v", err)
	}
	if user.Id != id {
		t.Errorf("Id = %v, want %v", user.Id, id)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	srv := newTestUserService()

	err := srv.UpdateUser(context.Background(), uuid.New(), &dto.UserUpdateParams{})
	if err != apperrors.UserNotFound {
		t.Errorf("error = %v, want UserNotFound", err)
	}
}

func TestUpdateUser_HashesPassword(t *testing.T) {
	srv := newTestUserService()
	repo := srv.userRepo.(*mockUserRepo)

	id := uuid.New()
	repo.users[id] = &entity.User{Id: id, Password: "old"}

	newPass := "new-password"
	err := srv.UpdateUser(context.Background(), id, &dto.UserUpdateParams{Password: &newPass})
	if err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}
	if repo.users[id].Password == newPass {
		t.Error("password should be hashed on update")
	}
}
