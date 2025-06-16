package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/internal/jwtutil"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := "http://localhost:8080" + r.URL.Path
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "backend unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 응답 헤더 복사
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	// 응답 바디 읽기
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read backend response", http.StatusInternalServerError)
		return
	}

	// 기본 응답 처리
	w.WriteHeader(resp.StatusCode)

	if r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/join" {
		var springResp map[string]any
		if err := json.Unmarshal(bodyBytes, &springResp); err != nil {
			http.Error(w, "failed to parse response", http.StatusInternalServerError)
			return
		}

		// userId -> float64 → int64
		userIDFloat, ok := springResp["id"].(float64)
		if !ok {
			http.Error(w, "userId missing or invalid", http.StatusInternalServerError)
			return
		}
		accessToken, err := jwtutil.GenerateToken(int64(userIDFloat))
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		// 토큰 추가
		springResp["accessToken"] = accessToken
		springResp["refreshToken"] = "dummy-refresh-token"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(springResp)
		return
	}

	// 로그인/회원가입이 아니면 원본 그대로 전달
	w.Write(bodyBytes)
}
