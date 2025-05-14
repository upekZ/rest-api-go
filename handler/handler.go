package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/upekZ/rest-api-go/datamanager"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	dbConnector *datamanager.PostgresConn
}

func NewHandler() (*Handler, error) {

	pgConnector, err := datamanager.NewPostgresConn()

	if err != nil {
		return nil, err
	}

	return &Handler{
		dbConnector: pgConnector,
	}, err
}

func (o *Handler) Create(writer http.ResponseWriter, req *http.Request) {

	var user datamanager.UserManager

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, fmt.Errorf("decoding failure %w", err).Error(), http.StatusInternalServerError)
		return
	}

	if datamanager.ValidateUser(&user) {
		err := o.dbConnector.CreateUser(req.Context(), &user)

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

	users, err := o.dbConnector.GetUsers(req.Context())

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	if err := WriteJSON(writer, http.StatusAccepted, users); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (o *Handler) GetByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")
	user, err := o.dbConnector.GetUserByID(req.Context(), userID)

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
	var user datamanager.UserManager

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, fmt.Errorf("decoding failure %w", err).Error(), http.StatusInternalServerError)
		return
	}

	err := o.dbConnector.UpdateUser(req.Context(), userID, &user)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (o *Handler) DeleteByID(writer http.ResponseWriter, req *http.Request) {

	userID := chi.URLParam(req, "id")

	err := o.dbConnector.DeleteUser(req.Context(), userID)
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
