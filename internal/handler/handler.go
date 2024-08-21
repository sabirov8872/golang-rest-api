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
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service service.IService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetAllUsers()
	if err != nil {
		statusInternalServerError(w, err)
	}

	jsonResponse(w, res)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := h.service.GetUserById(id)
	if err != nil {
		statusInternalServerError(w, err)
	}

	jsonResponse(w, res)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	res, err := h.service.CreateUser(req)
	if err != nil {
		statusInternalServerError(w, err)
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
		statusInternalServerError(w, err)
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteUser(id)
	if err != nil {
		statusInternalServerError(w, err)
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		statusInternalServerError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		statusInternalServerError(w, err)
	}
}

func statusInternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
