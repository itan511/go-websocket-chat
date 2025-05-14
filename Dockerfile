FROM golang:1.23-alpine as builder

RUN apk --no-cache add bash gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-websocket-chat ./cmd/main.go

# --- Production image ---
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/go-websocket-chat .

# Копируем миграции в нужную папку
COPY --from=builder /app/internal/db/migrations ./migrations

RUN chmod +x /root/go-websocket-chat

CMD ["./go-websocket-chat"]