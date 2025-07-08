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
)

var (
	secretKey       = []byte("secret-key")
	expirationTime  = 10 * time.Minute
	ErrTokenExpired = errors.New("access token expired")
	ErrTokenInvalid = errors.New("access token invalid or unexpected")
	PublicKey *rsa.PublicKey
	publicKeyURL = "http://localhost:8080/token/publicKey"
)

func VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 서명 알고리즘 체크
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrTokenInvalid
		}
		return PublicKey, nil
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


func FetchAndParsePublicKey() (*rsa.PublicKey, error) {
	resp, err := http.Get(publicKeyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %w", err)
	}
	defer resp.Body.Close()

	keyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key body: %w", err)
	}

	pubKey, err := ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PEM public key: %w", err)
	}

	return pubKey, nil
}

func ParseRSAPublicKeyFromPEM(pemBytes []byte) (*rsa.PublicKey, error) {
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