package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Side-Project-for-Sparrows/gateway/config"
	"github.com/Side-Project-for-Sparrows/gateway/internal/jwtutil"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL, _ := resolveTargetURL(r.URL.Path)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read backend response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	if r.URL.Path == "/user/auth/login" || r.URL.Path == "/user/auth/join" || strings.HasPrefix(r.URL.Path, "/index") {
		var respMap map[string]any
		if err := json.Unmarshal(bodyBytes, &respMap); err != nil {
			http.Error(w, "failed to parse response", http.StatusInternalServerError)
			return
		}

		userIDFloat, ok := respMap["id"].(float64)
		if !ok {
			http.Error(w, "userId missing or invalid", http.StatusInternalServerError)
			return
		}
		accessToken, err := jwtutil.GenerateToken(int64(userIDFloat), 10*time.Minute)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		respMap["accessToken"] = accessToken
		respMap["refreshToken"] = "dummy-refresh-token"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respMap)
		return
	}

	w.Write(bodyBytes)
}

func resolveTargetURL(path string) (string, bool) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev" // default fallback
	}

	routes := config.Conf.Routes[env]

	var base string
	switch {
	case strings.HasPrefix(path, "/user"):
		base = routes.User
	case strings.HasPrefix(path, "/board"):
		base = routes.Board
	case strings.HasPrefix(path, "/post"):
		base = routes.Board
	case strings.HasPrefix(path, "/school"):
		base = routes.School
	case strings.HasPrefix(path, "/index"):
		base = routes.Search
	default:
		return "", false
	}

	return base + path, true
}
