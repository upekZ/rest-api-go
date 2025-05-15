package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(writer http.ResponseWriter, reader *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	router.Route("/users", app.loadUserRoutes)

	app.router = router
}

func (app *App) loadUserRoutes(router chi.Router) {

	router.Post("/", app.handle.Create)
	router.Get("/", app.handle.List)
	router.Get("/{id}", app.handle.GetByID)
	router.Patch("/{id}", app.handle.UpdateByID)
	router.Delete("/{id}", app.handle.DeleteByID)

}
