package handler

import (
	"log"
	"net/http"
	"time"
)

func LoggingWrapper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC] %v", rec)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next(w, r)

		duration := time.Since(start)
		log.Printf("[LOG] %s %s took %s", r.Method, r.URL.Path, duration)
		w.Header().Set("X-Latency", duration.String())
	}
}
