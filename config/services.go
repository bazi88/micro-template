package config

import (
	"os"
	"strings"
)

type ServiceConfig struct {
	Name         string   `json:"name"`
	Enabled      bool     `json:"enabled"`
	LoadBalancer bool     `json:"load_balancer"`
	URLs         []string `json:"urls"`
	Prefixes     []string `json:"prefixes"`
	AuthRequired bool     `json:"auth_required"`
}

type ServicesConfig struct {
	ActiveServices []string                  `json:"active_services"`
	Services       map[string]*ServiceConfig `json:"services"`
}

func NewServicesConfig() *ServicesConfig {
	config := &ServicesConfig{
		Services: make(map[string]*ServiceConfig),
	}

	// Đọc danh sách services được kích hoạt từ biến môi trường
	if activeServices := os.Getenv("ACTIVE_SERVICES"); activeServices != "" {
		config.ActiveServices = strings.Split(activeServices, ",")
	}

	// Khởi tạo cấu hình cho từng service
	for _, serviceName := range config.ActiveServices {
		serviceEnvPrefix := strings.ToUpper(serviceName)

		config.Services[serviceName] = &ServiceConfig{
			Name:         serviceName,
			Enabled:      true,
			LoadBalancer: os.Getenv(serviceEnvPrefix+"_LOAD_BALANCER") == "true",
			URLs:         strings.Split(os.Getenv(serviceEnvPrefix+"_URLS"), ","),
			Prefixes:     strings.Split(os.Getenv(serviceEnvPrefix+"_PREFIXES"), ","),
			AuthRequired: os.Getenv(serviceEnvPrefix+"_AUTH_REQUIRED") == "true",
		}
	}

	return config
}
