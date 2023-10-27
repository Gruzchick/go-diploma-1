package handlers

import (
	"encoding/json"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/dbs/diplomadb"
	"net/http"
)

type WithdrawalDTO struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func GetWithdrawalsHandler(res http.ResponseWriter, req *http.Request) {
	tokenClaims := req.Context().Value(auth.TokenClaimsContextFieldName).(*auth.TokenClaims)

	withdrawals, withdrawalsErrors := diplomadb.GetWithdrawalsByUserID(tokenClaims.UserID)
	if withdrawalsErrors != nil {
		http.Error(res, withdrawalsErrors.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]WithdrawalDTO, 0)

	for _, v := range *withdrawals {
		response = append(response, WithdrawalDTO{Order: v.Order, Sum: v.Sum})
	}

	marshaledResp, err := json.Marshal(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")

	if len(response) == 0 {
		res.WriteHeader(http.StatusNoContent)
	} else {
		res.WriteHeader(http.StatusOK)
	}

	res.Write(marshaledResp)
}
