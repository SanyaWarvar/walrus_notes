package token

import (
	"context"
	"fmt"
	"time"
	apperrors "wn/internal/errors"
	"wn/internal/infrastructure/repository/tokens"
	"wn/pkg/util"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type tokenRepo interface {
	Create(ctx context.Context, token *tokens.RefreshToken) error
	GetByID(ctx context.Context, id uuid.UUID) (*tokens.RefreshToken, bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context, cutoffTime time.Time) error
}

func NewService(
	refreshTokenTTL time.Duration,
	accessTokenTTL time.Duration,
	secret string,
	tokenRepo tokenRepo,
) *Service {
	a := refreshTokenTTL.Seconds()
	b := accessTokenTTL.Seconds()
	_, _ = a, b
	fmt.Println(a, b)
	return &Service{
		refreshTokenTTL: refreshTokenTTL,
		accessTokenTTL:  accessTokenTTL,
		secret:          secret,
		tokenRepo:       tokenRepo,
	}
}

type Service struct {
	refreshTokenTTL time.Duration
	accessTokenTTL  time.Duration
	secret          string
	tokenRepo       tokenRepo
}

func (s *Service) CreateUserTokens(id uuid.UUID, role string) (*UserTokens, uuid.UUID, uuid.UUID, error) {
	jtiAccess := uuid.New()
	jtiRefresh := uuid.New()
	access, err := generateToken(jtiAccess, id, role, s.accessTokenTTL, s.secret)
	if err != nil {
		return nil, uuid.UUID{}, uuid.UUID{}, err
	}
	refresh, err := generateToken(jtiRefresh, id, role, s.refreshTokenTTL, s.secret)
	return &UserTokens{Access: access, Refresh: refresh}, jtiAccess, jtiRefresh, nil
}

func (s *Service) ParseToken(token string) (*CustomClaims, error) {
	return parseToken(s.secret, token)
}

func (s *Service) GenerateUserTokens(ctx context.Context, userId uuid.UUID, role string) (*UserTokens, error) {
	t, accessId, refreshId, err := s.CreateUserTokens(userId, role)
	if err != nil {
		return nil, errors.Wrap(err, ".GenerateUserTokens")
	}

	return t, s.tokenRepo.Create(ctx, &tokens.RefreshToken{
		Id:       refreshId,
		UserId:   userId,
		AccessId: accessId,
		ExpAt:    util.GetCurrentUTCTime().Add(s.refreshTokenTTL),
	})

}

func (s *Service) RefreshTokens(ctx context.Context, access, refresh string) (*UserTokens, error) {
	aToken, err := parseToken(s.secret, access)
	if err != nil {
		return nil, err
	}
	rToken, err := parseToken(s.secret, refresh)
	if err != nil {
		return nil, err
	}
	tokenId, err := uuid.Parse(rToken.ID)
	if err != nil {
		return nil, apperrors.TokenClaimsError
	}
	dbToken, ex, err := s.tokenRepo.GetByID(ctx, tokenId)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, apperrors.TokenDontExist
	}
	if dbToken.AccessId.String() != aToken.ID ||
		dbToken.Id != tokenId ||
		dbToken.UserId != aToken.UserId {
		return nil, apperrors.TokensDontMatch
	}
	t, err := s.GenerateUserTokens(ctx, aToken.UserId, aToken.Role)
	if err != nil {
		return nil, err
	}
	err = s.tokenRepo.Delete(ctx, dbToken.Id)
	if err != nil {
		return nil, err
	}
	return t, nil
}
