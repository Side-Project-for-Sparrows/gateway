FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ���� ����
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway

FROM scratch
COPY --from=builder /app/gateway /gateway

# ? config ���� ��������� ����
COPY config /config

ENTRYPOINT ["/gateway"]
