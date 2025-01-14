package main

import (
	"log"
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
	app.Get("/api/*", func(c *fiber.Ctx) error {
		apiURL, err := sd.GetServiceURL("api")
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("API service not available")
		}
		// Forward request to API service
		// TODO: Implement request forwarding
		return c.SendString("API service available at: " + apiURL)
	})

	// Forward requests to Logging service
	app.Get("/logs/*", func(c *fiber.Ctx) error {
		loggingURL, err := sd.GetServiceURL("logging-service")
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("Logging service not available")
		}
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
