# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p /app/bin && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/auth-service ./auth-service/cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/product-service ./product-service/cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/order-service ./order-service/cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/user-service ./user-service/cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/chat-service ./chat-service/cmd/server && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/notification-service ./notification-service/cmd/server

FROM alpine:3.19

WORKDIR /app
ENV GIN_MODE=release \
    BACKEND_PORT=8080

COPY --from=builder /app/bin /app/bin
COPY --from=builder /app/frontend/templates /app/templates
COPY --from=builder /app/frontend/static /app/static

EXPOSE 8080

CMD ["/app/bin/auth-service"]
