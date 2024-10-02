package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/sabirov8872/golang-rest-api/internal/handler"
	"github.com/sabirov8872/golang-rest-api/internal/types"
)

func Run(handler handler.IHandler, port string, secretKey string) {
	router := mux.NewRouter()
	router.HandleFunc("/sign_in", handler.SignIn).Methods("GET")
	router.HandleFunc("/user", AuthMiddleware(secretKey, handler.GetAllUsers)).Methods("GET")
	router.HandleFunc("/user/{id}", AuthMiddleware(secretKey, handler.GetUserById)).Methods("GET")
	router.HandleFunc("/user", AuthMiddleware(secretKey, handler.CreateUser)).Methods("POST")
	router.HandleFunc("/user/{id}", AuthMiddleware(secretKey, handler.UpdateUser)).Methods("PUT")
	router.HandleFunc("/user/{id}", AuthMiddleware(secretKey, handler.DeleteUser)).Methods("DELETE")

	log.Fatal(http.ListenAndServe("localhost:"+port, router))
}

func AuthMiddleware(secretKey string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		authHeader = authHeader[len("Bearer "):]
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secretKey), nil
		})
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: err.Error()})
			return
		}

		if !token.Valid {
			writeJSON(w, http.StatusUnauthorized, types.ErrorResponse{Message: "invalid token"})
			return
		}

		handler(w, r)
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
