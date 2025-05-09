# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bin/server ./server/cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/client/static ./client/static
COPY --from=builder /app/client/templates ./client/templates

RUN apk add --no-cache tzdata

EXPOSE 8080
CMD ["./server"]