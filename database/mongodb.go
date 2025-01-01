package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
	config *Config
}

func NewMongoDB(config *Config) *MongoDB {
	return &MongoDB{
		config: config,
	}
}

func (m *MongoDB) Connect() error {
	clientOptions := options.Client().ApplyURI(m.config.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to mongodb: %v", err)
	}

	m.client = client
	m.db = client.Database(m.config.Database)
	return nil
}

func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.Background())
}

func (m *MongoDB) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, nil)
}

// GetDB returns the underlying database connection
func (m *MongoDB) GetDB() *mongo.Database {
	return m.db
}

// Collection returns a handle to a MongoDB collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.db.Collection(name)
}
