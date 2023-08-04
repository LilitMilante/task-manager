package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"task-manager/internal/api/entity"
	"task-manager/internal/app"
	"task-manager/internal/repository"
	"task-manager/internal/service"

	"github.com/google/uuid"
)

const (
	dbHost   = "localhost"
	dbPort   = 5432
	user     = "postgres"
	password = "your-password"
	dbName   = "task-manager"
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

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, user, password, dbName)

	db, err := app.ConnectToPostgres(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	repo := repository.NewRepository(db)
	s := service.NewService(repo)
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

	if got.ID == uuid.Nil {
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
