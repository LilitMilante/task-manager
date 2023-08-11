package api

import (
	"context"
	"net/http"

	"task-manager/gen/proto/task/v1/todolistv1connect"

	"github.com/gorilla/mux"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	srv *http.Server
}

func NewServer(port string, h *Handler) *Server {
	r := mux.NewRouter()

	r.Use(h.LoggingMiddleware)

	path, handler := todolistv1connect.NewTaskServiceHandler(h)
	r.Handle(path, handler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h2c.NewHandler(r, &http2.Server{}),
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
