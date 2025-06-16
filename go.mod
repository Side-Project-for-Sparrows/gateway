module github.com/Side-Project-for-Sparrows/gateway

replace github.com/Side-Project-for-Sparrows/gateway => ./

go 1.24.1

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
)
