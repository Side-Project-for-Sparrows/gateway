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
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		// 요청 경로에서 인증 제외 대상이면 패스
		if excludedPaths[input.Path()] {
			return nil, nil
		}

		log.Printf("[JWT] input.Path=%v", input.Path())
		log.Printf("[JWT] input.HEADER=%q", input.Headers())
		authHeader := input.Headers().Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Print("엑세스토큰 없음")
			return nil, fmt.Errorf("unauthorized: missing bearer token")
		}

		// 토큰 파싱
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwtutil.VerifyToken(tokenString)
		if err != nil {
			log.Print("유효하지 않은 토큰")
			return nil, fmt.Errorf("unauthorized: invalid token")
		}

		log.Printf("인증된 유저 ID: %d", userID)

		return &middlewaretype.HeaderPatch{
			RequestAdd: http.Header{
				"X-Requester-Id": []string{fmt.Sprintf("%d", userID)},
			},
		}, nil
	}
}
