FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main cmd/main.go

FROM alpine
WORKDIR /app

COPY --from=builder /app .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

EXPOSE 50051
EXPOSE 80

CMD goose -dir ./migrations postgres "$DB_DSN" up && ./main --config config/staging.yaml
