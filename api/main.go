package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	err = sd.RegisterService("api", 8080)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Setup routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/api/test-log", func(c *fiber.Ctx) error {
		// Create test log
		testLog := fiber.Map{
			"service":   "api",
			"level":     "info",
			"message":   "Test log from API service",
			"timestamp": time.Now(),
		}

		// Convert log to JSON
		logJSON, err := json.Marshal(testLog)
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to marshal log: %v", err))
		}

		// Send log to logging service
		// Use container name for service discovery within Docker network
		url := fmt.Sprintf("http://logging-service:8082/logs")
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(logJSON))
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to create request: %v", err))
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to send log: %v", err))
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return c.Status(resp.StatusCode).SendString("Failed to send log to logging service")
		}

		return c.SendString("Log sent successfully")
	})

	app.Get("/api/users", func(c *fiber.Ctx) error {
		// TODO: Implement user listing
		return c.JSON(fiber.Map{
			"users": []string{"user1", "user2"},
		})
	})

	app.Get("/api/posts", func(c *fiber.Ctx) error {
		// TODO: Implement post listing
		return c.JSON(fiber.Map{
			"posts": []string{"post1", "post2"},
		})
	})

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		if err := sd.DeregisterService("api-8080"); err != nil {
			log.Printf("Error deregistering service: %v", err)
		}
		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
