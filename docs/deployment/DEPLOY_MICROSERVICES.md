# Hướng dẫn Deploy Microservices

## Yêu cầu hệ thống
- Docker
- Docker Compose
- Bash shell
- SSH access tới các host (nếu deploy distributed)

## Cấu trúc Microservices
Hệ thống bao gồm các service chính:
- Gateway Service (API Gateway)
- User Service
- Notification Service
- Logging Service
- Service Discovery (Consul)

## Các bước deploy

### 1. Chuẩn bị môi trường
```bash
# Copy file environment
cp .env.example .env

# Chỉnh sửa các biến môi trường trong .env
# Đặc biệt chú ý các biến:
- API_MODE=micro
- SERVICE_DISCOVERY_ENABLED=true
- CONSUL_ENABLED=true
```

### 2. Deploy Service Discovery (Consul)
```bash
# Deploy Consul trước
docker-compose -f deployments/docker-compose.yml up -d consul

# Kiểm tra Consul đã hoạt động
curl http://localhost:8500/v1/status/leader
```

### 3. Deploy từng Microservice

#### User Service
```bash
# Set các biến môi trường cho User Service
export USER_SERVICE_PORT=8081
export USER_DB_HOST=localhost
export USER_DB_PORT=5432
export USER_DB_NAME=users
export USER_REDIS_HOST=localhost
export USER_REDIS_PORT=6379

# Deploy
docker-compose -f deployments/docker-compose.user.yml up -d
```

#### Notification Service
```bash
# Set các biến môi trường cho Notification Service
export NOTIFICATION_SERVICE_PORT=8083
export NOTIFICATION_DB_HOST=localhost
export NOTIFICATION_DB_PORT=27017

# Deploy
docker-compose -f deployments/docker-compose.notification.yml up -d
```

#### Logging Service
```bash
# Set các biến môi trường cho Logging Service
export LOGGING_SERVICE_PORT=8082
export LOGGING_ES_HOST=localhost
export LOGGING_ES_PORT=9200

# Deploy
docker-compose -f deployments/docker-compose.logging.yml up -d
```

### 4. Deploy API Gateway
```bash
# Set các biến môi trường cho Gateway
export GATEWAY_MODE=auth
export GATEWAY_PORT=80
export CONSUL_ADDRESS=consul:8500

# Deploy
docker-compose -f deployments/docker-compose.gateway.yml up -d
```

## Kiểm tra hệ thống

### 1. Kiểm tra các service đã đăng ký với Consul
```bash
curl http://localhost:8500/v1/catalog/services
```

### 2. Kiểm tra health check của từng service
```bash
# User Service
curl http://localhost:8081/health

# Notification Service
curl http://localhost:8083/health

# Logging Service
curl http://localhost:8082/health

# Gateway
curl http://localhost:80/health
```

## Monitoring & Logging

### Metrics
- Prometheus endpoint: http://localhost:9090
- Grafana dashboard: http://localhost:3000

### Distributed Tracing
- Jaeger UI: http://localhost:16686

### Logs
- Logs được lưu tại thư mục `logs/` của mỗi service
- Elasticsearch + Kibana cho log aggregation

## Scaling
```bash
# Scale một service cụ thể (ví dụ: user-service)
docker-compose -f deployments/docker-compose.user.yml up -d --scale user-service=3
```

## Troubleshooting

### 1. Kiểm tra logs
```bash
# Xem logs của một service
docker-compose -f deployments/docker-compose.user.yml logs -f user-service
```

### 2. Kiểm tra service discovery
```bash
# Liệt kê các service đã đăng ký
curl http://localhost:8500/v1/catalog/services

# Kiểm tra chi tiết một service
curl http://localhost:8500/v1/catalog/service/user-service
```

### 3. Reset service
```bash
# Reset một service cụ thể
docker-compose -f deployments/docker-compose.user.yml rm -sf user-service
docker-compose -f deployments/docker-compose.user.yml up -d
```

## Backup & Restore

### Backup
```bash
# Chạy script backup
./scripts/backup.sh
```

### Restore
```bash
# Chạy script restore
./scripts/restore.sh [backup_file]
``` 