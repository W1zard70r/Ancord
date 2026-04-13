FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

COPY backend/ . 

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/api/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /main .
COPY --from=builder /app/backend/migrations ./migrations

EXPOSE 8080
CMD ["./main"]