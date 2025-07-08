package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/Side-Project-for-Sparrows/gateway/internal/router"
	"github.com/Side-Project-for-Sparrows/gateway/internal/jwtutil"
)

func main() {

	config.InitConfig()
	r := router.InitRoute()

	pubKey,err := jwtutil.FetchAndParsePublicKey()
	if err != nil {
		log.Fatalf("공개키 불러오기 실패: %v", err)
	}
	jwtutil.PublicKey = pubKey

	fmt.Println("? Gateway server is running on port 7080...")
	log.Fatal(http.ListenAndServe(":7080", r))
}
