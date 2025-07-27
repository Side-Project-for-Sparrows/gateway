package router

import (
	"github.com/Side-Project-for-Sparrows/gateway/internal/handler"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/security"
	"github.com/gorilla/mux"
)

func InitRoute() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/").Subrouter()
	//api.Use(middleware.RootMiddlewareHandler)
	api.Use(security.JWTAuthMiddleware().AdaptWithNext)
	//api.Use(middleware.JWTAuthMiddleware)
	//api.Use(middleware.RateLimiter)

	api.PathPrefix("/").HandlerFunc(handler.ProxyHandler)

	return r
}
