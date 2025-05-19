package router

import (
	"github.com/Side-Project-for-Sparrows/gateway/internal/handler"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	return r
}
