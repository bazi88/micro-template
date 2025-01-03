package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Api
	Cors
	Services *ServicesConfig
	Database
	Cache
	Elasticsearch

	OpenTelemetry
	Session
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	return &Config{
		Api:           API(),
		Cors:          NewCors(),
		Services:      NewServicesConfig(),
		Database:      DataStore(),
		Cache:         NewCache(),
		Elasticsearch: ElasticSearch(),
		Session:       NewSession(),
		OpenTelemetry: NewOpenTelemetry(),
	}
}
