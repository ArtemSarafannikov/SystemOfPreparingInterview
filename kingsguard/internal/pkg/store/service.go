package store

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/kingsguard/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DatabaseConfig) (*Storage, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnString())
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	return &Storage{
		pool: pool,
	}, nil
}
