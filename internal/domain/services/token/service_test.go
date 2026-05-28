package token

import (
	"context"
	"testing"
	"time"
	"wn/internal/domain/entity"
	apperrors "wn/internal/errors"

	"github.com/google/uuid"
)

const testSecret = "test-secret-key"

type mockTokenRepo struct {
	tokens map[uuid.UUID]*entity.RefreshToken
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[uuid.UUID]*entity.RefreshToken)}
}

func (m *mockTokenRepo) Create(_ context.Context, token *entity.RefreshToken) error {
	m.tokens[token.Id] = token
	return nil
}

func (m *mockTokenRepo) GetByID(_ context.Context, id uuid.UUID) (*entity.RefreshToken, bool, error) {
	t, ok := m.tokens[id]
	return t, ok, nil
}

func (m *mockTokenRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.tokens, id)
	return nil
}

func (m *mockTokenRepo) DeleteExpired(_ context.Context, _ time.Time) error {
	return nil
}

func newTestService() *Service {
	return NewService(time.Hour, time.Minute*15, testSecret, newMockTokenRepo())
}

func TestCreateUserTokens(t *testing.T) {
	srv := newTestService()
	userID := uuid.New()
	layoutID := uuid.New()

	tokens, accessID, refreshID, err := srv.CreateUserTokens(userID, layoutID, "ADMIN")
	if err != nil {
		t.Fatalf("CreateUserTokens() error = %v", err)
	}
	if tokens.Access == "" || tokens.Refresh == "" {
		t.Fatal("expected non-empty access and refresh tokens")
	}
	if accessID == uuid.Nil || refreshID == uuid.Nil {
		t.Fatal("expected non-nil token IDs")
	}

	accessClaims, err := srv.ParseToken(tokens.Access, false)
	if err != nil {
		t.Fatalf("ParseToken(access) error = %v", err)
	}
	if accessClaims.UserId != userID {
		t.Errorf("UserId = %v, want %v", accessClaims.UserId, userID)
	}
	if accessClaims.Role != "ADMIN" {
		t.Errorf("Role = %q, want %q", accessClaims.Role, "ADMIN")
	}
	if accessClaims.MainLayoutId != layoutID {
		t.Errorf("MainLayoutId = %v, want %v", accessClaims.MainLayoutId, layoutID)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	srv := newTestService()

	_, err := srv.ParseToken("invalid.token.here", true)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestGenerateUserTokens(t *testing.T) {
	srv := newTestService()
	userID := uuid.New()
	layoutID := uuid.New()

	tokens, err := srv.GenerateUserTokens(context.Background(), userID, layoutID, "CLIENT")
	if err != nil {
		t.Fatalf("GenerateUserTokens() error = %v", err)
	}
	if tokens.Access == "" || tokens.Refresh == "" {
		t.Fatal("expected non-empty tokens")
	}

	repo := srv.tokenRepo.(*mockTokenRepo)
	if len(repo.tokens) != 1 {
		t.Errorf("stored refresh tokens = %d, want 1", len(repo.tokens))
	}
}

func TestRefreshTokens(t *testing.T) {
	srv := newTestService()
	userID := uuid.New()
	layoutID := uuid.New()

	initial, err := srv.GenerateUserTokens(context.Background(), userID, layoutID, "CLIENT")
	if err != nil {
		t.Fatalf("GenerateUserTokens() error = %v", err)
	}

	refreshed, err := srv.RefreshTokens(context.Background(), initial.Access, initial.Refresh)
	if err != nil {
		t.Fatalf("RefreshTokens() error = %v", err)
	}
	if refreshed.Access == "" || refreshed.Refresh == "" {
		t.Fatal("expected non-empty refreshed tokens")
	}
	if refreshed.Access == initial.Access {
		t.Error("expected new access token after refresh")
	}
}

func TestRefreshTokens_TokensDontMatch(t *testing.T) {
	srv := newTestService()
	userID := uuid.New()
	layoutID := uuid.New()

	tokens, err := srv.GenerateUserTokens(context.Background(), userID, layoutID, "CLIENT")
	if err != nil {
		t.Fatalf("GenerateUserTokens() error = %v", err)
	}

	other, _, _, err := srv.CreateUserTokens(uuid.New(), layoutID, "CLIENT")
	if err != nil {
		t.Fatalf("CreateUserTokens() error = %v", err)
	}

	_, err = srv.RefreshTokens(context.Background(), other.Access, tokens.Refresh)
	if err != apperrors.TokensDontMatch {
		t.Errorf("error = %v, want TokensDontMatch", err)
	}
}

func TestRefreshTokens_TokenDontExist(t *testing.T) {
	srv := newTestService()
	userID := uuid.New()
	layoutID := uuid.New()

	tokens, _, _, err := srv.CreateUserTokens(userID, layoutID, "CLIENT")
	if err != nil {
		t.Fatalf("CreateUserTokens() error = %v", err)
	}

	_, err = srv.RefreshTokens(context.Background(), tokens.Access, tokens.Refresh)
	if err != apperrors.TokenDontExist {
		t.Errorf("error = %v, want TokenDontExist", err)
	}
}
