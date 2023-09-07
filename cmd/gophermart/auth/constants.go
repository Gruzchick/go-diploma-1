package auth

import "time"

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"
const AuthHeader = "Authorization"

type contextKey string

const TokenClaimsContextFieldName contextKey = "tokenClaims"
