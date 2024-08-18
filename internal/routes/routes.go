package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sabirov8872/golang-rest-api/internal/handler"
	"log"
	"net/http"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func Run(handler handler.IHandler, addr string) {
	router := mux.NewRouter()
	router.HandleFunc("/", helloWorld).Methods("GET")
	router.HandleFunc("/user", handler.GetAllUsers).Methods("GET")
	router.HandleFunc("/user/{id}", handler.GetUserById).Methods("GET")
	router.HandleFunc("/user/{first_name}/{username}/{phone}", handler.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}/{first_name}/{username}/{phone}", handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", handler.DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(addr, router))
}
