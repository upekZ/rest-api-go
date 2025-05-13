package servermanager

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/upekZ/rest-api-go/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(writer http.ResponseWriter, reader *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	router.Route("/users", loadUserRoutes)

	return router

}

func loadUserRoutes(router chi.Router) {
	userHandler, err := handler.NewHandler()

	if err != nil {
		fmt.Println("Failure to load user handler: %w", err)
	}

	router.Post("/", userHandler.Create)
	router.Get("/", userHandler.List)
	router.Get("/{id}", userHandler.GetByID)
	router.Patch("/{id}", userHandler.UpdateByID)
	router.Delete("/{id}", userHandler.DeleteByID)

}
