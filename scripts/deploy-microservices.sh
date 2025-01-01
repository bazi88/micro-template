#!/bin/bash

# Deploy script for microservices mode

# Load environment variables
set -a
source .env
set +a

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose is required but not installed. Aborting." >&2; exit 1; }

# Set environment for microservices mode
export API_MODE=micro
export SERVICE_DISCOVERY_ENABLED=true
export CONSUL_ENABLED=true

# Function to deploy a service
deploy_service() {
    local service=$1
    local host=$3
    
    echo "Deploying $service to $host..."
    
    # Copy necessary files
    ssh $host "mkdir -p ~/micro"
    scp -r . $host:~/micro/
    
    # Deploy service
    ssh $host "cd ~/micro && docker-compose -f deployments/docker-compose.${service}.yml pull && docker-compose -f deployments/docker-compose.${service}.yml up -d"
    
    # Wait for service to be healthy
    echo "Waiting for $service to be healthy..."
    while true; do
        status=$(ssh $host "cd ~/micro && docker-compose -f deployments/docker-compose.${service}.yml ps -q ${service}-service | xargs docker inspect -f '{{.State.Health.Status}}'")
        if [ "$status" = "healthy" ]; then
            echo "$service is healthy"
            break
        fi
        echo "Waiting for $service to be healthy..."
        sleep 5
    done
}

# Deploy Consul first
echo "Deploying Consul..."
docker-compose -f deployments/docker-compose.microservices.yml up -d consul

# Wait for Consul to be ready
echo "Waiting for Consul to be ready..."
while ! curl -s http://localhost:8500/v1/status/leader > /dev/null; do
    echo "Waiting for Consul..."
    sleep 5
done

# Deploy API Gateway
echo "Deploying API Gateway..."
docker-compose -f deployments/docker-compose.microservices.yml up -d api-gateway nginx

# Deploy individual services
deploy_service "user" "${USER_SERVICE_HOST}"
deploy_service "notification" "${NOTIFICATION_SERVICE_HOST}"

# Deploy monitoring stack
echo "Deploying monitoring stack..."
docker-compose -f deployments/docker-compose.microservices.yml up -d prometheus grafana jaeger

echo "Microservices deployment completed successfully!"
echo "You can access the services at:"
echo "- API Gateway: http://localhost:80"
echo "- Consul UI: http://localhost:8500"
echo "- Prometheus: http://localhost:9090"
echo "- Grafana: http://localhost:3000"
echo "- Jaeger UI: http://localhost:16686" 