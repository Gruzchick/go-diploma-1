package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/clients/accrualclient"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/dbs/diplomadb"
	"net/http"
	"sync"
)

type WithdrawUserBalanceHandlerRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func WithdrawUserBalanceHandler(res http.ResponseWriter, req *http.Request) {
	tokenClaims := req.Context().Value(auth.TokenClaimsContextFieldName).(*auth.TokenClaims)

	var unmarshalledBody WithdrawUserBalanceHandlerRequest

	if err := unmarshalBody(req.Body, &unmarshalledBody); err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	queryRows, queryRowError := diplomadb.DB.Query(`
	SELECT id FROM orders where userId = $1
	`, tokenClaims.UserID)
	if queryRowError != nil && !errors.Is(queryRowError, sql.ErrNoRows) {
		fmt.Println(queryRowError.Error())
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
			fmt.Println(queryRowError.Error())
			http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
			return
		}

		orderIDs = append(orderIDs, orderID)
	}

	rowsError := queryRows.Err()
	if rowsError != nil {
		fmt.Println(queryRowError.Error())
		http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup

	var accruals = make([]accrualclient.AccrualResponse, 0) // TODO ASK: Почему не получилось сделать через канал

	wg.Add(len(orderIDs)) // TODO ASK: Как занести это внутрь accrualclient.GetOrdersAccruals()

	go accrualclient.GetOrdersAccruals(orderIDs, &accruals, &wg)

	wg.Wait()

	var current float64 = 0

	for _, accrual := range accruals {
		if accrual.Error != nil {
			fmt.Println(queryRowError.Error())
			http.Error(res, queryRowError.Error(), http.StatusInternalServerError)
			return
		}

		if accrual.Accrual != nil && accrual.Accrual.Status == accrualclient.PROCESSED {
			current = current + accrual.Accrual.AccrualValue
		}
	}

	withdrawals, withdrawalsErrors := diplomadb.GetWithdrawalsByUserId(tokenClaims.UserID)
	if withdrawalsErrors != nil {
		fmt.Println(withdrawalsErrors.Error())
		http.Error(res, withdrawalsErrors.Error(), http.StatusInternalServerError)
		return
	}

	var withdrawalsSum float64

	for _, v := range *withdrawals {
		withdrawalsSum = withdrawalsSum + v.Sum
	}

	current = current - withdrawalsSum

	if unmarshalledBody.Sum > current {
		http.Error(res, queryRowError.Error(), http.StatusPaymentRequired)
		return
	}

	_, insertError := diplomadb.DB.Exec(`INSERT INTO withdrawals (userId, sum, orderid) values ($1, $2, $3)`, tokenClaims.UserID, unmarshalledBody.Sum, unmarshalledBody.Order)
	if insertError != nil {
		fmt.Println(insertError.Error())
		http.Error(res, insertError.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
