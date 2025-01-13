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
	err = sd.RegisterService("logging-service", 8082)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Setup routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/logs", func(c *fiber.Ctx) error {
		var logEntry struct {
			Level   string `json:"level"`
			Message string `json:"message"`
			Service string `json:"service"`
		}
		if err := c.BodyParser(&logEntry); err != nil {
			return err
		}

		// TODO: Store log in Elasticsearch
		log.Printf("[%s] %s: %s", logEntry.Service, logEntry.Level, logEntry.Message)

		return c.JSON(fiber.Map{
			"status": "logged",
		})
	})

	app.Get("/logs", func(c *fiber.Ctx) error {
		// TODO: Implement log retrieval from Elasticsearch
		return c.JSON(fiber.Map{
			"logs": []string{},
		})
	})

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		if err := sd.DeregisterService("logging-service-8082"); err != nil {
			log.Printf("Error deregistering service: %v", err)
		}
		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	if err := app.Listen(":8082"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
