package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sabirov8872/golang-rest-api/internal/service"
	"github.com/sabirov8872/golang-rest-api/internal/types"
)

type Handler struct {
	service service.IService
}

type IHandler interface {
	GetToken(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service service.IService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetToken(w http.ResponseWriter, r *http.Request) {
	var req types.GetToken
	json.NewDecoder(r.Body).Decode(&req)

	id, err := h.service.GetToken(req.Username)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusNoContent)
	}

	token, err := createToken(req.Username)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	writeJSON(w, types.TokenResponse{UserID: id, Token: token})
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) {
		errorResponse(w, "invalid token", http.StatusUnauthorized)
		return
	}

	res, err := h.service.GetAllUsers()
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	writeJSON(w, res)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := h.service.GetUserById(id)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusNoContent)
	}

	writeJSON(w, res)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	res, err := h.service.CreateUser(req)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	writeJSON(w, res)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req types.UpdateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := h.service.UpdateUser(id, req)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteUser(id)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusNoContent)
	}
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
	}

}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}
