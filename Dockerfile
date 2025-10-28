FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/gifka-bot ./cmd/gifka-bot

FROM alpine:latest

WORKDIR /app

RUN apk update && apk add --no-cache \
    ffmpeg \
    && rm -rf /var/cache/apk/*

RUN mkdir -p bin
COPY --from=builder /app/bin/gifka-bot ./bin/gifka-bot
COPY --from=builder /app/text_style ./text_style

CMD ["./bin/gifka-bot"]