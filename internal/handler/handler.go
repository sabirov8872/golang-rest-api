package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sabirov8872/golang-rest-api/internal/service"
	"github.com/sabirov8872/golang-rest-api/internal/types"
)

type Handler struct {
	service service.IService
}

type IHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service service.IService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req types.SignIn
	json.NewDecoder(r.Body).Decode(&req)

	id, err := h.service.SignIn(req)
	if err != nil {
		writeJSON(w, http.StatusNoContent, types.ErrorResponse{Message: "invalid username"})
		return
	}

	if r.Header.Get("Authorization") == "" || checkAuth(r) != nil {
		var token string
		token, err = createToken(req.Username, id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, types.ErrorResponse{Message: "internal server error"})
			return
		}

		w.Header().Add("Authorization", "Bearer "+token)
	}

	writeJSON(w, http.StatusOK, types.CheckUserResponse{UserID: id})
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if err := checkAuth(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "unauthorized"})
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
	if err := checkAuth(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "unauthorized"})
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
	if err := checkAuth(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "unauthorized"})
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
	if err := checkAuth(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "unauthorized"})
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

func checkAuth(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return errors.New("unauthorized")
	}

	authHeader = authHeader[len("Bearer "):]
	secret := os.Getenv("SECRET_KEY")
	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func createToken(username string, id int64) (string, error) {
	claims := &jwt.MapClaims{
		"username": username,
		"id":       id,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("SECRET_KEY")
	return token.SignedString([]byte(secret))
}
