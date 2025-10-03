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

func parseToken(secret, token string, skipExpireCheck bool) (*CustomClaims, error) {
	// Создаем кастомный парсер с нужными настройками
	parser := jwt.NewParser()

	if skipExpireCheck {
		// Отключаем проверку времени для этого случая
		parser = jwt.NewParser(jwt.WithoutClaimsValidation())
	}

	parsedToken, err := parser.ParseWithClaims(
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

	// При отключенной проверке времени, вручную проверяем подпись
	if skipExpireCheck {
		if !parsedToken.Valid {
			return nil, fmt.Errorf("token signature is not valid")
		}
	} else {
		if !parsedToken.Valid {
			return nil, fmt.Errorf("token is not valid")
		}
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
