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
	CheckUser(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service service.IService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CheckUser(w http.ResponseWriter, r *http.Request) {
	var req types.GetToken
	json.NewDecoder(r.Body).Decode(&req)

	id, err := h.service.CheckUser(req.Username)
	if err != nil {
		writeJSON(w, http.StatusNoContent, types.ErrorResponse{Message: "invalid username"})
		return
	}

	token, err := createToken(req.Username)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, types.ErrorResponse{Message: "internal server error"})
		return
	}

	w.Header().Add("Authorization", "Bearer "+token)
	writeJSON(w, http.StatusOK, types.CheckUserResponse{UserID: id})
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "invalid token"})
		return
	}

	res, err := h.service.GetAllUsers()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, types.ErrorResponse{Message: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "invalid token"})
		return
	}

	id := getID(r)
	res, err := h.service.GetUserById(id)
	if err != nil {
		writeJSON(w, http.StatusNoContent, types.ErrorResponse{Message: "no content"})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "invalid token"})
		return
	}

	var req types.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	res, err := h.service.CreateUser(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, types.ErrorResponse{Message: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "invalid token"})
		return
	}

	var req types.UpdateUserRequest
	json.NewDecoder(r.Body).Decode(&req)
	id := getID(r)
	err := h.service.UpdateUser(id, req)
	if err != nil {
		writeJSON(w, http.StatusNoContent, types.ErrorResponse{Message: "no content"})
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "invalid token"})
		return
	}

	id := getID(r)
	err := h.service.DeleteUser(id)
	if err != nil {
		writeJSON(w, http.StatusNoContent, types.ErrorResponse{Message: "no content"})
		return
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func getID(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["id"]
}
