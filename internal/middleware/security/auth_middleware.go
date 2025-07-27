package security

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Side-Project-for-Sparrows/gateway/internal/jwtutil"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
)

var excludedPaths = map[string]bool{
	"/user/auth/login":   true,
	"/user/auth/join":    true,
	"/user/auth/refresh": true,
	"/index/school":      true,
	"/dummy":             true,
}

func JWTAuthMiddleware() middlewaretype.Middleware {
	return func(w http.ResponseWriter, r *http.Request) error {
		if excludedPaths[r.URL.Path] {
			return nil
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Print("엑세스토큰 없음")
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return fmt.Errorf("unauthorized: missing bearer token")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := jwtutil.VerifyToken(tokenString)
		if err != nil {
			log.Print("유효하지 않은 토큰")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return fmt.Errorf("unauthorized: invalid token")
		}

		log.Print(userID)
		r.Header.Set("X-Requester-Id", fmt.Sprintf("%d", userID))

		return nil
	}
}
