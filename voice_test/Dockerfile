# Скачиваем зависимости и собираем бинарник
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник из папки sfu
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/sfu ./sfu

# Запускаем бинарник в минималистичном образе
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Копируем бинарник из первого этапа
COPY --from=builder /app/bin/sfu .

# Порт, на котором слушает приложение
EXPOSE 8081

# Команда запуска
CMD ["./sfu"]
