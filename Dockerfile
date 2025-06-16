# --- 1�ܰ�: Go ���� (���� ��������)
FROM golang:1.24 AS builder

WORKDIR /app

# ��� ĳ�ø� Ȱ���Ϸ��� go.mod�� go.sum�� ���� ����
COPY go.mod go.sum ./
RUN go mod download

# ������ �ҽ� ���� �� ����
COPY . .
RUN go build -o gateway

# --- 2�ܰ�: ����� �̹��� (�۰� ������)
FROM alpine:latest

WORKDIR /app

# �������ϸ� ����
COPY --from=builder /app/gateway .

# ���� ��Ʈ ����
EXPOSE 7080

# stg ȯ������ �����ؼ� ȯ�漳�� ����
ENV APP_ENV=stg

# ���� ���
ENTRYPOINT ["./gateway"]