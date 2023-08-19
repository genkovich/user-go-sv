package token

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Login string `json:"login"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}
