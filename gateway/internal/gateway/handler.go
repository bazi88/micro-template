package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony/gobreaker"
)

type Handler struct {
	config *Config
	cb     map[string]*gobreaker.CircuitBreaker
}

func NewHandler(config *Config) *Handler {
	cb := make(map[string]*gobreaker.CircuitBreaker)

	// Initialize circuit breakers for each service
	for serviceName := range config.ServiceRegistry {
		cb[serviceName] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        serviceName,
			MaxRequests: 3,
			Interval:    10 * time.Second,
			Timeout:     30 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= 3 && failureRatio >= 0.6
			},
		})
	}

	return &Handler{
		config: config,
		cb:     cb,
	}
}

func (h *Handler) HandleRequest(c *fiber.Ctx) error {
	path := c.Path()

	// Find matching service
	for _, service := range h.config.ServiceRegistry {
		for _, prefix := range service.Prefixes {
			if strings.HasPrefix(path, prefix) {
				return h.proxyRequest(c, service)
			}
		}
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "Service not found",
	})
}

func (h *Handler) proxyRequest(c *fiber.Ctx, service ServiceConfig) error {
	breaker := h.cb[service.Name]

	_, err := breaker.Execute(func() (interface{}, error) {
		// Simple round-robin load balancing
		targetURL, _ := url.Parse(service.URLs[0])

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Customize proxy behavior
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Set("X-Proxy", "API Gateway")
			return nil
		}

		// Handle proxy errors
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			c.Status(http.StatusBadGateway).JSON(fiber.Map{
				"error": "Proxy error",
			})
		}

		return nil, nil
	})

	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Service unavailable",
		})
	}

	return nil
}
