FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/gifka-bot ./cmd/gifka-bot

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/gifka-bot ./bin
CMD ["./bin/gifka-bot"]