package api

import (
	"context"
	"net/http"

	"task-manager/gen/proto/auth/v1/userv1connect"
	"task-manager/gen/proto/task/v1/todolistv1connect"

	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	srv *http.Server
}

func NewServer(l *zap.SugaredLogger, port string, taskHandler *TaskHandler, authHandler *AuthHandler) *Server {
	r := http.NewServeMux()

	logMW := LoggingMiddleware(l)

	authMW := AuthMiddleware(authHandler.s)
	path, handler := todolistv1connect.NewTaskServiceHandler(taskHandler)

	r.Handle(path, authMW(handler))
	r.Handle(userv1connect.NewAuthServiceHandler(authHandler))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h2c.NewHandler(logMW(r), &http2.Server{}),
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
