package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a bounded pgxpool tuned for the 3.6 GB ARM64 container.
// Note: SMtracker owns the PostgreSQL schema; this server never runs DDL migrations.
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	cfg.MinConns = 1
	cfg.MaxConns = 8
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
