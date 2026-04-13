package store

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
)

// ListTestsParams .
type ListTestsParams struct {
	ProblemIDEq *uuid.UUID
	WithHidden  bool

	OrderBy        string
	OrderDirection string
	Limit          uint64
	Offset         uint64
}

// ListTests .
func (s *Storage) ListTests(ctx context.Context, params ListTestsParams) ([]*model.Test, uint64, error) {
	query := Builder().
		Select(model.TestSchema.Columns()...).
		From(model.TestSchema.TableName())

	if params.ProblemIDEq != nil {
		query = query.Where(sq.Eq{"problem_id": params.ProblemIDEq})
	}
	if !params.WithHidden {
		query = query.Where(sq.Eq{"is_hidden": false})
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

	tests, err := Selectx[model.Test](ctx, s.pool, query)
	if err != nil {
		return nil, 0, err
	}
	return tests, count, nil
}
