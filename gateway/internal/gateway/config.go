package gateway

type Config struct {
	Port            string `env:"PORT" envDefault:"80"`
	ConsulAddress   string `env:"CONSUL_ADDRESS" envDefault:"consul:8500"`
	ServiceRegistry map[string]ServiceConfig
}

type ServiceConfig struct {
	Name     string   `json:"name"`
	URLs     []string `json:"urls"`
	Prefixes []string `json:"prefixes"`
}

func NewConfig() *Config {
	return &Config{
		ServiceRegistry: map[string]ServiceConfig{
			"user-service": {
				Name:     "user-service",
				URLs:     []string{"http://user-service:8081"},
				Prefixes: []string{"/api/users"},
			},
			"notification-service": {
				Name:     "notification-service",
				URLs:     []string{"http://notification-service:8082"},
				Prefixes: []string{"/api/notifications"},
			},
			"logging-service": {
				Name:     "logging-service",
				URLs:     []string{"http://logging-service:8082"},
				Prefixes: []string{"/api/logs"},
			},
		},
	}
}
