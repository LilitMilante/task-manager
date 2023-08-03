package main

import (
	"log"

	"task-manager/internal/api"
	"task-manager/internal/service"
)

const (
	port = "8080"
)

func main() {
	s := service.NewService()
	h := api.NewHandler(s)

	server := api.NewServer(port, h)
	err := server.Start()
	if err != nil {
		log.Fatal("start server error")
	}
}
