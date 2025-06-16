# --- 1단계: Go 빌드 (빌더 스테이지)
FROM golang:1.22 AS builder

WORKDIR /app

# 모듈 캐시를 활용하려면 go.mod와 go.sum만 먼저 복사
COPY go.mod go.sum ./
RUN go mod download

# 나머지 소스 복사 후 빌드
COPY . .
RUN go build -o gateway

# --- 2단계: 실행용 이미지 (작고 빠르게)
FROM alpine:latest

WORKDIR /app

# 실행파일만 복사
COPY --from=builder /app/gateway .

# 실행 포트 열기
EXPOSE 7080

# 실행 명령
ENTRYPOINT ["./gateway"]