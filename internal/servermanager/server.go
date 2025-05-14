package servermanager

import (
	"fmt"
	"net/http"
	"time"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: loadRoutes(),
	}

	return app
}

func (app *App) Start() error {

	server := &http.Server{
		Addr:         ":3000",
		Handler:      app.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server start failed %w", err)
	}

	return nil
}
