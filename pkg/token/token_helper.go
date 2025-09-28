package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserId uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
	jwt.StandardClaims
}

func ParseTokenWithoutKeyCheck(accessToken string) (*CustomClaims, error) {
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(accessToken, &CustomClaims{})
	if err != nil {
		return nil, err
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("Cannot serialise claim to custom")
	}
	return claims, nil
}

func GetUserRole(claims *CustomClaims) string {
	return claims.Role
}

func GetUserId(claims *CustomClaims) uuid.UUID {
	return claims.UserId
}

func GetTokenId(claims *CustomClaims) string {
	return claims.Id
}
