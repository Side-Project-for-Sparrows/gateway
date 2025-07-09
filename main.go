package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/Side-Project-for-Sparrows/gateway/internal/router"
	"github.com/Side-Project-for-Sparrows/gateway/internal/jwtutil"
)

func main() {
	initialize()
	route()
}

func initialize() {
	env := getEnv()
	config.InitConfig()
	jwtutil.Initialize(env)
}

func route(){
	r := router.InitRoute()
	fmt.Println("? Gateway server is running on port 7080...")
	log.Fatal(http.ListenAndServe(":7080", r))
}

func getEnv()(string){
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatalf("ENV 환경변수 없음. 'dev', 'prod' 중 하나를 설정하세요.")
	}

	return env
}