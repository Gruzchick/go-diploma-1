package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/clients/accrualclient"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/dbs/diplomadb"
	"github.com/theplant/luhn"
	"net/http"
	"sync"
)

func CreateOrderHandler(res http.ResponseWriter, req *http.Request) {
	tokenClaims := req.Context().Value(auth.TokenClaimsContextFieldName).(*auth.TokenClaims)

	var requestOrderID int64

	if err := unmarshalBody(req.Body, &requestOrderID); err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if !validateCreateOrderHandlerRequest(int(requestOrderID)) {
		fmt.Println(`create-order-handler.go:23`, `"не валидное значение ордера"`, "не валидное значение ордера")
		http.Error(res, "не валидное значение ордера", http.StatusUnprocessableEntity)
		return
	}

	var wg sync.WaitGroup

	var accruals = make([]accrualclient.AccrualResponse, 0)

	wg.Add(1)

	go accrualclient.GetOrdersAccruals([]int64{requestOrderID}, &accruals, &wg)

	wg.Wait()

	fmt.Println(`get-orders-handler.go:68`, `accruals`, accruals)

	var queryUserID int64

	if accruals[0].Accrual == nil {
		http.Error(res, "Не верный формат запроса", http.StatusUnprocessableEntity)
		return
	}

	queryRow := diplomadb.DB.QueryRow(`
		SELECT userId FROM orders where id = $1
	`, requestOrderID)

	queryRowError := queryRow.Scan(&queryUserID)
	if queryRowError != nil && !errors.Is(queryRowError, sql.ErrNoRows) {
		fmt.Println(`create-order-handler.go:37`, "here-1")

		http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
		return
	} else if queryRowError == nil {
		fmt.Println(`create-order-handler.go:37`, "here-2")

		if queryUserID == tokenClaims.UserID {
			res.WriteHeader(http.StatusOK)
			return
		} else {
			fmt.Println(`create-order-handler.go:37`, "here-3")

			http.Error(res, "номер заказа уже был загружен другим пользователем", http.StatusConflict)
			return
		}
	} else {
		fmt.Println(`create-order-handler.go:37`, "here-4")

		_, insertError := diplomadb.DB.Exec(`
			INSERT INTO orders (id, userId) values ($1, $2)
		`, requestOrderID, tokenClaims.UserID)
		if insertError != nil {
			fmt.Println(`create-order-handler.go:37`, "here-5")
			http.Error(res, insertError.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Println(`create-order-handler.go:37`, "here-6")
	res.WriteHeader(http.StatusAccepted)
}

func validateCreateOrderHandlerRequest(orderID int) bool {
	if luhn.Valid(orderID) {
		return true
	} else {
		return false
	}
}
