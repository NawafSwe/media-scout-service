package db

import (
	"fmt"

	"github.com/NawafSwe/media-scout-service/cmd/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// NewDBConn creates a new db connection and returns *sqlx.DB.
func NewDBConn(cfg config.DB) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create a db conn: %w", err)
	}
	if cfg.MaxOpenConnections != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	}
	if cfg.MaxIdleConnections != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	}
	if cfg.MaxConnectionsLifetime != 0 {
		db.SetConnMaxLifetime(cfg.MaxConnectionsLifetime)
	}
	return db, nil
}
