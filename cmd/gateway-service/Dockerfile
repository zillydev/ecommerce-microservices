FROM golang:1.23.4-alpine AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY /cmd/gateway-service ./cmd/gateway-service
COPY /internal/gateway ./internal/gateway
COPY /pkg ./pkg

RUN go mod tidy

RUN go build -o gateway-service ./cmd/gateway-service/main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/gateway-service .

EXPOSE 8080

CMD ./gateway-service
