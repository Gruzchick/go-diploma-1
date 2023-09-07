package auth

import (
	"context"
	"net/http"
	"strings"
)

func WithAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token := strings.Replace(req.Header.Get(AuthHeader), "Bearer ", "", 1)

		tokenClaims, tokenClaimsErr := getTokenClaims(token)
		if tokenClaimsErr != nil {
			http.Error(res, "ошибка авторизации", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), TokenClaimsContextFieldName, tokenClaims)

		h.ServeHTTP(res, req.WithContext(ctx))
	}
}
