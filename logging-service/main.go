package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/forum_golang/micro-template/internal/pkg/discovery"
	"github.com/gofiber/fiber/v2"
)

type LogEntry struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	// Initialize Elasticsearch client
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Create index if not exists
	indexName := "logs"
	mapping := `{
		"mappings": {
			"properties": {
				"level": { "type": "keyword" },
				"message": { "type": "text" },
				"service": { "type": "keyword" },
				"timestamp": { "type": "date" }
			}
		}
	}`

	// Check if index exists
	exists, err := esClient.Indices.Exists([]string{indexName})
	if err != nil {
		log.Printf("Error checking index existence: %s", err)
	} else if exists.StatusCode == 404 {
		// Create index if not exists
		res, err := esClient.Indices.Create(
			indexName,
			esClient.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			log.Printf("Cannot create index: %s", err)
		}
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}

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
		var logEntry LogEntry
		if err := c.BodyParser(&logEntry); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to parse log entry: %v", err),
			})
		}

		// Set timestamp if not provided
		if logEntry.Timestamp.IsZero() {
			logEntry.Timestamp = time.Now()
		}

		// Convert log entry to JSON
		logJSON, err := json.Marshal(logEntry)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to marshal log entry: %v", err),
			})
		}

		// Store log in Elasticsearch
		res, err := esClient.Index(
			indexName,
			strings.NewReader(string(logJSON)),
			esClient.Index.WithRefresh("true"),
		)
		if err != nil {
			log.Printf("Error indexing log: %s", err)
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to store log: %v", err),
			})
		}
		defer res.Body.Close()

		log.Printf("[%s] %s: %s", logEntry.Service, logEntry.Level, logEntry.Message)

		return c.JSON(fiber.Map{
			"status": "logged",
			"log":    logEntry,
		})
	})

	app.Get("/logs", func(c *fiber.Ctx) error {
		// Get query parameters
		service := c.Query("service")
		level := c.Query("level")
		fromStr := c.Query("from", "0")
		sizeStr := c.Query("size", "10")

		// Convert string parameters to int
		from, err := strconv.Atoi(fromStr)
		if err != nil {
			from = 0
		}
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			size = 10
		}

		// Build search query
		query := map[string]interface{}{
			"sort": []map[string]interface{}{
				{"timestamp": map[string]interface{}{"order": "desc"}},
			},
		}

		// Add filters if provided
		if service != "" || level != "" {
			must := []map[string]interface{}{}
			if service != "" {
				must = append(must, map[string]interface{}{
					"term": map[string]interface{}{"service": service},
				})
			}
			if level != "" {
				must = append(must, map[string]interface{}{
					"term": map[string]interface{}{"level": level},
				})
			}
			query["query"] = map[string]interface{}{
				"bool": map[string]interface{}{
					"must": must,
				},
			}
		}

		// Convert query to JSON
		queryJSON, err := json.Marshal(query)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to build search query: %v", err),
			})
		}

		// Search logs in Elasticsearch
		res, err := esClient.Search(
			esClient.Search.WithIndex(indexName),
			esClient.Search.WithBody(strings.NewReader(string(queryJSON))),
			esClient.Search.WithFrom(from),
			esClient.Search.WithSize(size),
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to search logs: %v", err),
			})
		}
		defer res.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to parse search results: %v", err),
			})
		}

		return c.JSON(fiber.Map{
			"total": result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"],
			"logs":  result["hits"].(map[string]interface{})["hits"],
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
