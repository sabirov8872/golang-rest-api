package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sabirov8872/golang-rest-api/internal/service"
	"github.com/sabirov8872/golang-rest-api/internal/types"
)

type Handler struct {
	service   service.IService
	secretKey string
}

type IHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service service.IService, secretKey string) *Handler {
	return &Handler{
		service:   service,
		secretKey: secretKey,
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req types.SignInRequest
	json.NewDecoder(r.Body).Decode(&req)

	s, err := h.service.SignIn(req.Username)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, types.ErrorResponse{Message: "invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(req.Password))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, types.ErrorResponse{Message: "invalid username or password"})
		return
	}

	var token string
	token, err = createToken(req.Username, s.ID, h.secretKey)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, types.ErrorResponse{Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, types.SignInResponse{UserID: s.ID, Token: token})
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetAllUsers()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, types.ErrorResponse{Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	res, err := h.service.GetUserById(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, types.ErrorResponse{Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	hashedPassword, err := hashingPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, types.ErrorResponse{Message: err.Error()})
		return
	}
	req.Password = hashedPassword

	res, err := h.service.CreateUser(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, types.ErrorResponse{Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	hashedPassword, err := hashingPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, types.ErrorResponse{Message: err.Error()})
		return
	}
	req.Password = hashedPassword

	id := getID(r)
	err = h.service.UpdateUser(id, req)
	if err != nil {
		writeJSON(w, http.StatusNoContent, nil)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	err := h.service.DeleteUser(id)
	if err != nil {
		writeJSON(w, http.StatusNoContent, nil)
		return
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
	}
}

func getID(r *http.Request) string {
	return mux.Vars(r)["id"]
}

func createToken(username string, id int64, secretKey string) (string, error) {
	claims := &jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func hashingPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
