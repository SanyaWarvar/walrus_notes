package config

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func configDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Dir(filename)
}

func TestNewConfig(t *testing.T) {
	cfg, err := NewConfig(filepath.Join(configDir(t), "main.yaml"))
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	if cfg.HTTP.Port != "80" {
		t.Errorf("HTTP.Port = %q, want %q", cfg.HTTP.Port, "80")
	}
	if cfg.Postgres.DBName != "postgres" {
		t.Errorf("Postgres.DBName = %q, want postgres", cfg.Postgres.DBName)
	}
	if cfg.Internal.EncryptKey == "" {
		t.Error("Internal.EncryptKey should not be empty")
	}
	if cfg.Jwt.AccessTTL <= 0 {
		t.Error("Jwt.AccessTTL should be positive")
	}
	if cfg.Email.CodeLenght != 6 {
		t.Errorf("Email.CodeLenght = %d, want 6", cfg.Email.CodeLenght)
	}
	if cfg.Email.CodeExp != 5*time.Minute {
		t.Errorf("Email.CodeExp = %v, want 5m", cfg.Email.CodeExp)
	}
}
