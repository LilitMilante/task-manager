package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_AddTask(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := NewHandler()

	handler.AddTask(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected: %d\ngot: %d", http.StatusOK, w.Code)
	}

	exp := "Hello"
	got := w.Body.String()
	if exp != got {
		t.Fatalf("expected: %q\ngot: %q", exp, got)
	}
}
