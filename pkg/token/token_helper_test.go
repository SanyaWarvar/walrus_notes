package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func buildTestToken(t *testing.T, claims CustomClaims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}
	return signed
}

func TestParseTokenWithoutKeyCheck(t *testing.T) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	layoutID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	tokenID := uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")

	claims := CustomClaims{
		UserId:       userID,
		Role:         "ADMIN",
		MainLayoutId: layoutID,
		StandardClaims: jwt.StandardClaims{
			Id:        tokenID.String(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}

	tokenStr := buildTestToken(t, claims)

	parsed, err := ParseTokenWithoutKeyCheck(tokenStr)
	if err != nil {
		t.Fatalf("ParseTokenWithoutKeyCheck() error = %v", err)
	}

	if parsed.UserId != userID {
		t.Errorf("UserId = %v, want %v", parsed.UserId, userID)
	}
	if parsed.Role != "ADMIN" {
		t.Errorf("Role = %q, want %q", parsed.Role, "ADMIN")
	}
	if parsed.MainLayoutId != layoutID {
		t.Errorf("MainLayoutId = %v, want %v", parsed.MainLayoutId, layoutID)
	}
	if parsed.Id != tokenID.String() {
		t.Errorf("Id = %q, want %q", parsed.Id, tokenID.String())
	}
}

func TestParseTokenWithoutKeyCheck_InvalidToken(t *testing.T) {
	_, err := ParseTokenWithoutKeyCheck("not.a.jwt")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestClaimsGetters(t *testing.T) {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	layoutID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	claims := &CustomClaims{
		UserId:       userID,
		Role:         "CLIENT",
		MainLayoutId: layoutID,
		StandardClaims: jwt.StandardClaims{
			Id: "token-id",
		},
	}

	if GetUserRole(claims) != "CLIENT" {
		t.Errorf("GetUserRole() = %q, want %q", GetUserRole(claims), "CLIENT")
	}
	if GetUserId(claims) != userID {
		t.Errorf("GetUserId() = %v, want %v", GetUserId(claims), userID)
	}
	if GetTokenId(claims) != "token-id" {
		t.Errorf("GetTokenId() = %q, want %q", GetTokenId(claims), "token-id")
	}
	if GetMainLayoutId(claims) != layoutID {
		t.Errorf("GetMainLayoutId() = %v, want %v", GetMainLayoutId(claims), layoutID)
	}
}
