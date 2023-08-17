package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	userv1 "task-manager/gen/proto/auth/v1"
	"task-manager/gen/proto/auth/v1/userv1connect"
	"task-manager/internal/app"
	"task-manager/internal/repository"
	"task-manager/internal/service"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func authClient(t *testing.T) userv1connect.AuthServiceClient {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName)

	db, err := app.ConnectToPostgres(dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	repo := repository.NewRepository(l, db)
	s := service.NewAuthService(repo)
	h := NewAuthHandler(l, s)

	_, handler := userv1connect.NewAuthServiceHandler(h)

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return userv1connect.NewAuthServiceClient(http.DefaultClient, server.URL)
}

func TestAuthHandler_SignUp(t *testing.T) {
	client := authClient(t)

	req := &userv1.SignUpRequest{
		Name:     "Nasta",
		Email:    uuid.NewString() + "@gmail.com",
		Password: "123test",
	}

	resp, err := client.SignUp(context.Background(), connect.NewRequest(req))
	require.NoError(t, err)

	got := resp.Msg.User

	require.NotEmpty(t, got.Id)
	require.Equal(t, req.Name, got.Name)
	require.Equal(t, req.Email, got.Email)
	require.NotZero(t, got.CreatedAt)

	signInResp, err := client.SignIn(context.Background(), connect.NewRequest(&userv1.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	}))
	require.NoError(t, err)
	require.NotEmpty(t, signInResp.Header().Get("Set-Cookie"))
}
