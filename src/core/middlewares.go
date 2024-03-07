package core

import (
	"net/http"
	"template-base-go/src/utils"
)

func LogRequest(log utils.ILogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Request received. Method: " + r.Method + " Url: " + r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
