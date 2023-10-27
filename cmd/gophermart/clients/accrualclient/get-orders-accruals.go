package accrualclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/configs"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	NEW        string = "NEW"
	REGISTERED string = "REGISTERED"
	INVALID    string = "INVALID"
	PROCESSING string = "PROCESSING"
	PROCESSED  string = "PROCESSED"
)

type Accrual struct {
	OrderID      string  `json:"order"`
	Status       string  `json:"status"`
	AccrualValue float64 `json:"accrual"`
}

type AccrualResponse struct {
	Accrual *Accrual
	Code    int
	Error   error
}

func GetOrdersAccruals(orderIDs []int64, accrualResponses *[]AccrualResponse, wg *sync.WaitGroup) {
	for i := 0; i < len(orderIDs); i++ {
		go GetOrderAccrual(orderIDs[i], accrualResponses, wg)
	}
}

func GetOrderAccrual(orderID int64, accrualResponses *[]AccrualResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	var m sync.Mutex

	orderIDString := strconv.FormatInt(orderID, 10)

	fmt.Println(configs.AccrualSystemAddress + "/api/orders/" + orderIDString)

	response, responseErr := http.Get(configs.AccrualSystemAddress + "/api/orders/" + orderIDString)

	if responseErr != nil {
		m.Lock()
		*accrualResponses = append(*accrualResponses, AccrualResponse{Error: responseErr})
		m.Unlock()
		return
	}

	if response.StatusCode == http.StatusInternalServerError {
		m.Lock()
		*accrualResponses = append(*accrualResponses, AccrualResponse{Code: http.StatusInternalServerError, Error: errors.New("ошибка системы расчёта начислений баллов")})
		m.Unlock()
		return
	}

	if response.StatusCode == http.StatusNoContent {
		stringOrderID := strconv.FormatInt(orderID, 10)

		m.Lock()
		*accrualResponses = append(*accrualResponses, AccrualResponse{Code: http.StatusNoContent, Accrual: &Accrual{OrderID: stringOrderID, AccrualValue: 0, Status: NEW}})
		m.Unlock()
		return
	}

	if response.StatusCode == http.StatusTooManyRequests {
		var repeatWg sync.WaitGroup

		delay, err := strconv.ParseInt(response.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			m.Lock()
			*accrualResponses = append(*accrualResponses, AccrualResponse{Error: err, Code: http.StatusTooManyRequests})
			m.Unlock()
			return
		}

		repeatWg.Add(1)

		time.AfterFunc(time.Duration(delay)*time.Second, func() { go GetOrderAccrual(orderID, accrualResponses, &repeatWg) })

		repeatWg.Wait()

		return
	}

	bodyBytes, readAllError := io.ReadAll(response.Body)
	response.Body.Close()
	if readAllError != nil {
		fmt.Println(readAllError)
		m.Lock()
		*accrualResponses = append(*accrualResponses, AccrualResponse{Error: readAllError})
		m.Unlock()
		return
	}

	var accrual Accrual

	if unmarshalErr := json.Unmarshal(bodyBytes, &accrual); unmarshalErr != nil {
		m.Lock()
		*accrualResponses = append(*accrualResponses, AccrualResponse{Error: unmarshalErr})
		m.Unlock()
		return
	}

	m.Lock()
	*accrualResponses = append(*accrualResponses, AccrualResponse{Accrual: &accrual, Code: response.StatusCode})
	m.Unlock()
}
