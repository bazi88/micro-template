# Go Microservices Template

## Mô hình triển khai

### 1. Monolithic Mode (Single Machine)

Trong mô hình monolithic, tất cả services chạy trên cùng một máy chủ. Phù hợp cho development và ứng dụng nhỏ.

#### 1.1. Cấu hình Monolithic
```bash
# Copy file môi trường mẫu
cp env.example .env

# Cấu hình cho monolithic mode
API_MODE=monolithic
SERVICE_DISCOVERY_ENABLED=false
CONSUL_ENABLED=false

# Cấu hình database và cache local
DB_HOST=localhost
REDIS_HOST=localhost
ELASTICSEARCH_URL=http://localhost:9200
```

#### 1.2. Khởi động Monolithic
```bash
# Khởi động toàn bộ stack
docker-compose -f docker-compose.monolithic.yml up -d

# Kiểm tra services
docker-compose -f docker-compose.monolithic.yml ps
```

#### 1.3. Scale trong Monolithic
```bash
# Scale một service cụ thể
docker-compose -f docker-compose.monolithic.yml up -d --scale api=3

# Cấu hình resource limits
API_CPU_LIMIT=0.5
API_MEMORY_LIMIT=512M
```

### 2. Microservices Mode (Multiple Machines)

Trong mô hình microservices, các services được phân tán trên nhiều máy chủ khác nhau. Phù hợp cho production và ứng dụng lớn.

#### 2.1. Setup Service Discovery (Consul Server)
```bash
# Trên máy chủ Consul
cp env.example .env

# Cấu hình Consul
CONSUL_ENABLED=true
CONSUL_ADDRESS=192.168.1.100:8500  # IP của Consul server

# Khởi động Consul
docker-compose -f docker-compose.microservices.yml up -d consul
```

#### 2.2. Setup API Gateway
```bash
# Trên máy chủ Gateway
cp env.example .env

# Cấu hình Gateway
GATEWAY_MODE=micro
API_MODE=micro
SERVICE_DISCOVERY_ENABLED=true
CONSUL_ADDRESS=192.168.1.100:8500

# Khởi động Gateway
docker-compose -f docker-compose.microservices.yml up -d api-gateway nginx
```

#### 2.3. Setup Individual Services

##### User Service (Máy 1)
```bash
# Copy config
cp env.example .env

# Cấu hình User Service
ACTIVE_SERVICES=user
USER_SERVICE_NAME=user-service
USER_SERVICE_HOST=192.168.1.10
USER_SERVICE_PORT=8081
USER_DB_HOST=localhost
USER_REDIS_HOST=localhost

# Service Discovery
CONSUL_ENABLED=true
CONSUL_ADDRESS=192.168.1.100:8500

# Khởi động User Service
docker-compose -f docker-compose.user.yml up -d
```

##### Notification Service (Máy 2)
```bash
# Copy config
cp env.example .env

# Cấu hình Notification Service
ACTIVE_SERVICES=notification
NOTIFICATION_SERVICE_NAME=notification-service
NOTIFICATION_SERVICE_HOST=192.168.1.11
NOTIFICATION_SERVICE_PORT=8083
NOTIFICATION_DB_HOST=localhost
NOTIFICATION_REDIS_HOST=localhost

# Service Discovery
CONSUL_ENABLED=true
CONSUL_ADDRESS=192.168.1.100:8500

# Khởi động Notification Service
docker-compose -f docker-compose.notification.yml up -d
```

#### 2.4. Monitoring Setup (Optional)
```bash
# Trên máy chủ monitoring
cp env.example .env

# Cấu hình monitoring
PROMETHEUS_ENABLED=true
GRAFANA_ENABLED=true
JAEGER_ENABLED=true

# Khởi động monitoring stack
docker-compose -f docker-compose.microservices.yml up -d prometheus grafana jaeger
```

### 3. Chuyển đổi giữa các mode

#### 3.1. Từ Monolithic sang Microservices
1. Backup dữ liệu từ monolithic database
2. Dừng monolithic stack
3. Setup Consul server
4. Migrate database cho từng service
5. Khởi động từng service riêng biệt
6. Cập nhật DNS/Load balancer

#### 3.2. Từ Microservices sang Monolithic
1. Backup dữ liệu từ tất cả services
2. Dừng tất cả microservices
3. Merge databases
4. Khởi động monolithic stack
5. Restore dữ liệu

### 4. Health Check và Monitoring

#### 4.1. Monolithic Health Check
```bash
# Kiểm tra tất cả services
curl localhost/health

# Kiểm tra service cụ thể
curl localhost:8080/health  # API
curl localhost:8081/health  # User
curl localhost:8083/health  # Notification
```

#### 4.2. Microservices Health Check
```bash
# Kiểm tra service discovery
curl http://consul:8500/v1/health/service/user-service

# Kiểm tra gateway
curl http://gateway/health

# Kiểm tra individual services
curl http://192.168.1.10:8081/health  # User Service
curl http://192.168.1.11:8083/health  # Notification Service
```

### 5. Troubleshooting

#### 5.1. Monolithic Troubleshooting
```bash
# Xem logs của tất cả services
docker-compose -f docker-compose.monolithic.yml logs

# Xem logs của service cụ thể
docker-compose -f docker-compose.monolithic.yml logs api
```

#### 5.2. Microservices Troubleshooting
```bash
# Kiểm tra service discovery
curl http://consul:8500/v1/catalog/services

# Kiểm tra service registration
curl http://consul:8500/v1/agent/services

# Xem logs của service cụ thể
docker-compose -f docker-compose.user.yml logs  # Trên máy User Service
docker-compose -f docker-compose.notification.yml logs  # Trên máy Notification Service
```

### 6. Lưu ý quan trọng

1. **Monolithic Mode**:
   - Dễ dàng setup và debug
   - Phù hợp cho development
   - Resource sharing giữa các services
   - Giới hạn về scale

2. **Microservices Mode**:
   - Phức tạp hơn trong setup
   - Cần quản lý network giữa các services
   - Độc lập trong scale và deploy
   - Cần monitoring tốt
   - Phù hợp cho production

3. **Security**:
   - Monolithic: Tập trung vào bảo mật external
   - Microservices: Cần bảo mật cả internal và external

4. **Backup & Recovery**:
   - Monolithic: Backup tập trung
   - Microservices: Backup phân tán