package servermanager

import (
	"fmt"
	"net/http"
	"time"
)

type App struct {
	router http.Handler
	handle reqHandle
}

type reqHandle interface {
	Create(writer http.ResponseWriter, req *http.Request)
	List(writer http.ResponseWriter, req *http.Request)
	GetByID(writer http.ResponseWriter, req *http.Request)
	UpdateByID(writer http.ResponseWriter, req *http.Request)
	DeleteByID(writer http.ResponseWriter, req *http.Request)
}

func NewServer(handle reqHandle) *App {

	app := &App{
		handle: handle,
	}
	app.loadRoutes()
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
