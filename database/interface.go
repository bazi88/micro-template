package database

import (
	"context"
)

// DataStore định nghĩa interface chung cho tất cả các loại database
type DataStore interface {
	Connect() error
	Close() error
	Ping(ctx context.Context) error
}

// Repository định nghĩa các phương thức CRUD cơ bản
type Repository interface {
	Create(ctx context.Context, collection string, data interface{}) error
	FindOne(ctx context.Context, collection string, filter interface{}) (interface{}, error)
	FindMany(ctx context.Context, collection string, filter interface{}) ([]interface{}, error)
	Update(ctx context.Context, collection string, filter interface{}, update interface{}) error
	Delete(ctx context.Context, collection string, filter interface{}) error
}
