package gateway

type Config struct {
	Port            string `env:"GATEWAY_PORT" envDefault:"80"`
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
			"logging-service": {
				Name:     "logging-service",
				URLs:     []string{"http://logging-service:8082"},
				Prefixes: []string{"/api/logs"},
			},
		},
	}
}
