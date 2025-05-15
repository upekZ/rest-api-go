package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/upekZ/rest-api-go/internal/database/models"
	"github.com/upekZ/rest-api-go/internal/types"
	"net/http"
)

type DB interface {
	GetUserByID(context.Context, string) (*types.UserManager, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, *types.UserManager) error
	GetUsers(context.Context) ([]models.User, error)
	CreateUser(context.Context, *types.UserManager) error
}

type Handler struct {
	db DB
}

func NewHandler(dbMnger DB) *Handler {
	return &Handler{
		db: dbMnger,
	}
}

func (o *Handler) Create(writer http.ResponseWriter, req *http.Request) {

	var user types.UserManager

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, fmt.Errorf("decoding failure %w", err).Error(), http.StatusInternalServerError)
		return
	}

	if types.ValidateUser(&user) {
		err := o.db.CreateUser(req.Context(), &user)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusCreated)
		return
	}

	http.Error(writer, "invalid user params", http.StatusBadRequest)

}

func (o *Handler) List(writer http.ResponseWriter, req *http.Request) {

	users, err := o.db.GetUsers(req.Context())

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	if err := WriteJSON(writer, http.StatusAccepted, users); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (o *Handler) GetByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")
	user, err := o.db.GetUserByID(req.Context(), userID)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	if err := WriteJSON(writer, http.StatusCreated, user); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (o *Handler) UpdateByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")
	var user types.UserManager

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, fmt.Errorf("decoding failure %w", err).Error(), http.StatusInternalServerError)
		return
	}

	err := o.db.UpdateUser(req.Context(), userID, &user)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (o *Handler) DeleteByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")

	err := o.db.DeleteUser(req.Context(), userID)
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
