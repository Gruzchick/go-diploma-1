package auth

import "github.com/golang-jwt/jwt/v4"

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID int64
}
