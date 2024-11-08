FROM golang:1.20-alpine AS builder

# Установка зависимостей
RUN apk update && apk add --no-cache git

# Создание рабочей директории
WORKDIR /app

# Копируем go.mod и go.sum, устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и собираем приложение
COPY . .
RUN go build -o main .

# Меньший базовый образ для запуска
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .

# Порт, на котором приложение будет доступно
EXPOSE 8080

CMD ["./main"]