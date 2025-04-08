package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreSQL struct {
	Pool *pgxpool.Pool
}

func NewPostgreSQL(ctx context.Context, dsn string) (*PostgreSQL, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе URL подключения к PostgreSQL: %w", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 5 * time.Minute
	cfg.MaxConnIdleTime = 1 * time.Minute

	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к PostgreSQL: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ошибка при проверке подключения к PostgreSQL: %w", err)
	}

	return &PostgreSQL{Pool: pool}, nil
}

func (p *PostgreSQL) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
