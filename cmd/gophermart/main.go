package main

import (
	"fmt"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/configs"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/dbs/diploma"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/router"
	"net/http"
)

func main() {
	if err := configs.Configure(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := diploma.Init(configs.DatabaseURI); err != nil {
		fmt.Println(err)
	}

	server := &http.Server{
		Addr:    configs.RunAddress,
		Handler: router.SetRoutes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
