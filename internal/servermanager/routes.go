package servermanager

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/upekZ/rest-api-go/api/handler"
	"net/http"
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
		fmt.Printf("Failure to load user handler: %s\n", err.Error())
	}

	router.Post("/", userHandler.Create)
	router.Get("/", userHandler.List)
	router.Get("/{id}", userHandler.GetByID)
	router.Patch("/{id}", userHandler.UpdateByID)
	router.Delete("/{id}", userHandler.DeleteByID)

}
