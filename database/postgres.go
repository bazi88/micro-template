package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db     *sqlx.DB
	config *Config
}

func NewPostgresDB(config *Config) *PostgresDB {
	return &PostgresDB{
		config: config,
	}
}

func (p *PostgresDB) Connect() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.config.Host,
		p.config.Port,
		p.config.User,
		p.config.Password,
		p.config.Database,
		p.config.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to postgres: %v", err)
	}

	p.db = db
	return nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// GetDB returns the underlying database connection
func (p *PostgresDB) GetDB() *sqlx.DB {
	return p.db
}
