package middleware

import (
	"net/http"
	"strings"
	"log"
	"fmt"
	"github.com/Side-Project-for-Sparrows/gateway/internal/jwtutil"
)

var excludedPaths = map[string]bool{
	"/user/auth/login": true,
	"/user/auth/join":  true,
	"/user/auth/refresh": true,
	"/index/school":    true,
	"/dummy":           true,
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if excludedPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Print("엑세스토큰 없음")
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := jwtutil.VerifyToken(tokenString)
		if err != nil {
			log.Print("유효하지 않은 토큰")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Print(userID)
		r.Header.Set("X-Requester-Id", fmt.Sprintf("%d", userID))

		next.ServeHTTP(w, r)
	})
}
