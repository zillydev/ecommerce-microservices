FROM golang:1.23.4-alpine AS builder

# Install migration tool
RUN apk add --no-cache curl && \
  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
  mv migrate /usr/local/bin/

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY /cmd/order-service ./cmd/order-service
COPY /internal/order ./internal/order
COPY /pkg ./pkg

RUN go mod tidy

RUN go build -o order-service ./cmd/order-service/main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/order-service .
COPY --from=builder /app/internal/order/migrations /migrations
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

ENV DB_URL=postgres://user:pass@order-db:5432/orderdb?sslmode=disable

EXPOSE 8084

CMD migrate -path=/migrations -database=$DB_URL up && ./order-service
