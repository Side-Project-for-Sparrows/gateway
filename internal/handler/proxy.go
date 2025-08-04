package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Side-Project-for-Sparrows/gateway/config/route"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL, ok := resolveTargetURL(r.URL.Path)

	if !ok {
		http.Error(w, "route to server failed", http.StatusInternalServerError)
		return
	}

	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}
	fmt.Println("call target : " + targetURL)
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header.Clone()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("backend unreachable: %v", err)
		fmt.Println(errMsg)
		http.Error(w, errMsg, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func resolveTargetURL(path string) (string, bool) {
	// /user/login → user
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(segments) == 0 {
		return "", false
	}
	key := segments[0]

	// key : user -> value : localhost:8080
	baseURL, ok := route.RouteMap[key]
	if !ok {
		return "", false
	}

	return baseURL + path, true
}
