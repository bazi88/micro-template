package main

import (
	"micro/internal/server"
)

// Version is injected using ldflags during build time
var Version = "v0.1.0"

// @title Go8
// @version 0.1.0
// @description Go + Postgres + Chi router + sqlx + ent + Testing starter kit for API development.
// @contact.name User Name
// @contact.url https://micro
// @contact.email email@example.com
// @host localhost:3080
// @BasePath /
func main() {
	s := server.New(server.WithVersion(Version))
	s.Init()
	s.Run()
}
