package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"task-manager/internal/api/entity"
)

type Service interface {
	AddTask(ctx context.Context, task entity.Task) (entity.Task, error)
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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
