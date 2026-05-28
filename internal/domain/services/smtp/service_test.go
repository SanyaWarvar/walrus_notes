package smtp

import (
	"context"
	"testing"
	"time"
	"wn/internal/domain/dto"
	"wn/internal/domain/enum"
	apperrors "wn/internal/errors"
	"wn/pkg/applogger"
)

type mockCacheRepo struct {
	codes map[string]dto.ConfirmationCode
}

func newMockCacheRepo() *mockCacheRepo {
	return &mockCacheRepo{codes: make(map[string]dto.ConfirmationCode)}
}

func (m *mockCacheRepo) GetConfirmCode(_ context.Context, email string) (*dto.ConfirmationCode, bool, error) {
	code, ok := m.codes[email]
	return &code, ok, nil
}

func (m *mockCacheRepo) SaveConfirmCode(_ context.Context, email string, item dto.ConfirmationCode, _ *time.Duration) error {
	m.codes[email] = item
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

func newTestService() *Service {
	cfg := NewConfig("from@test.com", "pass", "smtp:587", "api-key", 6, 5*time.Minute, time.Minute)
	return NewService(noopLogger{}, cfg, newMockCacheRepo())
}

func TestGenerateConfirmCode(t *testing.T) {
	srv := newTestService()

	code := srv.GenerateConfirmCode(enum.ConfirmCode)
	if len(code.Code) != 6 {
		t.Errorf("code length = %d, want 6", len(code.Code))
	}
	if code.Action != enum.ConfirmCode {
		t.Errorf("Action = %q, want %q", code.Action, enum.ConfirmCode)
	}
	if code.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestConfirmCode_Success(t *testing.T) {
	srv := newTestService()
	repo := srv.cacheRepo.(*mockCacheRepo)

	repo.codes["user@test.com"] = dto.ConfirmationCode{
		Code:   "123456",
		Action: enum.ConfirmCode,
	}

	got, err := srv.ConfirmCode(context.Background(), "user@test.com", "123456")
	if err != nil {
		t.Fatalf("ConfirmCode() error = %v", err)
	}
	if got.Code != "123456" {
		t.Errorf("Code = %q, want %q", got.Code, "123456")
	}
}

func TestConfirmCode_NotExist(t *testing.T) {
	srv := newTestService()

	_, err := srv.ConfirmCode(context.Background(), "missing@test.com", "123456")
	if err != apperrors.ConfirmCodeNotExist {
		t.Errorf("error = %v, want ConfirmCodeNotExist", err)
	}
}

func TestConfirmCode_Incorrect(t *testing.T) {
	srv := newTestService()
	repo := srv.cacheRepo.(*mockCacheRepo)

	repo.codes["user@test.com"] = dto.ConfirmationCode{Code: "123456"}

	_, err := srv.ConfirmCode(context.Background(), "user@test.com", "000000")
	if err != apperrors.ConfirmCodeIncorrect {
		t.Errorf("error = %v, want ConfirmCodeIncorrect", err)
	}
}

func TestSendConfirmEmailCode_SavesCode(t *testing.T) {
	srv := newTestService()
	repo := srv.cacheRepo.(*mockCacheRepo)

	err := srv.SendConfirmEmailCode(context.Background(), "user@test.com", enum.ForgotPassword)
	if err != nil {
		t.Fatalf("SendConfirmEmailCode() error = %v", err)
	}
	if _, ok := repo.codes["user@test.com"]; !ok {
		t.Error("code should be saved in cache")
	}
}
