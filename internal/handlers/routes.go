package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *Server) loadChiRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(writer http.ResponseWriter, reader *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	router.Get("/ws", app.HandleWebSocket)
	router.Route("/users", app.loadUserRoutes)

	return router
}

func (app *Server) loadUserRoutes(router chi.Router) {

	router.Post("/", app.Create)
	router.Get("/", app.List)
	router.Get("/{id}", app.GetByID)
	router.Patch("/{id}", app.UpdateByID)
	router.Delete("/{id}", app.DeleteByID)
}
