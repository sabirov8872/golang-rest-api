package routes

import (
	"github.com/gorilla/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sabirov8872/golang-rest-api/internal/handler"
)

func Run(handler handler.IHandler, port string, secretKey string) {
	router := mux.NewRouter()

	router.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", handler.GetUserById).Methods("GET")
	router.HandleFunc("/users", handler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)

	log.Fatal(http.ListenAndServe("localhost:"+port, cors(router)))
}
