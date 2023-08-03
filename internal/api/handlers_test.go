package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"task-manager/internal/api/entity"
	"task-manager/internal/service"
)

func TestHandler_AddTask(t *testing.T) {
	task := entity.Task{
		Name:        "test",
		Description: "test test test",
	}

	req, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(req))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	w := httptest.NewRecorder()

	s := service.NewService()
	handler := NewHandler(s)

	handler.AddTask(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected: %d\ngot: %d\nbody: %s", http.StatusOK, w.Code, w.Body)
	}

	var got entity.Task
	err = json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID == 0 {
		t.Error("ID is zero")
	}

	task.ID = got.ID

	if got.CreatedAt.IsZero() {
		t.Error("created at is zero")
	}

	task.CreatedAt = got.CreatedAt

	equal := reflect.DeepEqual(task, got)
	if !equal {
		t.Fatalf("expected: %+v\ngot: %+v", task, got)
	}
}
