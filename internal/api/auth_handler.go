package api

import (
	"context"
	"net/http"
	"time"

	v1 "task-manager/gen/proto/auth/v1"
	"task-manager/internal/entity"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService interface {
	SignUp(ctx context.Context, user entity.User) (entity.User, error)
	SignIn(ctx context.Context, email, password string) (uuid.UUID, error)
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

func (a *AuthHandler) SigneIn(ctx context.Context, c *connect.Request[v1.SigneInRequest]) (*connect.Response[v1.SigneInResponse], error) {
	id, err := a.s.SignIn(ctx, c.Msg.Email, c.Msg.Password)
	if err != nil {
		return nil, err
	}

	const secondsInHour = 3600

	cookie := http.Cookie{
		Name:     "sid",
		Value:    id.String(),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		MaxAge:   secondsInHour,
		HttpOnly: true,
	}

	resp := connect.NewResponse(&v1.SigneInResponse{})
	resp.Header().Set("Set-Cookie", cookie.String())

	return resp, nil
}
