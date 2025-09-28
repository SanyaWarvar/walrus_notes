package token

import (
	"fmt"
	"time"

	apperrors "wn/internal/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserTokens
// @schema
type UserTokens struct {
	Access  string `json:"accessToken" binding:"required"`
	Refresh string `json:"refreshToken" binding:"required"`
}

type CustomClaims struct {
	UserId uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func parseToken(secret, token string) (*CustomClaims, error) {

	parsedToken, err := jwt.ParseWithClaims(
		token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok {
		return nil, apperrors.TokenClaimsError
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

func generateToken(id, userId uuid.UUID, role string, ttl time.Duration, secret string) (string, error) {
	claims := CustomClaims{
		userId, role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Subject:   userId.String(),
			ID:        id.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(secret))
}
