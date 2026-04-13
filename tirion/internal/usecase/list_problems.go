package usecase

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
	"github.com/CodefriendOrg/tirion/internal/pkg/store"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// ListProblems .
func (s *Service) ListProblems(ctx context.Context, req *tirion.ListProblemsRequest) ([]*model.Problem, uint64, error) {
	params := store.ListProblemsParams{
		IDIn: lo.Map(req.Filter.IdIn, func(item string, _ int) uuid.UUID {
			return uuid.MustParse(item)
		}),
	}
	if req.Filter.AuthorIdEq != nil {
		params.AuthorIDEq = lo.ToPtr(uuid.MustParse(*req.Filter.AuthorIdEq))
	}
	if req.Filter.CreatedAtGte != nil {
		params.CreatedAtGte = lo.ToPtr(req.Filter.CreatedAtGte.AsTime())
	}
	if req.Filter.CreatedAtLte != nil {
		params.CreatedAtLte = lo.ToPtr(req.Filter.CreatedAtLte.AsTime())
	}

	params.Limit, params.Offset = store.GetLimitOffsetFromProtoPagination(req.Pagination)
	switch req.OrderBy {
	case tirion.ListProblemsRequest_CREATED_AT:
		params.OrderBy = "created_at"
	}
	params.OrderDirection = req.OrderDirection.String()

	problems, count, err := s.storage.ListProblems(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("storage.ListProblems: %w", err)
	}

	return problems, count, nil
}
