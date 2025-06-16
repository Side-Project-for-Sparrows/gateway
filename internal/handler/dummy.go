package handler

import (
	"fmt"
	"net/http"
)

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
