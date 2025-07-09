package middleware

import (
	"net/http"
	"strings"

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
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := jwtutil.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
