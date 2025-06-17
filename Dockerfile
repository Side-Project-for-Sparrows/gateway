FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gateway

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway .
RUN chmod +x ./gateway
EXPOSE 7080
ENV APP_ENV=stg

ENTRYPOINT ["./gateway"]