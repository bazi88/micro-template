package database

import "fmt"

// NewDatabase tạo instance database dựa trên cấu hình
func NewDatabase(config *Config) (DataStore, error) {
	switch config.Type {
	case "postgres":
		db := NewPostgresDB(config)
		if err := db.Connect(); err != nil {
			return nil, err
		}
		return db, nil
	case "mongodb":
		db := NewMongoDB(config)
		if err := db.Connect(); err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}
