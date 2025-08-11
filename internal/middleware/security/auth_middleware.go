package security

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/security/jwtutil"
)

func JWTAuthMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		// 요청 경로에서 인증 제외 대상이면 패스
		if jwtutil.IsExcluded(input.Path()) {
			return nil, nil
		}

		authHeader := input.Headers().Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Print("엑세스토큰 없음")
			return nil, fmt.Errorf("unauthorized: missing bearer token")
		}

		// 토큰 파싱
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwtutil.VerifyToken(tokenString)
		if err != nil {
			log.Print("유효하지 않은 토큰 %s", tokenString)
			return nil, fmt.Errorf("unauthorized: invalid token")
		}

		return &middlewaretype.HeaderPatch{
			RequestAdd: http.Header{
				"X-Requester-Id": []string{fmt.Sprintf("%d", userID)},
			},
		}, nil
	}
}
