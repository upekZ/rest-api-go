package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	service Service
}

func NewServer(service Service) *Server {
	return &Server{
		service: service,
	}
}

func (app *Server) Start() error {

	server := &http.Server{
		Addr:         ":3000",
		Handler:      app.loadChiRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server start failure: %w", err)
	}

	return err
}
