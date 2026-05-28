//go:build integration

package integration

import (
	"testing"
	"time"
	"wn/internal/application/auth"
	appuser "wn/internal/application/user"
	"wn/internal/domain/dto"
	domainlayout "wn/internal/domain/services/layout"
	domaintoken "wn/internal/domain/services/token"
	domainuser "wn/internal/domain/services/user"
	layoutrepo "wn/internal/infrastructure/repository/layout"
	permrepo "wn/internal/infrastructure/repository/permissions"
	tokenrepo "wn/internal/infrastructure/repository/tokens"
	userrepo "wn/internal/infrastructure/repository/user"

	"github.com/google/uuid"
)

func TestAuthFlow_RegisterAndLogin(t *testing.T) {
	resetDB(t)

	userRepo := userrepo.NewRepository(env.DB)
	layoutRepo := layoutrepo.NewRepository(env.DB)
	tokenRepo := tokenrepo.NewRepository(env.DB)
	permRepo := permrepo.NewRepository(env.DB)

	domainUser := domainuser.NewService("integration-salt", env.Trx, noopLogger{}, userRepo)
	domainLayout := domainlayout.NewService(env.Trx, noopLogger{}, layoutRepo, nil, nil, nil, nil, permRepo)
	tokenSvc := domaintoken.NewService(time.Hour*24, time.Minute*15, "integration-jwt-secret", tokenRepo)

	userApp := appuser.NewService(env.Trx, noopLogger{}, domainUser, nil, domainLayout)
	authApp := auth.NewService(env.Trx, noopLogger{}, domainUser, nil, tokenSvc, domainLayout)

	credentials := dto.RegisterCredentials{
		Username: "integration_user",
		Email:    "integration_user@test.com",
		Password: "SecurePass123",
	}

	reg, err := userApp.RegisterUser(env.Ctx, credentials)
	if err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}
	if reg.UserId == uuid.Nil {
		t.Fatal("expected user id")
	}

	layouts, err := layoutRepo.GetAvailableLayouts(env.Ctx, reg.UserId)
	if err != nil {
		t.Fatalf("GetAvailableLayouts() error = %v", err)
	}
	if len(layouts) != 1 {
		t.Fatalf("expected 1 main layout, got %d", len(layouts))
	}
	if !layouts[0].IsMain {
		t.Error("layout should be main")
	}

	login, err := authApp.Login(env.Ctx, dto.LoginRequest{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if login.Access == "" || login.Refresh == "" {
		t.Fatal("expected access and refresh tokens")
	}
	if login.UserId != reg.UserId {
		t.Errorf("UserId = %v, want %v", login.UserId, reg.UserId)
	}

	refreshed, err := authApp.RefreshTokens(env.Ctx, domaintoken.UserTokens{
		Access:  login.Access,
		Refresh: login.Refresh,
	})
	if err != nil {
		t.Fatalf("RefreshTokens() error = %v", err)
	}
	if refreshed.Access == "" || refreshed.Refresh == "" {
		t.Fatal("expected refreshed tokens")
	}
	if refreshed.Access == login.Access {
		t.Error("access token should change after refresh")
	}
}

func TestAuthFlow_Login_WrongPassword(t *testing.T) {
	resetDB(t)

	userRepo := userrepo.NewRepository(env.DB)
	layoutRepo := layoutrepo.NewRepository(env.DB)
	tokenRepo := tokenrepo.NewRepository(env.DB)
	permRepo := permrepo.NewRepository(env.DB)

	domainUser := domainuser.NewService("integration-salt", env.Trx, noopLogger{}, userRepo)
	domainLayout := domainlayout.NewService(env.Trx, noopLogger{}, layoutRepo, nil, nil, nil, nil, permRepo)
	tokenSvc := domaintoken.NewService(time.Hour*24, time.Minute*15, "integration-jwt-secret", tokenRepo)

	userApp := appuser.NewService(env.Trx, noopLogger{}, domainUser, nil, domainLayout)
	authApp := auth.NewService(env.Trx, noopLogger{}, domainUser, nil, tokenSvc, domainLayout)

	credentials := dto.RegisterCredentials{
		Username: "wrong_pass_user",
		Email:    "wrong_pass@test.com",
		Password: "CorrectPassword",
	}
	if _, err := userApp.RegisterUser(env.Ctx, credentials); err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}

	_, err := authApp.Login(env.Ctx, dto.LoginRequest{
		Email:    credentials.Email,
		Password: "WrongPassword",
	})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}
