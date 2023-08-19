package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const ExpirationDuration = 5 * time.Minute

type Generator interface {
	Generate(login string, role string) (*Token, error)
}

type JwtGenerator struct {
	jwtKey []byte
}

func NewJwtGenerator(jwtKey []byte) *JwtGenerator {
	return &JwtGenerator{jwtKey: jwtKey}
}

func (j *JwtGenerator) Generate(login string, role string) (*Token, error) {
	expirationTime := time.Now().Add(ExpirationDuration)
	claims := &Claims{
		Login: login,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(j.jwtKey)

	if err != nil {
		return nil, err
	}

	return NewToken(tokenString), nil
}
