package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/upekZ/rest-api-go/internal/database/queries" //To be removed after moving usage of queries.User --> types.UserEntity
	"github.com/upekZ/rest-api-go/internal/types"
	"log"
	"net/http"
)

type Service interface {
	CreateUser(ctx context.Context, user types.UserEntity) error
	ListUsers(ctx context.Context) ([]queries.User, error) //queries.User to be replaced with types.UserEntity
	GetUserByID(ctx context.Context, id string) (*types.UserEntity, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, user *types.UserEntity) error
}

type Channel interface {
	ServeWS() http.HandlerFunc
}

func (app *Server) Create(writer http.ResponseWriter, req *http.Request) {

	var user types.UserEntity

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, fmt.Errorf("decoding failure %w", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := app.service.CreateUser(req.Context(), user); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(writer).Encode(map[string]string{"message": "User created"}); err != nil {
		log.Printf("Failed to write JSON response in user creation: %v", err)
	}

	app.broadcastUserEvent("created", user)
	writer.WriteHeader(http.StatusCreated)

}

func (app *Server) List(writer http.ResponseWriter, req *http.Request) {

	users, err := app.service.ListUsers(req.Context())

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	if err := WriteJSON(writer, http.StatusAccepted, users); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Server) GetByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")
	user, err := app.service.GetUserByID(req.Context(), userID)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	if err := WriteJSON(writer, http.StatusCreated, user); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Server) UpdateByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")
	var user types.UserEntity

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, fmt.Errorf("decoding failure %w", err).Error(), http.StatusInternalServerError)
		return
	}

	err := app.service.UpdateUser(req.Context(), userID, &user)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (app *Server) DeleteByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")

	err := app.service.DeleteUser(req.Context(), userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
