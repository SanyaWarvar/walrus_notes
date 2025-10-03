package token

import (
	"errors"
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

func parseToken(secret, token string, skipExpireCheck bool) (*CustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	// Если нужно пропустить проверку Expire и есть ошибка
	if skipExpireCheck && err != nil {
		// В jwt/v5 проверяем через errors.Is
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Токен просрочен, но мы игнорируем эту ошибку
			// Проверяем, что токен в целом валиден (подпись и т.д.)
			if parsedToken != nil && parsedToken.Valid {
				claims, ok := parsedToken.Claims.(*CustomClaims)
				if ok {
					return claims, nil
				}
			}
		}
	}

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
