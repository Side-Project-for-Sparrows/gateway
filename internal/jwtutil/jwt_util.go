package jwtutil

import (
	"time"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"github.com/golang-jwt/jwt/v5"
	"encoding/pem"
	"github.com/Side-Project-for-Sparrows/gateway/config"
	"log"
	"sync/atomic"
)

var (
	secretKey       = []byte("secret-key")
	expirationTime  = 10 * time.Minute
	ErrTokenExpired = errors.New("access token expired")
	ErrTokenInvalid = errors.New("access token invalid or unexpected")
	ErrPublicKeyInvalid = errors.New("public key invalid or unexpected")
	atomicKey atomic.Value
	//PublicKey *rsa.PublicKey
)

func VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 서명 알고리즘 체크
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrTokenInvalid
		}
		keyAny := atomicKey.Load()
		if keyAny == nil{
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

func Initialize(env string){
	go func(){
		for {
			log.Print("[DEBUG] 공개키 polling 시작")
			keyBytes, err := fetchPublicKeyPEM(env)
			if err != nil {
				log.Printf("[WARN] 공개키 불러오기 실패: %v", err)
				continue
			}

			pubKey, err := parseRSAPublicKeyFromPEM(keyBytes)
			if err != nil {
				log.Printf("[WARN] 공개키 파싱 실패: %v", err)
				continue
			}

			atomicKey.Store(pubKey)
			log.Print("[INFO] 공개키 갱신 성공")

			time.Sleep(1 * time.Minute)
		}
	}()
}

func fetchPublicKeyPEM(env string)([]byte, error){
	log.Printf("[DEBUG] env: %s", env)
	log.Printf("[DEBUG] full jwt config: %+v", config.Conf.JwtConfig)
	log.Printf("[DEBUG] resolved URL: %s", config.Conf.JwtConfig[env].PublicKeyUrl)

	resp, err := http.Get(config.Conf.JwtConfig[env].PublicKeyUrl)

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