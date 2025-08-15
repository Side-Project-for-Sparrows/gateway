package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Side-Project-for-Sparrows/gateway/internal/handler"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware"
	"github.com/gorilla/mux"
)

func InitRoute() *mux.Router {

	cert, key := getTlsFilePath()

	// (선택) 실행 시 실제 경로 검증
	if _, err := os.Stat(cert); err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(key); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.RootMiddlewareHandler)

	//미들웨어 실행시간 측정용 더미 핸들러
	api.PathPrefix("/user/dummy").HandlerFunc(handler.DummyHandler)
	api.PathPrefix("/").HandlerFunc(handler.LoggingWrapper(handler.ProxyHandler))

	fmt.Println("? Gateway server is running on port 443...")
	//log.Fatal(http.ListenAndServe(":7080", r))
	log.Fatal(http.ListenAndServeTLS(
		":8443",
		cert,
		key,
		r,
	))
	return r
}

func getTlsFilePath() (string, string) {
	dir := os.Getenv("CERT_DIR")
	if dir == "" {
		dir = "/etc/tls"
	}
	cert := filepath.Join(dir, "tls.crt")
	key := filepath.Join(dir, "tls.key")
	if _, err := os.Stat(cert); err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(key); err != nil {
		log.Fatal(err)
	}
	return cert, key
}
