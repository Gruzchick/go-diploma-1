package handlers

import (
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

	res.WriteHeader(http.StatusOK)
}
