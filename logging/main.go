package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/olivere/elastic/v7"
)

type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Service     string                 `json:"service"`
	Message     string                 `json:"message"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	TraceID     string                 `json:"trace_id,omitempty"`
	SpanID      string                 `json:"span_id,omitempty"`
	Environment string                 `json:"environment"`
}

func main() {
	// Connect to Elasticsearch
	client, err := elastic.NewClient(
		elastic.SetURL(os.Getenv("ELASTICSEARCH_URL")),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create Fiber app
	app := fiber.New()
	app.Use(logger.New())

	// Create log index if not exists
	indexName := "logs"
	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		mapping := `{
			"mappings": {
				"properties": {
					"timestamp": { "type": "date" },
					"level": { "type": "keyword" },
					"service": { "type": "keyword" },
					"message": { "type": "text" },
					"metadata": { "type": "object" },
					"trace_id": { "type": "keyword" },
					"span_id": { "type": "keyword" },
					"environment": { "type": "keyword" }
				}
			}
		}`
		_, err = client.CreateIndex(indexName).Body(mapping).Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Handle log ingestion
	app.Post("/api/logs", func(c *fiber.Ctx) error {
		var logEntry LogEntry
		if err := c.BodyParser(&logEntry); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid log entry format",
			})
		}

		logEntry.Timestamp = time.Now()

		// Index log entry
		_, err := client.Index().
			Index(indexName).
			BodyJson(logEntry).
			Do(context.Background())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to store log entry",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"message": "Log entry stored successfully",
		})
	})

	// Query logs
	app.Get("/api/logs", func(c *fiber.Ctx) error {
		service := c.Query("service")
		level := c.Query("level")
		fromStr := c.Query("from", "0")
		sizeStr := c.Query("size", "10")

		// Convert string to int
		from, err := strconv.Atoi(fromStr)
		if err != nil {
			from = 0
		}
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			size = 10
		}

		// Build search query
		query := elastic.NewBoolQuery()
		if service != "" {
			query.Must(elastic.NewTermQuery("service", service))
		}
		if level != "" {
			query.Must(elastic.NewTermQuery("level", level))
		}

		// Execute search
		searchResult, err := client.Search().
			Index(indexName).
			Query(query).
			Sort("timestamp", false).
			From(from).
			Size(size).
			Do(context.Background())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to query logs",
			})
		}

		var logs []LogEntry
		for _, hit := range searchResult.Hits.Hits {
			var log LogEntry
			err := json.Unmarshal(hit.Source, &log)
			if err != nil {
				continue
			}
			logs = append(logs, log)
		}

		return c.JSON(logs)
	})

	// Add health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	log.Fatal(app.Listen(":8082"))
}
