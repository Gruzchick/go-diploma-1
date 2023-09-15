package handlers

import (
	"database/sql"
	"errors"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/dbs/diplomadb"
	"net/http"
)

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	var unmarshalledBody credentialsBody

	if err := unmarshalBody(req.Body, &unmarshalledBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateLoginHandlerBody(unmarshalledBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	queryRow := diplomadb.DB.QueryRow(`
		SELECT id, password FROM users where login = $1
	`, unmarshalledBody.Login)

	var userID int64
	var password string

	queryRowError := queryRow.Scan(&userID, &password)

	if queryRowError != nil && !errors.Is(queryRowError, sql.ErrNoRows) {
		http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
		return
	}

	if errors.Is(queryRowError, sql.ErrNoRows) {
		http.Error(res, "неверная пара логин/пароль", http.StatusUnauthorized)
		return
	}

	if password != unmarshalledBody.Password {
		http.Error(res, "неверная пара логин/пароль", http.StatusUnauthorized)
		return
	}

	newJwt, jwtError := auth.CreateJwtToken(userID)
	if jwtError != nil {
		http.Error(res, jwtError.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set(auth.AuthHeader, newJwt)
	res.WriteHeader(http.StatusOK)
}

func validateLoginHandlerBody(body credentialsBody) error {
	if len(body.Login) == 0 {
		return errors.New("не указан логин")
	}

	if len(body.Password) == 0 {
		return errors.New("не указан пароль")
	}

	return nil
}
