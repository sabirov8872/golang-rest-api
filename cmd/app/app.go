package app

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sabirov8872/golang-rest-api/internal/config"
	"github.com/sabirov8872/golang-rest-api/internal/database"
	"github.com/sabirov8872/golang-rest-api/internal/handler"
	"github.com/sabirov8872/golang-rest-api/internal/routes"
	"github.com/sabirov8872/golang-rest-api/internal/service"
)

func Run() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbPath := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer db.Close()

	repo := database.NewRepository(db)
	serv := service.NewService(repo)
	hand := handler.NewHandler(serv, cfg.SecretKey)
	routes.Run(hand, cfg.ServerPort, cfg.SecretKey)
}
