package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func LoggingMiddleware(l *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Infof("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}
