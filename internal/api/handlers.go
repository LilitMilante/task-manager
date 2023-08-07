package api

import (
	"context"
	"encoding/json"
	"net/http"

	"task-manager/internal/api/entity"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Service interface {
	AddTask(ctx context.Context, task entity.Task) (entity.Task, error)
	TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error)
	Tasks(ctx context.Context) ([]entity.Task, error)
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

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task entity.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		h.SendJsonError(w, http.StatusBadRequest, err)
		return
	}

	task, err = h.s.AddTask(r.Context(), task)
	if err != nil {
		h.SendJsonError(w, http.StatusInternalServerError, err)
		return
	}

	h.SendJson(w, task)
}

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
