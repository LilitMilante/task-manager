package main

import (
	"fmt"

	"task-manager/internal/api"
)

const (
	port = "8080"
)

func main() {
	h := api.NewHandler()

	server := api.NewServer(port, h)
	err := server.Start()
	if err != nil {
		fmt.Errorf("start server error")
		return
	}
}
