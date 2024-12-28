package database

import (
	"os"
	"strconv"
)

type Config struct {
	Type     string
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string

	// MongoDB specific
	MongoURI string

	// Redis specific
	RedisEnabled bool
	RedisHost    string
	RedisPort    int
	RedisPass    string
	RedisDB      int
}

func NewConfig() *Config {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisEnabled, _ := strconv.ParseBool(os.Getenv("REDIS_ENABLED"))

	return &Config{
		Type:     os.Getenv("DB_TYPE"),
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),

		MongoURI: os.Getenv("MONGO_URI"),

		RedisEnabled: redisEnabled,
		RedisHost:    os.Getenv("REDIS_HOST"),
		RedisPort:    redisPort,
		RedisPass:    os.Getenv("REDIS_PASSWORD"),
		RedisDB:      redisDB,
	}
}
