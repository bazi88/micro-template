package gateway

import (
	"io"
	"net/http"
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

	// Handle health check endpoint
	if path == "/health" {
		return c.SendString("OK")
	}

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

	result, err := breaker.Execute(func() (interface{}, error) {
		// Simple round-robin load balancing
		targetURL, err := url.Parse(service.URLs[0])
		if err != nil {
			return nil, err
		}

		// Create HTTP request
		req, err := http.NewRequest(c.Method(), targetURL.String()+c.Path(), strings.NewReader(string(c.Body())))
		if err != nil {
			return nil, err
		}

		// Copy headers
		c.Request().Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Response().Header.Set(key, value)
			}
		}

		// Copy response status
		c.Status(resp.StatusCode)

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})

	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Service unavailable",
		})
	}

	if body, ok := result.([]byte); ok {
		return c.Send(body)
	}

	return nil
}
