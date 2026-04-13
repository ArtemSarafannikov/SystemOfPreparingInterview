package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/CodefriendOrg/kingsguard/internal/model"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) CreateUser(ctx context.Context, username, password string) (*model.User, error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(model.User{}.TableName()).
		SetMap(map[string]any{
			"username": username,
			"password": password,
		}).
		Suffix("RETURNING *")

	return Getx[model.User](ctx, s.pool, query)
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From(model.User{}.TableName()).
		Where(sq.Eq{"username": username})

	return Getx[model.User](ctx, s.pool, query)
}

func (s *Storage) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("*").
		From(model.User{}.TableName()).
		Where(sq.Eq{"id": userID})

	return Getx[model.User](ctx, s.pool, query)
}
