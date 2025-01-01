# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final stage
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./api"]
