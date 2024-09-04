package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sabirov8872/golang-rest-api/internal/handler"
)

func Run(handler handler.IHandler, port string) {
	router := mux.NewRouter()
	router.HandleFunc("/token", handler.CheckUser).Methods("GET")
	router.HandleFunc("/user", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/user/{id}", handler.GetUserById).Methods("GET")
	router.HandleFunc("/user", handler.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", handler.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe("localhost:"+port, router))
}
