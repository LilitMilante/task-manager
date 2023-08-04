package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-manager/internal/api/entity"
	"task-manager/internal/app"
	"task-manager/internal/repository"
	"task-manager/internal/service"

	"github.com/stretchr/testify/require"
)

const (
	dbHost   = "localhost"
	dbPort   = 5432
	user     = "postgres"
	password = "your-password"
	dbName   = "task-manager"
)

func TestHandler_AddTask(t *testing.T) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName)

	db, err := app.ConnectToPostgres(dsn)
	require.NoError(t, err)

	t.Cleanup(func() { db.Close() })

	repo := repository.NewRepository(db)
	s := service.NewService(repo)
	handler := NewHandler(s)

	task := entity.Task{
		Name:        "test",
		Description: "test test test",
	}

	req, err := json.Marshal(task)
	require.NoError(t, err)

	r, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(req))
	require.NoError(t, err)
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	w := httptest.NewRecorder()
	handler.AddTask(w, r)
	require.Equal(t, w.Code, http.StatusOK)

	var got entity.Task
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)

	task.ID = got.ID
	task.CreatedAt = got.CreatedAt

	require.NotZero(t, got.ID)
	require.NotZero(t, got.CreatedAt)
	require.Equal(t, task, got)
}
