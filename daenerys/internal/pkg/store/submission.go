package store

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// CreateSubmission .
func (s *Storage) CreateSubmission(ctx context.Context, submission *model.Submission) (*model.Submission, error) {
	query := Builder().
		Insert(model.SubmissionSchema.TableName()).
		SetMap(map[string]any{
			"problem_id": submission.ProblemID,
			"user_id":    submission.UserID,
			"status":     daenerys.SubmissionStatus_STATUS_NEW,
			"code":       submission.Code,
			"language":   submission.Language,
		}).
		Suffix(model.SubmissionSchema.Returning())

	return Getx[model.Submission](ctx, s.pool, query)
}

// UpdateSubmissionStatus .
func (s *Storage) UpdateSubmissionStatus(ctx context.Context, submissionID uuid.UUID, status daenerys.SubmissionStatus) (*model.Submission, error) {
	query := Builder().
		Update(model.SubmissionSchema.TableName()).
		SetMap(map[string]any{
			"status": status,
		}).
		Where(sq.Eq{"id": submissionID}).
		Suffix(model.SubmissionSchema.Returning())

	return Getx[model.Submission](ctx, s.pool, query)
}

// GetSubmissionByID .
func (s *Storage) GetSubmissionByID(ctx context.Context, submissionID uuid.UUID) (*model.Submission, error) {
	query := Builder().
		Select(model.SubmissionSchema.Columns()...).
		From(model.SubmissionSchema.TableName()).
		Where(sq.Eq{"id": submissionID})

	return Getx[model.Submission](ctx, s.pool, query)
}

// ListSubmissionsParams .
type ListSubmissionsParams struct {
	IDIn        []uuid.UUID
	UserIDEq    *uuid.UUID
	ProblemIDEq *uuid.UUID

	OrderBy        string
	OrderDirection string
	Limit          uint64
	Offset         uint64
}

// ListSubmissions .
func (s *Storage) ListSubmissions(ctx context.Context, params ListSubmissionsParams) ([]*model.Submission, uint64, error) {
	query := Builder().
		Select(model.SubmissionSchema.Columns()...).
		From(model.SubmissionSchema.TableName()).
		Limit(params.Limit).
		Offset(params.Offset)

	if len(params.IDIn) > 0 {
		query = query.Where(sq.Eq{"id": params.IDIn})
	}
	if params.UserIDEq != nil {
		query = query.Where(sq.Eq{"user_id": params.UserIDEq})
	}
	if params.ProblemIDEq != nil {
		query = query.Where(sq.Eq{"problem_id": params.ProblemIDEq})
	}

	count, err := Count(ctx, s.pool, query)
	if err != nil {
		return nil, 0, err
	}

	if params.OrderBy != "" {
		query = query.OrderBy(fmt.Sprintf("%s %s", params.OrderBy, params.OrderDirection))
	}

	submissions, err := Selectx[model.Submission](ctx, s.pool, query)
	if err != nil {
		return nil, 0, err
	}
	return submissions, count, nil
}
