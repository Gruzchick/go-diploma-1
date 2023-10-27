package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/clients/accrualclient"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/dbs/diplomadb"
	"net/http"
	"sync"
)

type OrderDTO struct {
	Number  string  `json:"number"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func GetOrdersHandler(res http.ResponseWriter, req *http.Request) {
	tokenClaims := req.Context().Value(auth.TokenClaimsContextFieldName).(*auth.TokenClaims)

	queryRows, queryRowError := diplomadb.DB.Query(`
		SELECT id FROM orders where userId = $1
	`, tokenClaims.UserID)
	if queryRowError != nil && !errors.Is(queryRowError, sql.ErrNoRows) {
		http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
		return
	}
	defer queryRows.Close()

	if errors.Is(queryRowError, sql.ErrNoRows) {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	var orderIDs = make([]int64, 0)

	for queryRows.Next() {
		var orderID int64

		err := queryRows.Scan(&orderID)
		if err != nil {
			http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
			return
		}

		orderIDs = append(orderIDs, orderID)
	}

	rowsError := queryRows.Err()
	if rowsError != nil {
		http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup

	var accruals = make([]accrualclient.AccrualResponse, 0)

	wg.Add(len(orderIDs))

	go accrualclient.GetOrdersAccruals(orderIDs, &accruals, &wg)

	wg.Wait()

	responseOrders := make([]OrderDTO, 0)

	for _, accrual := range accruals {
		if accrual.Error != nil {
			http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
			return
		}

		if accrual.Accrual != nil {
			responseOrders = append(responseOrders, OrderDTO{
				Number:  accrual.Accrual.OrderID,
				Accrual: accrual.Accrual.AccrualValue,
				Status:  accrual.Accrual.Status,
			})
		}
	}

	marshaledResp, err := json.Marshal(responseOrders)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")

	if len(responseOrders) == 0 {
		res.WriteHeader(http.StatusNoContent)
	} else {
		res.WriteHeader(http.StatusOK)
	}

	res.Write(marshaledResp)
}
