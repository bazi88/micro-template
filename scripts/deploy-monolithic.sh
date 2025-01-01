#!/bin/bash

# Deploy script for monolithic mode

# Load environment variables
set -a
source .env
set +a

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose is required but not installed. Aborting." >&2; exit 1; }

# Set environment for monolithic mode
export API_MODE=monolithic
export SERVICE_DISCOVERY_ENABLED=false
export CONSUL_ENABLED=false

# Backup existing data if exists
if [ -d "./data" ]; then
    echo "Backing up existing data..."
    timestamp=$(date +%Y%m%d_%H%M%S)
    mkdir -p ./backups
    tar -czf "./backups/data_backup_$timestamp.tar.gz" ./data
fi

# Pull latest images
echo "Pulling latest images..."
docker-compose -f docker-compose.monolithic.yml pull

# Start services
echo "Starting services..."
docker-compose -f docker-compose.monolithic.yml up -d

# Wait for services to be healthy
echo "Waiting for services to be healthy..."
for service in api user-service notification-service logging-service; do
    echo "Checking $service..."
    while true; do
        status=$(docker-compose -f docker-compose.monolithic.yml ps -q $service | xargs docker inspect -f '{{.State.Health.Status}}')
        if [ "$status" = "healthy" ]; then
            echo "$service is healthy"
            break
        fi
        echo "Waiting for $service to be healthy..."
        sleep 5
    done
done

# Run database migrations if needed
echo "Running database migrations..."
docker-compose -f docker-compose.monolithic.yml exec api ./migrate up
docker-compose -f docker-compose.monolithic.yml exec user-service ./migrate up
docker-compose -f docker-compose.monolithic.yml exec notification-service ./migrate up

echo "Monolithic deployment completed successfully!"
echo "You can access the services at:"
echo "- API Gateway: http://localhost:80"
echo "- API Service: http://localhost:8080"
echo "- User Service: http://localhost:8081"
echo "- Notification Service: http://localhost:8083"
echo "- Logging Service: http://localhost:8082" 