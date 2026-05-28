package auth

import (
	"context"
	"testing"
	"time"
	"wn/internal/domain/dto"
	"wn/internal/domain/enum"
	"wn/internal/domain/services/token"
	apperrors "wn/internal/errors"
	"wn/pkg/applogger"
	"wn/pkg/constants"

	"github.com/google/uuid"
)

type mockAppUserService struct {
	user *dto.User
	err  error
	updateErr error
}

func (m *mockAppUserService) GetUserByEmail(_ context.Context, _ string, _ string) (*dto.User, error) {
	return m.user, m.err
}

func (m *mockAppUserService) UpdateUser(_ context.Context, _ uuid.UUID, _ *dto.UserUpdateParams) error {
	return m.updateErr
}

type mockAppTokenService struct {
	tokens *token.UserTokens
	err    error
}

func (m *mockAppTokenService) GenerateUserTokens(_ context.Context, _, _ uuid.UUID, _ string) (*token.UserTokens, error) {
	return m.tokens, m.err
}

func (m *mockAppTokenService) ParseToken(_ string, _ bool) (*token.CustomClaims, error) {
	return nil, nil
}

func (m *mockAppTokenService) RefreshTokens(_ context.Context, _, _ string) (*token.UserTokens, error) {
	return m.tokens, m.err
}

type mockAppSmtpService struct {
	sendErr   error
	code      *dto.ConfirmationCode
	confirmErr error
}

func (m *mockAppSmtpService) SendConfirmEmailCode(_ context.Context, _ string, _ enum.EmailCodeAction) error {
	return m.sendErr
}

func (m *mockAppSmtpService) ConfirmCode(_ context.Context, _ string, _ string) (*dto.ConfirmationCode, error) {
	return m.code, m.confirmErr
}

type mockLayoutService struct {
	layouts []dto.Layout
	err     error
}

func (m *mockLayoutService) GetAvailableLayouts(_ context.Context, _ uuid.UUID) ([]dto.Layout, error) {
	return m.layouts, m.err
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

func newTestAppService(user *mockAppUserService, smtp *mockAppSmtpService, tokenSvc *mockAppTokenService, layout *mockLayoutService) *Service {
	return NewService(noopTx{}, noopLogger{}, user, smtp, tokenSvc, layout)
}

func TestConfirmCode_ConfirmEmail(t *testing.T) {
	userID := uuid.New()
	srv := newTestAppService(
		&mockAppUserService{user: &dto.User{Id: userID, Email: "a@test.com"}},
		&mockAppSmtpService{code: &dto.ConfirmationCode{Action: enum.ConfirmCode, Code: "123"}},
		nil,
		nil,
	)

	err := srv.ConfirmCode(context.Background(), dto.ConfirmationCodeRequest{
		Email: "a@test.com",
		Code:  "123",
	})
	if err != nil {
		t.Fatalf("ConfirmCode() error = %v", err)
	}
}

func TestConfirmCode_ForgotPassword_NoNewPassword(t *testing.T) {
	srv := newTestAppService(
		&mockAppUserService{user: &dto.User{Id: uuid.New()}},
		&mockAppSmtpService{code: &dto.ConfirmationCode{Action: enum.ForgotPassword, Code: "123"}},
		nil,
		nil,
	)

	err := srv.ConfirmCode(context.Background(), dto.ConfirmationCodeRequest{
		Email: "a@test.com",
		Code:  "123",
	})
	if err != apperrors.NoNewPassword {
		t.Errorf("error = %v, want NoNewPassword", err)
	}
}

func TestLogin(t *testing.T) {
	userID := uuid.New()
	layoutID := uuid.New()

	srv := newTestAppService(
		&mockAppUserService{user: &dto.User{Id: userID, Role: constants.ClientRole}},
		nil,
		&mockAppTokenService{tokens: &token.UserTokens{Access: "access", Refresh: "refresh"}},
		&mockLayoutService{layouts: []dto.Layout{{Id: layoutID, OwnerId: userID, IsMain: true}}},
	)

	resp, err := srv.Login(context.Background(), dto.LoginRequest{Email: "a@test.com", Password: "pass"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if resp.UserId != userID {
		t.Errorf("UserId = %v, want %v", resp.UserId, userID)
	}
	if resp.Access != "access" || resp.Refresh != "refresh" {
		t.Errorf("tokens = %+v", resp)
	}
}

func TestSendConfirmationCode(t *testing.T) {
	srv := newTestAppService(
		&mockAppUserService{user: &dto.User{Id: uuid.New()}},
		&mockAppSmtpService{},
		nil,
		nil,
	)

	resp, err := srv.SendConfirmationCode(context.Background(), dto.LoginRequest{Email: "a@test.com"}, enum.ConfirmCode)
	if err != nil {
		t.Fatalf("SendConfirmationCode() error = %v", err)
	}
	if resp.NextCodeDelay != time.Minute {
		t.Errorf("NextCodeDelay = %v, want 1m", resp.NextCodeDelay)
	}
}

func TestConfirmCode_ForgotPassword(t *testing.T) {
	userID := uuid.New()
	srv := newTestAppService(
		&mockAppUserService{user: &dto.User{Id: userID, Email: "a@test.com"}},
		&mockAppSmtpService{code: &dto.ConfirmationCode{Action: enum.ForgotPassword, Code: "123"}},
		nil,
		nil,
	)

	err := srv.ConfirmCode(context.Background(), dto.ConfirmationCodeRequest{
		Email:       "a@test.com",
		Code:        "123",
		NewPassword: "new-pass",
	})
	if err != nil {
		t.Fatalf("ConfirmCode() error = %v", err)
	}
}

func TestRefreshTokens(t *testing.T) {
	srv := newTestAppService(
		nil,
		nil,
		&mockAppTokenService{tokens: &token.UserTokens{Access: "new-access", Refresh: "new-refresh"}},
		nil,
	)

	tokens, err := srv.RefreshTokens(context.Background(), token.UserTokens{Access: "old", Refresh: "old"})
	if err != nil {
		t.Fatalf("RefreshTokens() error = %v", err)
	}
	if tokens.Access != "new-access" {
		t.Errorf("Access = %q, want new-access", tokens.Access)
	}
}
