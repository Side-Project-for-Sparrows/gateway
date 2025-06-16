package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/internal/router"
)

func main() {
	r := router.NewRouter()

	fmt.Println("? Gateway server is running on port 7080...")
	log.Fatal(http.ListenAndServe(":7080", r))
}
