package handlers

import (
	"database/sql"
	"errors"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/dbs/diplomadb"
	"net/http"
)

func CreateOrderHandler(res http.ResponseWriter, req *http.Request) {
	tokenClaims := req.Context().Value(auth.TokenClaimsContextFieldName).(*auth.TokenClaims)

	var requestOrderID int64

	if err := unmarshalBody(req.Body, &requestOrderID); err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	var queryUserID int64

	queryRow := diplomadb.DB.QueryRow(`
	SELECT userId FROM orders where id = $1
	`, requestOrderID)

	queryRowError := queryRow.Scan(&queryUserID)
	if queryRowError != nil && !errors.Is(queryRowError, sql.ErrNoRows) {
		http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
		return
	} else if queryRowError == nil {
		if queryUserID == tokenClaims.UserID {
			res.WriteHeader(http.StatusOK)
			return
		} else {
			http.Error(res, "номер заказа уже был загружен другим пользователем", http.StatusConflict)
			return
		}
	} else {
		_, insertError := diplomadb.DB.Exec(`INSERT INTO orders (id, userId) values ($1, $2)`, requestOrderID, tokenClaims.UserID)
		if insertError != nil {
			http.Error(res, insertError.Error(), http.StatusInternalServerError)
			return
		}
	}

	res.WriteHeader(http.StatusAccepted)
}
