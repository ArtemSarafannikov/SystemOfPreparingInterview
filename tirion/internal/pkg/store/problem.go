package store

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
)

// GetProblemByID .
func (s *Storage) GetProblemByID(ctx context.Context, id uuid.UUID) (*model.Problem, error) {
	query := Builder().
		Select(model.ProblemSchema.Columns()...).
		From(model.ProblemSchema.TableName()).
		Where(sq.Eq{"id": id.String()})

	return Getx[model.Problem](ctx, s.pool, query)
}

// CreateProblem .
func (s *Storage) CreateProblem(ctx context.Context, problem *model.Problem) (*model.Problem, error) {
	query := Builder().
		Insert(model.ProblemSchema.TableName()).
		SetMap(map[string]any{
			"summary":           problem.Summary,
			"description":       problem.Description,
			"author_id":         problem.AuthorID,
			"time_limit_ms":     problem.TimeLimitMs,
			"memory_limit_kb":   problem.MemoryLimitKb,
		}).
		Suffix(model.ProblemSchema.Returning())

	return Getx[model.Problem](ctx, s.pool, query)
}

// ListProblemsParams .
type ListProblemsParams struct {
	IDIn         []uuid.UUID
	AuthorIDEq   *uuid.UUID
	CreatedAtLte *time.Time
	CreatedAtGte *time.Time

	OrderBy        string
	OrderDirection string
	Limit          uint64
	Offset         uint64
}

// ListProblems .
func (s *Storage) ListProblems(ctx context.Context, params ListProblemsParams) ([]*model.Problem, uint64, error) {
	query := Builder().
		Select(model.ProblemSchema.Columns()...).
		From(model.ProblemSchema.TableName())

	if len(params.IDIn) > 0 {
		query = query.Where(sq.Eq{"id": params.IDIn})
	}
	if params.AuthorIDEq != nil {
		query = query.Where(sq.Eq{"author_id": params.AuthorIDEq})
	}
	if params.CreatedAtLte != nil {
		query = query.Where(sq.LtOrEq{"created_at": params.CreatedAtLte})
	}
	if params.CreatedAtGte != nil {
		query = query.Where(sq.GtOrEq{"created_at": params.CreatedAtGte})
	}
	count, err := Count(ctx, s.pool, query)
	if err != nil {
		return nil, 0, err
	}

	if params.OrderBy != "" {
		query = query.OrderBy(fmt.Sprintf("%s %s", params.OrderBy, params.OrderDirection))
	}
	if params.Limit != 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}

	problems, err := Selectx[model.Problem](ctx, s.pool, query)
	if err != nil {
		return nil, 0, err
	}

	return problems, count, nil
}
