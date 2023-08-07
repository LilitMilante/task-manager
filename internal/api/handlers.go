package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"task-manager/internal/api/entity"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Service interface {
	AddTask(ctx context.Context, task entity.Task) (entity.Task, error)
	TaskByID(ctx context.Context, id uuid.UUID) (entity.Task, error)
	Tasks(ctx context.Context) ([]entity.Task, error)
}

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task entity.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	task, err = h.s.AddTask(r.Context(), task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
}

func (h *Handler) TaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	task, err := h.s.TaskByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.s.Tasks(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
}
