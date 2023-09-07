package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func getTokenClaims(token string) (*TokenClaims, error) {
	tokenClaims := &TokenClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, tokenClaims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		fmt.Println("Token is not valid")
		return nil, err
	}

	return tokenClaims, nil
}
