package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	service Service
	channel Channel
}

func NewServer(service Service, channel Channel) *Server {
	return &Server{
		service: service,
		channel: channel,
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
