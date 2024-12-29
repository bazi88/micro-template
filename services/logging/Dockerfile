FROM golang:1.21-alpine

WORKDIR /app

# Copy the source code
COPY . .

# Initialize module and download dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Create log directory
RUN mkdir -p /var/log/micro

# Expose port
EXPOSE 8082

# Command to run the executable
CMD ["./main"] 