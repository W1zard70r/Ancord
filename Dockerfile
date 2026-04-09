FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Финальный легковесный образ
FROM alpine:latest
WORKDIR /root/
# Копируем бинарник, конфиги и миграции
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./main"]