package usecase

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
	"github.com/google/uuid"
)

// CreateProblemParams .
type CreateProblemParams struct {
	Summary       string
	Description   string
	AuthorID      uuid.UUID
	TimeLimitMs   int64
	MemoryLimitKb int64
}

// CreateProblem .
func (s *Service) CreateProblem(ctx context.Context, params CreateProblemParams) (*tirion.Problem, error) {
	createdProblem, err := s.storage.CreateProblem(ctx, &model.Problem{
		Summary:       params.Summary,
		Description:   params.Description,
		AuthorID:      params.AuthorID,
		TimeLimitMs:   params.TimeLimitMs,
		MemoryLimitKb: params.MemoryLimitKb,
	})
	if err != nil {
		return nil, fmt.Errorf("s.storage.CreateProblem: %w", err)
	}

	return createdProblem.ToProto(), nil
}
