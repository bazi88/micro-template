FROM golang:1.21-alpine

WORKDIR /app

# Copy the source code
COPY . .

# Initialize module and download dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 