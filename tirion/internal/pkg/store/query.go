package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Getx executes a query built by squirrel and scans the first row into a struct T using db tags.
// Returns pgx.ErrNoRows if no rows are found.
func Getx[T any](ctx context.Context, pool *pgxpool.Pool, query squirrel.Sqlizer) (*T, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[T])
}

// Selectx executes a query built by squirrel and scans all rows into a slice of *T using db tags.
// Returns an empty slice if no rows are found.
func Selectx[T any](ctx context.Context, pool *pgxpool.Pool, query squirrel.Sqlizer) ([]*T, error) {
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[T])
}
