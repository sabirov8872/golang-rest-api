package app

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/sabirov8872/golang-rest-api/internal/config"
	"github.com/sabirov8872/golang-rest-api/internal/database"
	"github.com/sabirov8872/golang-rest-api/internal/handler"
	"github.com/sabirov8872/golang-rest-api/internal/routes"
	"github.com/sabirov8872/golang-rest-api/internal/service"
	"log"
)

func Run() {
	cnf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cnf.DatabasePath)
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer db.Close()

	repo := database.NewRepository(db)
	serv := service.NewService(repo)
	hand := handler.NewHandler(serv)
	routes.Run(hand, cnf.ServerAddress)
}
