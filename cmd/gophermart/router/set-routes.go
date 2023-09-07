package router

import (
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/auth"
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/handlers"
	"github.com/go-chi/chi/v5"
)

func SetRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/user/register", handlers.RegisterHandler)
	router.Post("/api/user/orders", auth.WithAuth(handlers.PostOrdersHandler))

	return router
}
