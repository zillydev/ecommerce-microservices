package user

import (
	"database/sql"
	"ecommerce-microservices/pkg/kafka"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service *Service
	Kafka   *kafka.Kafka
}

func NewHandler(db *sql.DB, kafka *kafka.Kafka) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{Service: service, Kafka: kafka}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userRequest NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, fmt.Sprintf("error parsing register request: %v", err), http.StatusBadRequest)
		return
	}

	user, err := h.Service.RegisterUser(&userRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("error registering user: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	var updatePreferencesRequest UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&updatePreferencesRequest); err != nil {
		http.Error(w, fmt.Sprintf("error parsing update preferences request: %v", err), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	updatePreferencesRequest.UserID = userID

	if err := h.Service.UpdateUserPreferences(updatePreferencesRequest); err != nil {
		http.Error(w, fmt.Sprintf("error updating user preferences: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting user: %v", err), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUsersByPreferences(w http.ResponseWriter, r *http.Request) {
	var preferencesRequest GetUsersByPreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&preferencesRequest); err != nil {
		http.Error(w, fmt.Sprintf("error parsing preferences request: %v", err), http.StatusBadRequest)
		return
	}

	users, err := h.Service.GetUsersByPreferences(preferencesRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting users: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
