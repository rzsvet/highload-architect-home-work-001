# Build stage
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости
RUN apk add --no-cache git

WORKDIR /app

# Копируем только файлы модулей для лучшего кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем статический бинарник с оптимизацией
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s -extldflags '-static'" \
    -o /api \
    ./cmd/server
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
#     go build \
#     -ldflags="-w -s -X main.buildVersion=$(git describe --tags --always || echo 'dev') \
#                -X main.buildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
#                -X main.buildCommit=$(git rev-parse HEAD || echo 'unknown') \
#                -extldflags '-static'" \
#     -trimpath \
#     -o /api \
#     ./cmd/server

# Создаем минимальный образ
FROM scratch

# Копируем необходимые файлы из builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /api /api

# Необязательно: копируем документацию если нужна
COPY --from=builder /app/docs /docs

# Переменные окружения
ENV TZ=UTC \
    SERVER_PORT=8080

# Порт приложения
EXPOSE 8080

# Здоровье проверка
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/api", "healthcheck"] || exit 1

# Запуск приложения
ENTRYPOINT ["/api"]