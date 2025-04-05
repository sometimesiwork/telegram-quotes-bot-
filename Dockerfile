# Сборочный образ
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем модули
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Проверяем содержимое директории
RUN ls -l ./cmd

# Собираем проект
RUN mkdir -p cmd && CGO_ENABLED=0 GOOS=linux go build -o cmd/bot ./cmd

# Финальный образ
FROM alpine:latest

# Копируем исполняемый файл
COPY --from=builder /app/cmd/bot /bot

# Запускаем бота
CMD ["/bot"]