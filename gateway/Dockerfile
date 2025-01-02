FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 80

# Copy config file
COPY config/config.yaml /app/config/config.yaml

CMD ["./main"] 