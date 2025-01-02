package gateway

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Gateway struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"gateway"`
	ConsulAddress   string                   `mapstructure:"consul_address"`
	ServiceRegistry map[string]ServiceConfig `mapstructure:"services"`
	Auth            struct {
		Enabled       bool     `mapstructure:"enabled"`
		JWTSecret     string   `mapstructure:"jwt_secret"`
		JWTExpiration int      `mapstructure:"jwt_expiration"`
		APIKeyHeader  string   `mapstructure:"api_key_header"`
		APIKeySecret  string   `mapstructure:"api_key_secret"`
		ExcludedPaths []string `mapstructure:"excluded_paths"`
		Permissions   struct {
			CacheTTL int `mapstructure:"cache_ttl"`
		} `mapstructure:"permissions"`
	} `mapstructure:"auth"`
	CORS CORSConfig `mapstructure:"cors"`
}

type ServiceConfig struct {
	Name         string   `mapstructure:"name"`
	URLs         []string `mapstructure:"urls"`
	Prefixes     []string `mapstructure:"prefixes"`
	AuthRequired bool     `mapstructure:"auth_required"`
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}

	return &config
}
