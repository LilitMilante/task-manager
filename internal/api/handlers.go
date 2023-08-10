package api

import (
	"context"
	"encoding/json"
	"net/http"

	v1 "task-manager/gen/proto/task/v1"
	"task-manager/internal/api/entity"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Service interface {
	AddTask(ctx context.Context, task entity.Task) (entity.Task, error)
	TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error)
	Tasks(ctx context.Context) ([]entity.Task, error)
	UpdateTask(ctx context.Context, id uuid.UUID, updateTask entity.TaskUpdated) error
}

type Handler struct {
	l *zap.SugaredLogger
	s Service
}

func NewHandler(l *zap.SugaredLogger, s Service) *Handler {
	return &Handler{
		l: l,
		s: s,
	}
}

func (h *Handler) AddTask(ctx context.Context, c *connect.Request[v1.AddTaskRequest]) (*connect.Response[v1.AddTaskResponse], error) {
	task := TaskFromAPI(c.Msg)

	task, err := h.s.AddTask(ctx, task)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(TaskToAPI(task))

	return resp, nil
}

//func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
//	var task entity.Task
//
//	err := json.NewDecoder(r.Body).Decode(&task)
//	if err != nil {
//		h.SendJsonError(w, http.StatusBadRequest, err)
//		return
//	}
//
//	task, err = h.s.AddTask(r.Context(), task)
//	if err != nil {
//		h.SendJsonError(w, http.StatusInternalServerError, err)
//		return
//	}
//
//	h.SendJson(w, task)
//}

func (h *Handler) TaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.SendJsonError(w, http.StatusBadRequest, err)
		return
	}

	task, err := h.s.TaskByID(r.Context(), id)
	if err != nil {
		h.SendJsonError(w, http.StatusNotFound, err)
		return
	}

	h.SendJson(w, task)
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.s.Tasks(r.Context())
	if err != nil {
		h.SendJsonError(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJson(w, tasks)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.SendJsonError(w, http.StatusBadRequest, err)
		return
	}

	var updateTask entity.TaskUpdated

	err = json.NewDecoder(r.Body).Decode(&updateTask)
	if err != nil {
		h.SendJsonError(w, http.StatusBadRequest, err)
		return
	}

	err = h.s.UpdateTask(r.Context(), id, updateTask)
	if err != nil {
		h.SendJsonError(w, http.StatusInternalServerError, err)
		return
	}
}

//Helpers

type JSONErr struct {
	Error string `json:"error"`
}

func (h *Handler) SendJsonError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)

	resp := JSONErr{
		err.Error(),
	}

	h.SendJson(w, resp)
}

func (h *Handler) SendJson(w http.ResponseWriter, data any) {
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

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.l.Infof("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
