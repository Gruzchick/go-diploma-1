package router

import (
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/handlers"
	"github.com/go-chi/chi/v5"
)

func SetRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/user/register", handlers.RegisterHandler)
	router.Post("/api/user/login", handlers.LoginHandler)

	router.Post("/api/user/orders", auth.WithAuth(handlers.CreateOrderHandler))
	router.Get("/api/user/orders", auth.WithAuth(handlers.GetOrdersHandler))
	router.Get("/api/user/withdrawals", auth.WithAuth(handlers.GetWithdrawalsHandler))

	router.Get("/api/user/balance", auth.WithAuth(handlers.GetUserBalanceHandler))
	router.Post("/api/user/balance/withdraw", auth.WithAuth(handlers.WithdrawUserBalanceHandler))

	return router
}
