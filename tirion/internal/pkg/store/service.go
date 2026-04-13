package store

import (
	"context"
	"fmt"
	"math"

	"github.com/CodefriendOrg/tirion/internal/config"
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
	countLimit    = 10000
)

// Storage .
type Storage struct {
	pool *pgxpool.Pool
}

// New .
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

// Close closes the database connection pool.
func (s *Storage) Close() {
	s.pool.Close()
}

// Builder .
func Builder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

// GetLimitOffsetFromProtoPagination .
func GetLimitOffsetFromProtoPagination(pagination *tirion.Pagination) (limit uint64, offset uint64) {
	if pagination == nil {
		return defaultLimit, defaultOffset
	}

	limit = pagination.PerPage
	offset = pagination.PerPage * (pagination.Page - 1)

	return
}

// Count use before pagination, limit, offset
func Count(ctx context.Context, pool *pgxpool.Pool, queryBuilder squirrel.SelectBuilder) (uint64, error) {
	var count uint64
	query, args, _ := prepareCountQuery(queryBuilder).ToSql()
	if err := pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}

	if count == uint64(countLimit) {
		count = math.MaxInt32
	}

	return count, nil
}

func prepareCountQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	subQuery := query.
		RemoveColumns().Column("1").
		RemoveOffset().Limit(uint64(countLimit))

	return Builder().
		Select("COUNT(*)").
		FromSelect(subQuery, "rows")
}

// NewTest .
func NewTest(ctx context.Context) (*Storage, error) {
	return New(ctx, config.DatabaseConfig{
		Host:               "localhost",
		Port:               5432,
		User:               "postgres",
		Password:           "postgres",
		Name:               "tirion_test",
		SslMode:            "disable",
		MaxOpenConnections: 10,
	})
}
