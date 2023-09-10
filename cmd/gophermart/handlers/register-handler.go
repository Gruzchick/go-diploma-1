package handlers

import (
	"database/sql"
	"errors"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/dbs/diplomadb"
	"net/http"
)

type requestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterHandler(res http.ResponseWriter, req *http.Request) {
	var unmarshalledBody requestBody

	if err := unmarshalBody(req.Body, &unmarshalledBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateRegisterHandlerBody(unmarshalledBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	queryRow := diplomadb.DB.QueryRow(`
	SELECT id FROM users where login = $1
	`, unmarshalledBody.Login)

	var id int64

	queryRowError := queryRow.Scan(&id)

	if queryRowError != nil && !errors.Is(queryRowError, sql.ErrNoRows) {
		http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
		return
	}

	if queryRowError == nil {
		http.Error(res, "логин уже занят", http.StatusConflict)
		return
	}

	insertQueryRow := diplomadb.DB.QueryRow(`
	INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id
	`, unmarshalledBody.Login, unmarshalledBody.Password)

	var newID int64 = 1

	insertQueryRowError := insertQueryRow.Scan(&newID)

	if insertQueryRowError != nil {
		http.Error(res, insertQueryRowError.Error(), http.StatusInternalServerError)
		return
	}

	newJwt, jwtError := auth.CreateJwtToken(newID)
	if jwtError != nil {
		http.Error(res, jwtError.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set(auth.AuthHeader, newJwt)
	res.WriteHeader(http.StatusOK)
}

func validateRegisterHandlerBody(body requestBody) error {
	if len(body.Login) == 0 {
		return errors.New("не указан логин")
	}

	if len(body.Password) == 0 {
		return errors.New("не указан пароль")
	}

	return nil
}
