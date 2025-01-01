package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheStore định nghĩa interface cho caching
type CacheStore interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// RedisCache implementation của CacheStore sử dụng Redis
type RedisCache struct {
	client *redis.Client
}

// Config cấu hình cho Redis cache
type Config struct {
	Host        string        `env:"REDIS_HOST" envDefault:"redis"`
	Port        int           `env:"REDIS_PORT" envDefault:"6379"`
	Password    string        `env:"REDIS_PASSWORD" envDefault:""`
	DB          int           `env:"REDIS_DB" envDefault:"0"`
	DefaultTTL  time.Duration `env:"CACHE_DEFAULT_TTL" envDefault:"3600s"`
	Compression bool          `env:"CACHE_COMPRESSION_ENABLED" envDefault:"true"`
	Strategy    string        `env:"CACHE_STRATEGY" envDefault:"distributed"`
}

// NewRedisCache tạo instance mới của RedisCache
func NewRedisCache(config *Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisCache{
		client: client,
	}, nil
}

// Get lấy giá trị từ cache
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

// Set lưu giá trị vào cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var data []byte
	var err error

	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// Delete xóa key khỏi cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Clear xóa toàn bộ cache
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Close đóng kết nối Redis
func (c *RedisCache) Close() error {
	return c.client.Close()
}
