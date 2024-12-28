package gateway

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSConfig cấu hình CORS
type CORSConfig struct {
	Enabled        bool     `env:"CORS_ENABLED" envDefault:"true"`
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	AllowedMethods []string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders []string `env:"CORS_ALLOWED_HEADERS" envDefault:"*"`
	MaxAge         int      `env:"CORS_MAX_AGE" envDefault:"86400"`
}

// NewCORSConfig tạo cấu hình CORS mới từ environment
func NewCORSConfig() *CORSConfig {
	return &CORSConfig{
		Enabled:        true,
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         86400,
	}
}

// ConfigureCORS cấu hình CORS middleware cho Fiber app
func ConfigureCORS(app *fiber.App, config *CORSConfig) {
	if !config.Enabled {
		return
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(config.AllowedOrigins, ","),
		AllowMethods: strings.Join(config.AllowedMethods, ","),
		AllowHeaders: strings.Join(config.AllowedHeaders, ","),
		MaxAge:       config.MaxAge,
	}))
}
