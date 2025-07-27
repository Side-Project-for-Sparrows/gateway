package router

import (
	"github.com/Side-Project-for-Sparrows/gateway/internal/handler"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware"
	"github.com/gorilla/mux"
)

func InitRoute() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.RootMiddlewareHandler)

	api.PathPrefix("/").HandlerFunc(handler.ProxyHandler)

	return r
}
