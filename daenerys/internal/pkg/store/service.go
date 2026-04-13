package store

import (
	"context"
	"fmt"
	"math"

	"github.com/CodefriendOrg/daenerys/internal/config"
	desc "github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
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
func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

// Builder .
func Builder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

// NewTest .
func NewTest(ctx context.Context) (*Storage, error) {
	cfg := config.DatabaseConfig{
		Host:               "localhost",
		Port:               5432,
		User:               "postgres",
		Password:           "postgres",
		Name:               "daenerys_test",
		SslMode:            "disable",
		MaxOpenConnections: 10,
	}
	pool, err := cfg.GetConn(ctx)
	return &Storage{pool: pool}, err
}

// GetLimitOffsetFromProtoPagination .
func GetLimitOffsetFromProtoPagination(pagination *desc.Pagination) (limit uint64, offset uint64) {
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

// WithTransaction .
func (s *Storage) WithTransaction(ctx context.Context, f func() error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if err = f(); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
