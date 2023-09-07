package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateJwtToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},

		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
