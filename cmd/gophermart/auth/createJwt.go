package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

func CreateJwt(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
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
