package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/Side-Project-for-Sparrows/gateway/internal/router"
)

func main() {

	config.InitConfig()

	r := router.InitRoute()

	fmt.Println("? Gateway server is running on port 7080...")
	log.Fatal(http.ListenAndServe(":7080", r))
}
