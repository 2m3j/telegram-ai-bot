FROM golang:1.24-alpine AS base

USER root

# Configure non-root user

# Create a group and user
RUN addgroup -S app-group && adduser -S app -G app-group

RUN apk --no-cache add shadow && \
    groupmod -o -g 1000 app-group && \
    usermod -o -u 1000 -g app-group app

USER app

# Устанавливаем Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Установка Air для горячей перезагрузки
RUN go install github.com/air-verse/air@latest

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go install github.com/vektra/mockery/v2@v2.52.2

# Создаем рабочую директорию
WORKDIR /app

# Устанавливаем зависимости (опционально для кеширования модулей)
COPY ../../go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]