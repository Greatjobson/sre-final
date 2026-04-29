# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o clothes-web ./frontend/cmd/web

FROM alpine:3.19

WORKDIR /app
ENV GIN_MODE=release \
    FRONTEND_PORT=8081 \
    AUTH_SERVICE_URL=http://auth-service:8082 \
    PRODUCT_SERVICE_URL=http://product-service:8083 \
    ORDER_SERVICE_URL=http://order-service:8084 \
    USER_SERVICE_URL=http://user-service:8085 \
    CHAT_SERVICE_URL=http://chat-service:8086 \
    ANALYTICS_SERVICE_URL=http://order-service:8084

COPY --from=builder /app/clothes-web /app/clothes-web
COPY --from=builder /app/frontend/templates /app/templates
COPY --from=builder /app/frontend/static /app/static

EXPOSE 8081

CMD ["/app/clothes-web"]

