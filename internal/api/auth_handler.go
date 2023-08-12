package api

import (
	"context"

	v1 "task-manager/gen/proto/auth/v1"
	"task-manager/internal/entity"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

type AuthService interface {
	SignUp(ctx context.Context, user entity.User) (entity.User, error)
}

type AuthHandler struct {
	l *zap.SugaredLogger
	s AuthService
}

func NewAuthHandler(l *zap.SugaredLogger, s AuthService) *AuthHandler {
	return &AuthHandler{
		l: l,
		s: s,
	}
}

func (a *AuthHandler) SigneUp(ctx context.Context, c *connect.Request[v1.SigneUpRequest]) (*connect.Response[v1.SigneUpResponse], error) {
	user := UserFromAPI(c.Msg)

	resp, err := a.s.SignUp(ctx, user)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.SigneUpResponse{User: UserToAPI(resp)}), nil
}
