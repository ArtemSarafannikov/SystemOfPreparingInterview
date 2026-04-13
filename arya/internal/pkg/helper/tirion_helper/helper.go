package tirion_helper

import (
	"context"

	"github.com/CodefriendOrg/arya/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/arya/internal/pkg/model"
	"github.com/CodefriendOrg/arya/internal/pkg/user_error"
	"github.com/CodefriendOrg/arya/internal/pkg/utils"
	"github.com/google/uuid"
)

// Service .
type Service struct {
	tirionClient tirion.TirionClient
}

// NewService .
func NewService(tirionClient tirion.TirionClient) *Service {
	return &Service{
		tirionClient: tirionClient,
	}
}

// ListProblems .
func (s *Service) ListProblems(ctx context.Context, filter *model.ListProblemsFilter, pagination model.Pagination) ([]*tirion.Problem, uint64, error) {
	if filter == nil {
		filter = &model.ListProblemsFilter{}
	}
	resp, err := s.tirionClient.ListProblems(ctx, &tirion.ListProblemsRequest{
		Filter: &tirion.ListProblemsRequest_Filter{
			IdIn: utils.UUIDsToStrings(filter.Ids),
		},
		OrderBy:        tirion.ListProblemsRequest_CREATED_AT,
		OrderDirection: tirion.OrderDirection_ASC,
		Pagination: &tirion.Pagination{
			Page:    uint64(pagination.Page),    //nolint:gosec
			PerPage: uint64(pagination.PerPage), //nolint:gosec
		},
	})
	if err != nil {
		return nil, 0, user_error.New(user_error.InternalError, "s.tirionClient.ListProblems: %v", err)
	}

	return resp.Problems, resp.TotalItems, nil
}

// GetTestsByProblem .
func (s *Service) GetTestsByProblem(ctx context.Context, problemID uuid.UUID) ([]*tirion.Test, error) {
	page := uint64(1)
	perPage := uint64(1000)
	tests := make([]*tirion.Test, 0)

	for {
		resp, err := s.tirionClient.ListTests(ctx, &tirion.ListTestsRequest{
			Filter: &tirion.ListTestsRequest_Filter{
				ProblemIdEq: problemID.String(),
			},
			Pagination: &tirion.Pagination{
				Page:    page,
				PerPage: perPage,
			},
		})
		if err != nil {
			return nil, user_error.FromGRPCError(err)
		}
		tests = append(tests, resp.Tests...)

		if len(resp.Tests) < int(perPage) {
			break
		}

		page++
	}

	return tests, nil
}

// GetProblemByID .
func (s *Service) GetProblemByID(ctx context.Context, problemID uuid.UUID) (*tirion.Problem, error) {
	resp, err := s.tirionClient.GetProblem(ctx, &tirion.GetProblemRequest{
		Id:        problemID.String(),
		WithTests: false,
	})
	if err != nil {
		return nil, user_error.FromGRPCError(err)
	}

	return resp.Problem, nil
}
