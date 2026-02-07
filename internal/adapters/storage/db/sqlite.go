package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

    // Import the generated SQLC code
	"waya/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/001_initial_schema.sql
var schema string

type Database struct {
	Conn *sql.DB
	Q    *Queries
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	var dsn string
	if cfg.Driver == "sqlite3" {
		// Enable Write-Ahead Logging (WAL) for concurrency speed
		dsn = fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000", cfg.Source)
	} else {
		dsn = cfg.Source
	}

	conn, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection Pool Settings (Important for high volume)
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("✅ Database connected", "driver", cfg.Driver)

	// AUTO-MIGRATE (Run schema on startup)
	if _, err := conn.Exec(schema); err != nil {
		// In production, use a real migration tool. For hackathon, this is fine.
		slog.Warn("⚠️ Migration warning (tables might exist):", "err", err)
	} else {
        slog.Info("✅ Schema applied successfully")
    }

	return &Database{
		Conn: conn,
		Q:    New(conn),
	}, nil
}