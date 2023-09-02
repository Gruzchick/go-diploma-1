package router

import (
	"github.com/Gruzchick/go-diploma-1/cmd/gophermart/handlers"
	"github.com/go-chi/chi/v5"
)

func SetRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/user/register", handlers.RegisterHandler)

	return router
}
