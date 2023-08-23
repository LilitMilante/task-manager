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
	"github.com/google/uuid"
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

func taskClient(t *testing.T) todolistv1connect.TaskServiceClient {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName)

	db, err := app.ConnectToPostgres(dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	repo := repository.NewRepository(l, db)
	auth := service.NewAuthService(repo) // заменить на mock
	s := service.NewTaskService(repo, auth)
	h := NewTaskHandler(l, s)

	_, handler := todolistv1connect.NewTaskServiceHandler(h)

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return todolistv1connect.NewTaskServiceClient(http.DefaultClient, server.URL)
}

func TestHandler_AddTask(t *testing.T) {
	client := taskClient(t)

	// Add task
	addTaskReq := &todolistv1.AddTaskRequest{
		Name:        "Test",
		Description: "Test test",
	}

	resp, err := client.AddTask(context.Background(), connect.NewRequest(addTaskReq))
	require.NoError(t, err)

	require.Equal(t, addTaskReq.Name, resp.Msg.Task.Name)
	require.Equal(t, addTaskReq.Description, resp.Msg.Task.Description)
	require.NotEmpty(t, resp.Msg.Task.Id)
	require.NotZero(t, resp.Msg.Task.CreatedAt)
	require.NotZero(t, resp.Msg.Task.UpdatedAt)

	// Get task by ID
	getTaskReq := &todolistv1.TaskByIDRequest{
		Id: resp.Msg.Task.Id,
	}

	taskByID, err := client.TaskByID(context.Background(), connect.NewRequest(getTaskReq))
	require.NoError(t, err)

	require.Equal(t, resp.Msg.Task.Id, taskByID.Msg.Task.Id)
}

func TestHandler_UpdateTask(t *testing.T) {
	client := taskClient(t)

	// Add task
	addTaskReq := &todolistv1.AddTaskRequest{
		Name:        "Test",
		Description: "Test test",
	}

	task, err := client.AddTask(context.Background(), connect.NewRequest(addTaskReq))
	require.NoError(t, err)

	// Update task
	updateTask := &todolistv1.UpdateTaskRequest{
		Id:          task.Msg.Task.Id,
		Name:        "TEST",
		Description: "TEST test TEST",
		IsCompleted: true,
	}

	_, err = client.UpdateTask(context.Background(), connect.NewRequest(updateTask))
	require.NoError(t, err)

	// Get update task
	getTaskReq := &todolistv1.TaskByIDRequest{
		Id: task.Msg.Task.Id,
	}

	taskByID, err := client.TaskByID(context.Background(), connect.NewRequest(getTaskReq))
	require.NoError(t, err)

	require.Equal(t, updateTask.Id, taskByID.Msg.Task.Id)
	require.Equal(t, updateTask.Name, taskByID.Msg.Task.Name)
	require.Equal(t, updateTask.Description, taskByID.Msg.Task.Description)
	require.Equal(t, updateTask.IsCompleted, taskByID.Msg.Task.IsCompleted)
	require.Greater(t, taskByID.Msg.Task.UpdatedAt.Nanos, task.Msg.Task.UpdatedAt.Nanos)
}

func TestHandler_UpdateTask_Error(t *testing.T) {
	client := taskClient(t)

	// Update task
	updateTask := &todolistv1.UpdateTaskRequest{
		Id:          uuid.NewString(),
		Name:        "TEST",
		Description: "TEST test TEST",
		IsCompleted: true,
	}

	_, err := client.UpdateTask(context.Background(), connect.NewRequest(updateTask))
	require.Error(t, err)
}

func TestHandler_DeleteTask(t *testing.T) {
	client := taskClient(t)

	// Add task
	addTaskReq := &todolistv1.AddTaskRequest{
		Name:        "Test",
		Description: "Test test",
	}

	task, err := client.AddTask(context.Background(), connect.NewRequest(addTaskReq))
	require.NoError(t, err)

	// Delete task
	deleteTask := &todolistv1.DeleteTaskRequest{
		Id: task.Msg.Task.Id,
	}
	_, err = client.DeleteTask(context.Background(), connect.NewRequest(deleteTask))
	require.NoError(t, err)

	// Get delete task
	getTaskReq := &todolistv1.TaskByIDRequest{
		Id: task.Msg.Task.Id,
	}

	_, err = client.TaskByID(context.Background(), connect.NewRequest(getTaskReq))
	require.Error(t, err)
}

func TestHandler_DeleteTask_Error(t *testing.T) {
	client := taskClient(t)

	// Delete task
	deleteTask := &todolistv1.DeleteTaskRequest{
		Id: uuid.NewString(),
	}
	_, err := client.DeleteTask(context.Background(), connect.NewRequest(deleteTask))
	require.Error(t, err)
}

func TestHandler_Tasks(t *testing.T) {
	client := taskClient(t)

	// Add tasks
	got := []*todolistv1.AddTaskRequest{
		{
			Name:        "Task1",
			Description: "test task1",
		},
		{
			Name:        "Task2",
			Description: "test task2",
		},
	}

	var exp []*todolistv1.Task

	for _, v := range got {
		task, err := client.AddTask(context.Background(), connect.NewRequest(v))
		require.NoError(t, err)
		exp = append(exp, task.Msg.Task)
	}

	// Get tasks
	tasks, err := client.Tasks(context.Background(), connect.NewRequest(&todolistv1.TasksRequest{}))
	require.NoError(t, err)

	for _, v := range exp {
		require.Contains(t, tasks.Msg.Tasks, v)
	}
}
