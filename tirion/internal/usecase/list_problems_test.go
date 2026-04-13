package usecase

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
)

func (s *testSuite) Test_ListProblems() {
	problem1, errCreate := s.storage.CreateProblem(s.ctx, &model.Problem{
		Summary:       "test summary",
		Description:   "test description",
		TimeLimitMs:   1000,
		MemoryLimitKb: 65536,
	})
	s.Require().NoError(errCreate)

	problem2, errCreate := s.storage.CreateProblem(s.ctx, &model.Problem{
		Summary:       "test summary 2",
		Description:   "test description",
		AuthorID:      uuid.New(),
		TimeLimitMs:   1000,
		MemoryLimitKb: 65536,
	})
	s.Require().NoError(errCreate)

	service := s.newMockedService()

	s.Run("filter by id", func() {
		problems, count, err := service.ListProblems(s.ctx, &tirion.ListProblemsRequest{
			Filter: &tirion.ListProblemsRequest_Filter{
				IdIn: []string{
					problem1.ID.String(),
				},
			},
		})
		s.Require().NoError(err)
		s.Require().Len(problems, 1)
		s.Require().Equal(uint64(1), count)
		s.Require().Equal(problems[0], problem1)
	})

	s.Run("filter by author_id", func() {
		problems, count, err := service.ListProblems(s.ctx, &tirion.ListProblemsRequest{
			Filter: &tirion.ListProblemsRequest_Filter{
				AuthorIdEq: lo.ToPtr(problem2.AuthorID.String()),
			},
		})
		s.Require().NoError(err)
		s.Require().Len(problems, 1)
		s.Require().Equal(uint64(1), count)
		s.Require().Equal(problems[0], problem2)
	})
}
