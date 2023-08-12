package api

import (
	"context"
	"encoding/json"
	"net/http"

	v1 "task-manager/gen/proto/task/v1"
	"task-manager/internal/entity"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TaskService interface {
	AddTask(ctx context.Context, task entity.Task) (entity.Task, error)
	TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error)
	Tasks(ctx context.Context) ([]entity.Task, error)
	UpdateTask(ctx context.Context, updateTask entity.TaskUpdated) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
}

type TaskHandler struct {
	l *zap.SugaredLogger
	s TaskService
}

func NewTaskHandler(l *zap.SugaredLogger, s TaskService) *TaskHandler {
	return &TaskHandler{
		l: l,
		s: s,
	}
}

func (h *TaskHandler) AddTask(ctx context.Context, c *connect.Request[v1.AddTaskRequest]) (*connect.Response[v1.AddTaskResponse], error) {
	task := TaskFromAPI(c.Msg)

	task, err := h.s.AddTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.AddTaskResponse{Task: TaskToAPI(task)}), nil
}

func (h *TaskHandler) TaskByID(ctx context.Context, c *connect.Request[v1.TaskByIDRequest]) (*connect.Response[v1.TaskByIDResponse], error) {
	id, err := uuid.Parse(c.Msg.Id)
	if err != nil {
		return nil, err
	}

	task, err := h.s.TaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.TaskByIDResponse{Task: TaskToAPI(task)}), nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, c *connect.Request[v1.UpdateTaskRequest]) (*connect.Response[v1.UpdateTaskResponse], error) {
	updateTask, err := UpdateTaskFromAPI(c.Msg)
	if err != nil {
		return nil, err
	}

	err = h.s.UpdateTask(ctx, updateTask)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.UpdateTaskResponse{}), nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, c *connect.Request[v1.DeleteTaskRequest]) (*connect.Response[v1.DeleteTaskResponse], error) {
	id, err := uuid.Parse(c.Msg.Id)
	if err != nil {
		return nil, err
	}

	err = h.s.DeleteTask(ctx, id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.DeleteTaskResponse{}), nil
}

func (h *TaskHandler) Tasks(ctx context.Context, _ *connect.Request[v1.TasksRequest]) (*connect.Response[v1.TasksResponse], error) {
	tasks, err := h.s.Tasks(ctx)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.TasksResponse{Tasks: TasksToAPI(tasks)}), nil
}

//Helpers

type JSONErr struct {
	Error string `json:"error"`
}

func (h *TaskHandler) SendJsonError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)

	resp := JSONErr{
		err.Error(),
	}

	h.SendJson(w, resp)
}

func (h *TaskHandler) SendJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.l.Error(err)
		return
	}

	h.l.Infow("send response", "response", string(resp))

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.l.Error(err)
		return
	}
}
