# Stage 1: Базовый образ для сборки приложения
FROM golang:1.19-alpine AS builder

# Добавьте необходимые шаги для сборки приложения
WORKDIR /$HOME/pipeline
COPY . .
RUN go mod download
RUN go build -o pipeline

# Stage 2: Образ для конечного приложения
FROM alpine

# Скопируйте собранные файлы приложения из предыдущего этапа
COPY --from=builder /$HOME/pipeline/ /$HOME/pipeline/

# Установите необходимые зависимости и настройки
RUN apk update && apk add --no-cache ca-certificates

# Укажите команду запуска приложения
CMD ["$HOME//pipeline"]