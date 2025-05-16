package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/upekZ/rest-api-go/internal/types"
	"net/http"
	"time"
)

type Server struct {
	service   Service
	wsHandler WebSocketHandler
}

type WebSocketHandler interface {
	HandleWebSocket(w http.ResponseWriter, r *http.Request) error
	Broadcast(message []byte)
}

func NewServer(service Service, wsHandler WebSocketHandler) *Server {
	return &Server{
		service:   service,
		wsHandler: wsHandler,
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

func (app *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if err := app.wsHandler.HandleWebSocket(w, r); err != nil {
		http.Error(w, "Could not handle WebSocket", http.StatusBadRequest)
		return
	}
}

func (app *Server) broadcastUserEvent(eventType string, user types.UserEntity) {
	event := map[string]interface{}{
		"event": eventType,
		"user":  user,
	}
	data, err := json.Marshal(event)
	if err != nil {
		// Log error in production.
		return
	}
	app.wsHandler.Broadcast(data)
}
