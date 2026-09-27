package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go-api-practice/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *service.UserService
}

type userRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := readUserID(w, r)
	if !ok {
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input userRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	user, err := h.service.Create(r.Context(), input.Name, input.Email)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Location", "/users/"+strconv.Itoa(user.ID))
	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := readUserID(w, r)
	if !ok {
		return
	}

	var input userRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	user, err := h.service.Update(r.Context(), id, input.Name, input.Email)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := readUserID(w, r)
	if !ok {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func readUserID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "id must be a positive integer"})
		return 0, false
	}

	return id, true
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrNameRequired),
		errors.Is(err, service.ErrNameTooLong),
		errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrEmailTooLong):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrEmailExists):
		writeJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
