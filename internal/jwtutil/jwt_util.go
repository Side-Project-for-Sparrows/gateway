package jwtutil

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey       = []byte("secret-key")
	expirationTime  = 10 * time.Minute
	errTokenExpired = errors.New("access token expired")
	errTokenInvalid = errors.New("access token invalid or unexpected")
)

func GenerateToken(userID int64, expirationTime time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {

		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, errTokenExpired
		}
		return 0, errTokenInvalid
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userID, parseErr := parseUserID(claims.Subject)
		if parseErr != nil {
			return 0, errTokenInvalid
		}
		return userID, nil
	}

	return 0, errTokenInvalid
}

func parseUserID(subject string) (int64, error) {
	return strconv.ParseInt(subject, 10, 64)
}
