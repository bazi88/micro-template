# Go-micro-template

Go-micro-template là một boilerplate cho microservices sử dụng Golang, được thiết kế với kiến trúc module hóa cao, cho phép chạy độc lập từng service hoặc kết hợp nhiều services tùy theo nhu cầu.

## Tính năng

### Core Services
- **API Service**: RESTful API service với khả năng chạy standalone hoặc microservice mode
- **Gateway Service**: API Gateway với load balancing và routing
- **Logging Service**: Centralized logging service, có thể chạy độc lập hoặc kết hợp với ELK stack

### Infrastructure Services (Optional)
- **Database**: PostgreSQL/MongoDB (tùy chọn)
- **Cache**: Redis (tùy chọn)
- **Message Queue**: RabbitMQ/Kafka (tùy chọn)
- **Service Discovery**: Consul (tùy chọn)

### Monitoring & Logging Stack (Optional)
- **ELK Stack**: Elasticsearch, Logstash, Kibana
- **Monitoring Stack**: Prometheus, Grafana, Alertmanager

## Cấu trúc dự án

```
go-micro-template/
├── services/              # Core Services
│   ├── api/              # API Service
│   ├── gateway/          # API Gateway
│   └── logging/          # Logging Service
├── infrastructure/        # Optional Infrastructure
│   ├── database/         # Database Services
│   ├── cache/            # Cache Service
│   ├── elk/              # ELK Stack
│   └── monitoring/       # Monitoring Stack
└── docker-compose.yml    # Main compose file
```

## Yêu cầu

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.21+ (cho development)

## Hướng dẫn sử dụng

### 1. API Service

#### Standalone Mode
```bash
cd services/api
cat > .env << EOL
API_MODE=standalone
DB_ENABLED=false
CACHE_ENABLED=false
DISCOVERY_ENABLED=false
EOL
docker-compose up -d
```

#### With Database
```bash
cd services/api
cat > .env << EOL
API_MODE=micro
DB_ENABLED=true
DB_TYPE=postgres
DB_HOST=postgres
EOL
docker-compose -f docker-compose.yml -f ../../infrastructure/database/postgres/docker-compose.yml up -d
```

### 2. Logging Service

#### File-based Logging
```bash
cd services/logging
cat > .env << EOL
ES_ENABLED=false
LOG_OUTPUT=file
EOL
docker-compose up -d
```

#### With ELK Stack
```bash
cd infrastructure/elk
cat > .env << EOL
ES_VERSION=7.9.3
ES_PORT=9200
ES_JVM_MIN=512m
ES_JVM_MAX=512m
ELASTIC_USERNAME=elastic
ELASTIC_PASSWORD=changeme
KIBANA_VERSION=7.9.3
KIBANA_PORT=5601
EOL
docker-compose up -d
```

### 3. Monitoring Stack

```bash
cd infrastructure/monitoring
cat > .env << EOL
# Prometheus
PROMETHEUS_VERSION=v2.45.0
PROMETHEUS_PORT=9090
PROMETHEUS_RETENTION_TIME=15d

# Grafana
GRAFANA_VERSION=10.0.3
GRAFANA_PORT=3000
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=admin123
GRAFANA_ALLOW_SIGNUP=false

# Alertmanager
ALERTMANAGER_VERSION=v0.25.0
ALERTMANAGER_PORT=9093
EOL
docker-compose up -d
```

## Truy cập các Service

### Core Services
- API Service: http://localhost:8080
- Gateway: http://localhost:80
- Logging Service: http://localhost:8082

### Monitoring & Logging
- Elasticsearch: http://localhost:9200
- Kibana: http://localhost:5601
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
- Alertmanager: http://localhost:9093

## Scaling Services

### Horizontal Scaling
```bash
# Scale API service
docker-compose up -d --scale api=3

# Scale Logging service
docker-compose up -d --scale logging=2
```

### Resource Limits
```bash
# API Service
API_CPU_LIMIT=0.5
API_MEMORY_LIMIT=512M

# Logging Service
LOGGING_CPU_LIMIT=0.3
LOGGING_MEMORY_LIMIT=256M
```

## Environment Variables

### API Service
| Variable | Default | Description |
|----------|---------|-------------|
| API_MODE | standalone | standalone/micro |
| DB_ENABLED | false | Enable database |
| CACHE_ENABLED | false | Enable Redis cache |
| DISCOVERY_ENABLED | false | Enable service discovery |

### Logging Service
| Variable | Default | Description |
|----------|---------|-------------|
| ES_ENABLED | false | Enable Elasticsearch |
| LOG_OUTPUT | file | file/elasticsearch |
| LOG_FORMAT | json | json/text |

### ELK Stack
| Variable | Default | Description |
|----------|---------|-------------|
| ES_VERSION | 7.9.3 | Elasticsearch version |
| ES_PORT | 9200 | Elasticsearch port |
| KIBANA_VERSION | 7.9.3 | Kibana version |
| KIBANA_PORT | 5601 | Kibana port |

### Monitoring Stack
| Variable | Default | Description |
|----------|---------|-------------|
| PROMETHEUS_VERSION | v2.45.0 | Prometheus version |
| PROMETHEUS_PORT | 9090 | Prometheus port |
| GRAFANA_VERSION | 10.0.3 | Grafana version |
| GRAFANA_PORT | 3000 | Grafana port |
| ALERTMANAGER_VERSION | v0.25.0 | Alertmanager version |
| ALERTMANAGER_PORT | 9093 | Alertmanager port |

## Development

1. Clone repository:
```bash
git clone https://github.com/yourusername/go-micro-template.git
cd go-micro-template
```

2. Start specific service:
```bash
cd services/api  # hoặc services khác
cp .env.example .env
docker-compose up -d
```

3. Development với hot-reload:
```bash
cd services/api
go run main.go
```

## Testing

```bash
# Unit tests
go test ./...

# Integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## Production Deployment

1. Cấu hình security:
```bash
# Đặt passwords mạnh
ELASTIC_PASSWORD=your-secure-password
GRAFANA_ADMIN_PASSWORD=your-secure-password

# Bật TLS/SSL
ENABLE_SSL=true
SSL_CERT_PATH=/path/to/cert
```

2. Cấu hình monitoring:
```bash
# Tăng retention time
PROMETHEUS_RETENTION_TIME=30d

# Cấu hình alerts
ALERTMANAGER_CONFIG_PATH=/path/to/config
```

3. Deploy:
```bash
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