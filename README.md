# Go-micro-template

Go-micro-template là một boilerplate cho microservices sử dụng Golang, được thiết kế với kiến trúc module hóa cao và dễ dàng tùy chỉnh.

## Tính năng

- **API Gateway**: Reverse proxy, load balancing, và routing
- **Service Discovery**: Sử dụng Consul
- **Database**: Hỗ trợ PostgreSQL và MongoDB
- **Caching**: Redis cache với nhiều chiến lược
- **Authentication**: JWT và API Key
- **Rate Limiting**: Global và per-client
- **Circuit Breaker**: Fault tolerance
- **Logging**: Centralized logging với ELK stack
- **Monitoring**: Prometheus và Grafana
- **Tracing**: Distributed tracing với Jaeger
- **Migration**: Database migration và seeding
- **Hot Reload**: Tự động reload khi code thay đổi
- **Security**: CORS, rate limiting, và các best practices

## Cấu trúc dự án

```
.
├── api/            # API service
├── gateway/        # API Gateway
├── logging/        # Logging service
├── cache/          # Cache layer
├── database/       # Database layer
├── config/         # Configuration
├── docker/         # Docker files
└── scripts/        # Helper scripts
```

## Yêu cầu

- Go 1.21+
- Docker và Docker Compose
- Make (optional)

## Cài đặt

1. Clone repository:
```bash
git clone https://github.com/yourusername/go-micro-template.git
cd go-micro-template
```

2. Sao chép file môi trường:
```bash
cp env.example .env
```

3. Khởi động hệ thống:
```bash
docker-compose up -d
```

## Cấu hình

### Database

Chọn loại database trong `.env`:

```env
# PostgreSQL
DB_TYPE=postgres
DB_HOST=postgres
DB_PORT=5432
DB_NAME=micro
DB_USER=postgres
DB_PASSWORD=your-password

# MongoDB
DB_TYPE=mongodb
MONGO_URI=mongodb://mongodb:27017
MONGO_DB_NAME=micro
```

### Cache

Cấu hình Redis cache:

```env
REDIS_ENABLED=true
REDIS_HOST=redis
REDIS_PORT=6379
CACHE_STRATEGY=distributed
```

### Authentication

Bật/tắt xác thực:

```env
AUTH_ENABLED=true
JWT_SECRET=your-secret
API_KEY_ENABLED=true
```

### Service Discovery

Cấu hình Consul:

```env
CONSUL_ADDRESS=consul:8500
CONSUL_DATACENTER=dc1
```

## API Endpoints

### API Gateway
- `GET /health`: Health check
- `GET /metrics`: Prometheus metrics

### API Service
- `GET /api/v1/users`: List users
- `POST /api/v1/users`: Create user
- `GET /api/v1/users/:id`: Get user
- `PUT /api/v1/users/:id`: Update user
- `DELETE /api/v1/users/:id`: Delete user

### Logging Service
- `POST /api/logs`: Store log
- `GET /api/logs`: Query logs

## Monitoring

- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Kibana: http://localhost:5601
- Consul UI: http://localhost:8500

## Development

1. Bật hot reload:
```env
HOT_RELOAD_ENABLED=true
HOT_RELOAD_INTERVAL=1s
```

2. Chạy migrations:
```bash
make migrate
```

3. Seed database:
```bash
make seed
```

## Testing

1. Unit tests:
```bash
make test
```

2. Integration tests:
```bash
make test-integration
```

## Production

1. Cập nhật các biến môi trường:
- Đặt passwords mạnh
- Bật SSL/TLS
- Cấu hình rate limiting
- Bật authentication

2. Build và deploy:
```bash
make build
docker-compose -f docker-compose.prod.yml up -d
```

## Contributing

1. Fork repository
2. Tạo feature branch
3. Commit changes
4. Push to branch
5. Tạo Pull Request

## License

MIT License - see [LICENSE](LICENSE) file

## Authors

- Your Name - Initial work