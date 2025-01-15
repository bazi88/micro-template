package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/forum_golang/micro-template/internal/pkg/discovery"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Initialize service discovery
	sd, err := discovery.NewServiceDiscovery()
	if err != nil {
		log.Fatalf("Failed to create service discovery: %v", err)
	}

	// Register service with Consul
	err = sd.RegisterService("api-gateway", 80)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Setup routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Forward requests to API service
	app.All("/api/*", func(c *fiber.Ctx) error {
		// Use container name for service discovery within Docker network
		apiURL := fmt.Sprintf("http://api:8080%s", c.Path())

		// Create HTTP client
		client := &http.Client{}

		// Create new request
		req, err := http.NewRequest(
			c.Method(),
			apiURL,
			bytes.NewReader(c.Body()),
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Copy headers
		for key, values := range c.GetReqHeaders() {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Forward the request
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Response().Header.Set(key, value)
			}
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Send response
		return c.Status(resp.StatusCode).Send(body)
	})

	// Forward requests to Logging service
	app.Get("/logs/*", func(c *fiber.Ctx) error {
		loggingURL := fmt.Sprintf("http://logging-service:8082%s", c.Path())
		// Forward request to Logging service
		// TODO: Implement request forwarding
		return c.SendString("Logging service available at: " + loggingURL)
	})

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		if err := sd.DeregisterService("api-gateway-80"); err != nil {
			log.Printf("Error deregistering service: %v", err)
		}
		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	if err := app.Listen(":80"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
