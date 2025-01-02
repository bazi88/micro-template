# Hướng Dẫn Deploy Monolithic

## Cấu trúc
Hệ thống monolithic bao gồm các thành phần chính:
- App chính (monolithic service)
- PostgreSQL database
- Redis cache
- Network bridge cho kết nối giữa các service

## File Cấu Hình
### 1. File Docker Compose
File `deployments/docker-compose-monolithic.yml` định nghĩa:
- App service với các cấu hình về resource limits và volumes
- PostgreSQL service với healthcheck
- Redis service với healthcheck và password
- Network bridge cho kết nối nội bộ
- Các volume để lưu trữ dữ liệu

### 2. File Environment
File `.env.monolithic` chứa các biến môi trường:
```env
# Deployment Configuration
DEPLOY_MODE=monolithic
GATEWAY_MODE=monolithic

# API Gateway
GATEWAY_PORT=80

# Database Configuration
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=postgres
DB_PORT=5432

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis

# Logging Configuration
LOG_LEVEL=info

# Performance Configuration
CPU_LIMIT=0.5
MEMORY_LIMIT=512M
CPU_RESERVATION=0.25
MEMORY_RESERVATION=256M
```

## Các Bước Deploy

### 1. Chuẩn Bị Môi Trường
```bash
# Copy file environment
cp .env.monolithic .env
```

### 2. Tạo Thư Mục Logs và Config
```bash
# Tạo thư mục logs
mkdir -p logs

# Tạo thư mục config
mkdir -p config
```

### 3. Deploy Hệ Thống
```bash
# Build và start các services
docker-compose -f deployments/docker-compose-monolithic.yml up --build
```

### 4. Kiểm Tra Hoạt Động
- API Gateway: http://localhost:80
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## Monitoring và Logging
- Logs được lưu trong thư mục `logs`
- Cấu hình được lưu trong thư mục `config`
- Health check tự động cho database và redis

## Scaling và Performance
- CPU và Memory được giới hạn thông qua các biến môi trường
- Có thể điều chỉnh trong file .env.monolithic:
  - CPU_LIMIT và CPU_RESERVATION
  - MEMORY_LIMIT và MEMORY_RESERVATION

## Khác Biệt với Microservices
1. Không sử dụng:
   - Service Discovery (Consul)
   - Load Balancer (Nginx)
   - Multiple Database Init Script
   
2. Database:
   - Sử dụng một database duy nhất thay vì nhiều database cho từng service
   - Cấu trúc đơn giản hơn

3. Deployment:
   - Đơn giản hơn, không cần orchestration
   - Không cần service discovery
   - Không cần load balancing

## Troubleshooting
1. Lỗi Permission:
```bash
# Sử dụng sudo nếu gặp lỗi permission
sudo docker-compose -f deployments/docker-compose-monolithic.yml up --build
```

2. Lỗi Port Conflict:
- Kiểm tra và đảm bảo các port không bị conflict:
  - 80: API Gateway
  - 5432: PostgreSQL
  - 6379: Redis

3. Lỗi Volume:
```bash
# Xóa volumes cũ nếu cần
docker-compose -f deployments/docker-compose-monolithic.yml down -v
``` 