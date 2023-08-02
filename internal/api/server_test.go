package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {
	server := NewServer("8081", &Handler{})
	t.Cleanup(func() {
		err := server.Shutdown(context.Background())
		if err != nil {
			t.Error(err)
		}
	})

	go func() {
		err := server.Start()
		if !errors.Is(err, http.ErrServerClosed) {
			t.Error(err)
		}
	}()

	time.Sleep(time.Millisecond * 100)

	_, err := http.Get("http://localhost:8081")
	if err != nil {
		t.Error(err)
	}
}
