package main

import (
	"fmt"
	"log"

	"task-manager/internal/api"
	"task-manager/internal/app"
	"task-manager/internal/repository"
	"task-manager/internal/service"
)

const (
	port     = "8080"
	dbHost   = "localhost"
	dbPort   = 5432
	user     = "postgres"
	password = "your-password"
	dbName   = "task-manager"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName)

	db, err := app.ConnectToPostgres(dsn)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	s := service.NewService(repo)
	h := api.NewHandler(s)

	server := api.NewServer(port, h)
	err = server.Start()
	if err != nil {
		log.Fatal("start server error")
	}
}
