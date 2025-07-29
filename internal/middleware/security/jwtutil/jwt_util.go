package jwtutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	jwtConfig "github.com/Side-Project-for-Sparrows/gateway/config/jwt"
	"github.com/Side-Project-for-Sparrows/gateway/lifecycle"
	"github.com/golang-jwt/jwt/v5"
)

type jwtInitializer struct{}

func init() {
	lifecycle.Register(&jwtInitializer{})
}

func (j *jwtInitializer) Construct() error {
	log.Println("[Construct] JWTUtil Initialize 호출")
	Initialize()
	return nil
}

func Initialize() {
	fetchAndParse()

	go func() {
		for {
			fetchAndParse()
			time.Sleep(1 * time.Minute)
		}
	}()
}

var (
	secretKey           = []byte("secret-key")
	expirationTime      = 10 * time.Minute
	ErrTokenExpired     = errors.New("access token expired")
	ErrTokenInvalid     = errors.New("access token invalid or unexpected")
	ErrPublicKeyInvalid = errors.New("public key invalid or unexpected")
	atomicKey           atomic.Value
	//PublicKey *rsa.PublicKey
)

func VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 서명 알고리즘 체크
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrTokenInvalid
		}
		keyAny := atomicKey.Load()
		if keyAny == nil {
			return 0, ErrPublicKeyInvalid
		}

		publicKey := keyAny.(*rsa.PublicKey)

		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrTokenExpired
		}
		return 0, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return 0, ErrTokenInvalid
		}
		var userID int64
		_, err := fmt.Sscan(sub, &userID)
		if err != nil {
			return 0, ErrTokenInvalid
		}
		return userID, nil
	}

	return 0, ErrTokenInvalid
}

func fetchAndParse() {
	log.Print("[DEBUG] 공개키 polling 시작")
	keyBytes, err := fetchPublicKeyPEM()
	if err != nil {
		log.Printf("[WARN] 공개키 불러오기 실패: %v", err)
	}

	pubKey, err := parseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		log.Printf("[WARN] 공개키 파싱 실패: %v", err)
	} else {
		atomicKey.Store(pubKey)
		log.Print("[INFO] 공개키 갱신 성공")
	}
}

func fetchPublicKeyPEM() ([]byte, error) {
	var config = jwtConfig.Config
	log.Printf("[DEBUG] full jwt config: %+v", config)
	log.Printf("[DEBUG] resolved URL: %s", config.PublicKeyUrl)

	resp, err := http.Get(config.PublicKeyUrl)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %w", err)
	}
	defer resp.Body.Close()
	keyBytes, err := io.ReadAll(resp.Body)
	return keyBytes, err
}

func parseRSAPublicKeyFromPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return rsaPub, nil
}

// jwt 인증 불필요 여부 확인
func IsExcluded(path string) bool {
	for _, excluded := range jwtConfig.Config.ExcludedPaths {
		if strings.HasPrefix(path, excluded) {
			return true
		}
	}
	return false
}
