package diplomadb

import (
	"database/sql"
	"errors"
)

type WithdrawalRow struct {
	Id     int64
	userID int64
	Sum    float64
	Order  string
}

func GetWithdrawalsByUserId(userID int64) (*[]WithdrawalRow, error) {
	queryRows, queryRowError := DB.Query(`
	SELECT id, sum, orderid FROM withdrawals where userId = $1
	`, userID)
	if queryRowError != nil && !errors.Is(queryRowError, sql.ErrNoRows) {
		return nil, queryRowError
	}
	defer queryRows.Close()

	var withdrawalRows = make([]WithdrawalRow, 0)

	if errors.Is(queryRowError, sql.ErrNoRows) {
		return &withdrawalRows, nil
	}

	for queryRows.Next() {
		var r WithdrawalRow

		err := queryRows.Scan(&r.Id, &r.Sum, &r.Order)
		if err != nil {
			return nil, err
		}

		withdrawalRows = append(withdrawalRows, r)
	}

	rowsError := queryRows.Err()
	if rowsError != nil {
		return nil, rowsError
	}

	return &withdrawalRows, nil
}
