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
