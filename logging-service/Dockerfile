FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o logging-service ./logging-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/logging-service .

EXPOSE 8082

CMD ["./logging-service"]
