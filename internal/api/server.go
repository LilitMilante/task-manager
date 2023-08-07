package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	srv *http.Server
}

func NewServer(port string, h *Handler) *Server {
	r := mux.NewRouter()

	r.Use(h.LoggingMiddleware)

	r.HandleFunc("/tasks", h.AddTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{id}", h.TaskByID).Methods(http.MethodGet)
	r.HandleFunc("/tasks", h.Tasks).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
