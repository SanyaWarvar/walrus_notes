//go:build integration

package integration

import (
	"context"
	"testing"
	"time"
	"wn/internal/domain/dto"
	"wn/internal/domain/enum"
	smtpcache "wn/internal/infrastructure/cache/smtp"
	"wn/pkg/applogger"
	"wn/pkg/util"
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

func TestSMTPCache_SaveAndGetConfirmCode(t *testing.T) {
	cache := smtpcache.NewCache(noopLogger{}, env.Redis)
	email := "integration@test.com"
	ttl := time.Minute

	code := dto.ConfirmationCode{
		Code:      "123456",
		Action:    enum.ConfirmCode,
		CreatedAt: util.GetCurrentUTCTime(),
	}

	if err := cache.SaveConfirmCode(env.Ctx, email, code, &ttl); err != nil {
		t.Fatalf("SaveConfirmCode() error = %v", err)
	}

	got, exists, err := cache.GetConfirmCode(env.Ctx, email)
	if err != nil {
		t.Fatalf("GetConfirmCode() error = %v", err)
	}
	if !exists {
		t.Fatal("code should exist")
	}
	if got.Code != code.Code {
		t.Errorf("Code = %q, want %q", got.Code, code.Code)
	}
	if got.Action != enum.ConfirmCode {
		t.Errorf("Action = %q, want CONFIRM_CODE", got.Action)
	}
}

func TestSMTPCache_GetConfirmCode_NotFound(t *testing.T) {
	cache := smtpcache.NewCache(noopLogger{}, env.Redis)

	_, exists, err := cache.GetConfirmCode(env.Ctx, "missing@test.com")
	if err != nil {
		t.Fatalf("GetConfirmCode() error = %v", err)
	}
	if exists {
		t.Fatal("code should not exist")
	}
}
