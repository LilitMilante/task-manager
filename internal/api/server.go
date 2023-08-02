package api

import (
	"context"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(port string, h *Handler) *Server {
	r := http.NewServeMux()

	r.HandleFunc("/tasks", h.AddTask)

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
