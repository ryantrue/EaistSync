# Stage 1: Сборка приложения
FROM golang:1.24-alpine AS builder

LABEL maintainer="zufirka@icloud.com"

RUN apk update && \
    apk add --no-cache git build-base ca-certificates && \
    update-ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /build/cmd/eaistsync
RUN CGO_ENABLED=1 go build -o eaistsync .

# Stage 2: Финальный образ
FROM alpine:latest

LABEL maintainer="zufirka@icloud.com"

# Устанавливаем ca-certificates и wget (для HEALTHCHECK)
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Копируем собранный бинарник из предыдущего этапа
COPY --from=builder /build/cmd/eaistsync/eaistsync .

# Копируем миграции (если нужны)
COPY --from=builder /build/migrations ./migrations

EXPOSE 8080

RUN adduser -D appuser
USER appuser

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./eaistsync"]