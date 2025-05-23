package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/upekZ/rest-api-go/internal/database/queries" //To be removed after moving usage of queries.User --> model.UserEntity
	"github.com/upekZ/rest-api-go/internal/model"
	"net/http"
)

type Service interface {
	CreateUser(ctx context.Context, user *model.UserEntity) error
	ListUsers(ctx context.Context) ([]queries.User, error) //queries.User to be replaced with model.UserEntity
	GetUserByID(ctx context.Context, id string) (*model.UserEntity, error)
	DeleteUser(ctx context.Context, id string) (*model.UserEntity, error)
	UpdateUser(ctx context.Context, id string, user *model.UserEntity) (*model.UserEntity, error)
	HandleWebSocket(w http.ResponseWriter, r *http.Request)
}

type Channel interface {
	ServeWS() http.HandlerFunc
}

func (app *Server) Create(writer http.ResponseWriter, req *http.Request) {

	var user model.UserEntity

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, "users creation failure", http.StatusInternalServerError)
		return
	}

	if err := app.service.CreateUser(req.Context(), &user); err != nil {
		http.Error(writer, "users creation failure", http.StatusInternalServerError)
		return
	}

	if err := WriteJSON(writer, http.StatusCreated, user); err != nil {
		http.Error(writer, "users creation failure", http.StatusInternalServerError)
	}
}

func (app *Server) List(writer http.ResponseWriter, req *http.Request) {

	users, err := app.service.ListUsers(req.Context())

	if err != nil {
		http.Error(writer, "no users found", http.StatusInternalServerError)
		return
	}

	if err := WriteJSON(writer, http.StatusOK, users); err != nil {
		http.Error(writer, "no users found", http.StatusInternalServerError)
	}
}

func (app *Server) GetByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")
	user, err := app.service.GetUserByID(req.Context(), userID)

	if err != nil {
		http.Error(writer, "user not found", http.StatusNotFound)
		return
	}

	if err := WriteJSON(writer, http.StatusCreated, user); err != nil {
		http.Error(writer, "user not found", http.StatusInternalServerError)
	}
}

func (app *Server) UpdateByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")
	var user model.UserEntity

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, "user not found", http.StatusInternalServerError)
		return
	}

	_, err := app.service.UpdateUser(req.Context(), userID, &user)

	if err != nil {
		http.Error(writer, "users update failure", http.StatusInternalServerError)
		return
	}

	if err := WriteJSON(writer, http.StatusOK, user); err != nil {
		http.Error(writer, "users update failure", http.StatusInternalServerError)
	}
}

func (app *Server) DeleteByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")

	user, err := app.service.DeleteUser(req.Context(), userID)
	if err != nil {
		http.Error(writer, "users deletion failure", http.StatusInternalServerError)
		return
	}

	if err := WriteJSON(writer, http.StatusOK, user); err != nil {
		http.Error(writer, "users update failure", http.StatusInternalServerError)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
