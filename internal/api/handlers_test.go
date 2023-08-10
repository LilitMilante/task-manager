package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	todolistv1 "task-manager/gen/proto/task/v1"
	"task-manager/gen/proto/task/v1/todolistv1connect"
	"task-manager/internal/app"
	"task-manager/internal/repository"
	"task-manager/internal/service"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	dbHost   = "localhost"
	dbPort   = 5432
	user     = "postgres"
	password = "your-password"
	dbName   = "task-manager"
)

var l = zap.NewNop().Sugar()

//	func TestHandler_AddTask(t *testing.T) {
//		handler := newHandler(t)
//
//		// Add task
//		exp := entity.Task{
//			Name:        "test",
//			Description: "test test test",
//		}
//
//		got := addTask(t, handler, exp)
//
//		exp.ID = got.ID
//		exp.CreatedAt = got.CreatedAt
//		exp.UpdatedAt = got.UpdatedAt
//
//		require.NotZero(t, got.ID)
//		require.NotZero(t, got.CreatedAt)
//		require.NotZero(t, got.UpdatedAt)
//		require.Equal(t, exp, got)
//
//		// Get created task by ID
//
//		id := exp.ID.String()
//
//		r, err := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
//		require.NoError(t, err)
//		r = mux.SetURLVars(r, map[string]string{"id": id})
//
//		w := httptest.NewRecorder()
//		handler.TaskByID(w, r)
//		require.Equal(t, http.StatusOK, w.Code)
//
//		got = entity.Task{}
//		err = json.NewDecoder(w.Body).Decode(&got)
//		require.NoError(t, err)
//
//		require.Equal(t, exp, got)
//	}
//
//	func TestHandler_Tasks(t *testing.T) {
//		expected := []entity.Task{
//			{
//				Name:        "test1",
//				Description: "test t1",
//			},
//			{
//				Name:        "test2",
//				Description: "test t2",
//			},
//		}
//
//		handler := newHandler(t, func(db *sql.DB) {
//			//goland:noinspection SqlWithoutWhere
//			_, err := db.Exec("DELETE FROM tasks")
//			require.NoError(t, err)
//		})
//
//		r, err := http.NewRequest(http.MethodGet, "/tasks", nil)
//		require.NoError(t, err)
//
//		w := httptest.NewRecorder()
//
//		for i, v := range expected {
//			task, err := handler.s.AddTask(r.Context(), v)
//			require.NoError(t, err)
//
//			expected[i].ID = task.ID
//			expected[i].CreatedAt = task.CreatedAt
//			expected[i].UpdatedAt = task.CreatedAt
//		}
//
//		handler.Tasks(w, r)
//		require.Equal(t, http.StatusOK, w.Code, w.Body)
//
//		var got []entity.Task
//		err = json.NewDecoder(w.Body).Decode(&got)
//		require.NoError(t, err)
//
//		require.Equal(t, expected, got)
//	}
//
//	func TestHandler_TaskByID_Error(t *testing.T) {
//		handler := newHandler(t)
//
//		id := uuid.New().String()
//		r, err := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
//		require.NoError(t, err)
//		r = mux.SetURLVars(r, map[string]string{"id": id})
//
//		w := httptest.NewRecorder()
//		handler.TaskByID(w, r)
//		require.Equal(t, http.StatusNotFound, w.Code)
//
//		exp := "sql: no rows in result set"
//		var got JSONErr
//
//		err = json.NewDecoder(w.Body).Decode(&got)
//		require.NoError(t, err)
//
//		require.Equal(t, exp, got.Error)
//	}
//
//	func TestHandler_UpdateTask(t *testing.T) {
//		handler := newHandler(t)
//
//		// Add new task
//		exp := entity.Task{
//			Name:        "test",
//			Description: "test test test",
//		}
//
//		got := addTask(t, handler, exp)
//
//		// Create task for update
//		updateTask := entity.TaskUpdated{
//			Name:        "TesT",
//			Description: "Test Test",
//			IsCompleted: true,
//		}
//
//		// Update task
//		req, err := json.Marshal(updateTask)
//		require.NoError(t, err)
//
//		id := got.ID.String()
//		r, err := http.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewBuffer(req))
//		require.NoError(t, err)
//		r = mux.SetURLVars(r, map[string]string{"id": id})
//		r.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		w := httptest.NewRecorder()
//		handler.UpdateTask(w, r)
//		require.Equal(t, http.StatusOK, w.Code)
//
//		task, err := handler.s.TaskByID(r.Context(), got.ID)
//
//		require.NotZero(t, task.UpdatedAt)
//		require.True(t, task.UpdatedAt.After(got.UpdatedAt))
//
//		require.Equal(t, updateTask.Name, task.Name)
//		require.Equal(t, updateTask.Description, task.Description)
//		require.Equal(t, updateTask.IsCompleted, task.IsCompleted)
//	}
//
//	func addTask(t *testing.T, handler *Handler, exp entity.Task) (got entity.Task) {
//		t.Helper()
//
//		req, err := json.Marshal(exp)
//		require.NoError(t, err)
//
//		r, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(req))
//		require.NoError(t, err)
//		r.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		w := httptest.NewRecorder()
//		handler.AddTask(w, r)
//		require.Equal(t, http.StatusOK, w.Code)
//
//		err = json.NewDecoder(w.Body).Decode(&got)
//		require.NoError(t, err)
//
//		return got
//	}
func newClient(t *testing.T) todolistv1connect.TaskServiceClient {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName)

	db, err := app.ConnectToPostgres(dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	repo := repository.NewRepository(l, db)
	s := service.NewService(repo)
	h := NewHandler(l, s)

	_, handler := todolistv1connect.NewTaskServiceHandler(h)

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return todolistv1connect.NewTaskServiceClient(http.DefaultClient, server.URL)
}

func TestHandler_AddTask(t *testing.T) {
	client := newClient(t)

	req := &todolistv1.AddTaskRequest{
		Name:        "Test",
		Description: "Test test",
	}

	resp, err := client.AddTask(context.Background(), connect.NewRequest(req))
	require.NoError(t, err)

	require.Equal(t, req.Name, resp.Msg.Name)
	require.Equal(t, req.Description, resp.Msg.Description)
	require.NotEmpty(t, resp.Msg.Id)
	require.NotZero(t, resp.Msg.CreatedAt)
	require.NotZero(t, resp.Msg.UpdatedAt)
}
