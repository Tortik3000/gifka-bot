FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/gifka-bot ./cmd/gifka-bot

# Используем образ с предустановленным ffmpeg
FROM jrottenberg/ffmpeg:4.1-alpine

WORKDIR /app
RUN mkdir -p bin
COPY --from=builder /app/bin/gifka-bot ./bin/gifka-bot

CMD ["./bin/gifka-bot"]