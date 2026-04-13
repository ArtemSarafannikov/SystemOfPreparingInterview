package usecase

import (
	"github.com/CodefriendOrg/daenerys/internal/pkg/store"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// Service .
type Service struct {
	storage     *store.Storage
	riverClient *river.Client[pgx.Tx]
}

// NewService .
func NewService(storage *store.Storage, riverClient *river.Client[pgx.Tx]) *Service {
	return &Service{
		storage:     storage,
		riverClient: riverClient,
	}
}
