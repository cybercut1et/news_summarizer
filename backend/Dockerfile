# Dockerfile для Go бэкенда
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Устанавливаем Python и pip
RUN apk --no-cache add python3 py3-pip sqlite

WORKDIR /app

# Копируем скомпилированное приложение
COPY --from=builder /app/main .

# Делаем исполняемым
RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]