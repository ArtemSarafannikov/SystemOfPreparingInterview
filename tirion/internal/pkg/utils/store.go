package utils

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// ErrorIsNotFound .
func ErrorIsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, pgx.ErrNoRows)
}
