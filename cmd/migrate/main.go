package main

import (
	"log"

	"micro/config"
	"micro/database"
	db "micro/third_party/database"
)

// Version is injected using ldflags during build time
var Version string

func main() {
	log.Printf("Version: %s\n", Version)

	cfg := config.New()
	store := db.NewSqlx(cfg.Database)
	migrator := database.Migrator(store.DB)

	// todo: accept cli flag for other operations
	migrator.Up()
}
