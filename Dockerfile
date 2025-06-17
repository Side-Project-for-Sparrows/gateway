FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# 정적 빌드
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway

FROM scratch
COPY --from=builder /app/gateway /gateway

# ? config 파일 명시적으로 복사
COPY config /config

ENTRYPOINT ["/gateway"]
