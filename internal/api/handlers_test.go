package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-manager/internal/api/entity"
	"task-manager/internal/app"
	"task-manager/internal/repository"
	"task-manager/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	handler := newHandler(t)

	// Add task

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
	require.Equal(t, http.StatusOK, w.Code)

	var got entity.Task
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)

	task.ID = got.ID
	task.CreatedAt = got.CreatedAt

	require.NotZero(t, got.ID)
	require.NotZero(t, got.CreatedAt)
	require.Equal(t, task, got)

	// Get created task by ID

	id := task.ID.String()

	r, err = http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	require.NoError(t, err)
	r = mux.SetURLVars(r, map[string]string{"id": id})

	w = httptest.NewRecorder()
	handler.TaskByID(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	got = entity.Task{}
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, task, got)
}
func TestHandler_Tasks(t *testing.T) {
	expected := []entity.Task{
		{
			Name:        "test1",
			Description: "test t1",
		},
		{
			Name:        "test2",
			Description: "test t2",
		},
	}

	handler := newHandler(t, func(db *sql.DB) {
		_, err := db.Exec("DELETE FROM tasks")
		require.NoError(t, err)
	})

	r, err := http.NewRequest(http.MethodGet, "/alltasks", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	for i, v := range expected {
		task, err := handler.s.AddTask(r.Context(), v)
		require.NoError(t, err)

		expected[i].ID = task.ID
		expected[i].CreatedAt = task.CreatedAt
	}

	handler.Tasks(w, r)
	require.Equal(t, http.StatusOK, w.Code, w.Body)

	var got []entity.Task
	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, expected, got)
}

func TestHandler_TaskByID(t *testing.T) {
	handler := newHandler(t)

	id := uuid.New().String()
	r, err := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	require.NoError(t, err)
	r = mux.SetURLVars(r, map[string]string{"id": id})

	w := httptest.NewRecorder()
	handler.TaskByID(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)

	exp := "sql: no rows in result set"
	var got JSONErr

	err = json.NewDecoder(w.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, exp, got.Error)
}

func newHandler(t *testing.T, fns ...func(db *sql.DB)) *Handler {
	t.Helper()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName)

	db, err := app.ConnectToPostgres(dsn)
	require.NoError(t, err)

	for _, fn := range fns {
		fn(db)
	}

	t.Cleanup(func() { db.Close() })

	repo := repository.NewRepository(db)
	s := service.NewService(repo)
	return NewHandler(s)
}
