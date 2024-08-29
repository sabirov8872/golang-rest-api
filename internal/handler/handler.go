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
	var req types.GetTokenResponse
	json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	id, err := h.service.GetToken(req)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusNoContent)
	}

	token, err := createToken(req.Username)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	jsonResponse(w, types.TokenResponse{UserID: id, Token: token})

}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		newErrorResponse(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err := verifyToken(tokenString)
	if err != nil {
		newErrorResponse(w, "Invalid token", http.StatusUnauthorized)
	}

	res, err := h.service.GetAllUsers()
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	jsonResponse(w, res)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := h.service.GetUserById(id)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusNoContent)
	}

	jsonResponse(w, res)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	res, err := h.service.CreateUser(req)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	jsonResponse(w, res)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req types.UpdateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := h.service.UpdateUser(id, req)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusNoContent)
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteUser(id)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusNoContent)
	}
}
