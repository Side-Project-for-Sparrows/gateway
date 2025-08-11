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

	//미들웨어 실행시간 측정용 더미 핸들러
	api.PathPrefix("/user/dummy").HandlerFunc(handler.DummyHandler)
	api.PathPrefix("/").HandlerFunc(handler.LoggingWrapper(handler.ProxyHandler))

	return r
}
