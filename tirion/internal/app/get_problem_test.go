package app

import (
	"github.com/google/uuid"
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
)

func (s *testSuite) Test_GetProblem() {
	s.Run("success", func() {
		impl := s.newMockedImplementation()

		problem, err := s.storage.CreateProblem(s.ctx, &model.Problem{
			Summary:       uuid.NewString(),
			TimeLimitMs:   1000,
			MemoryLimitKb: 65536,
		})
		s.Require().NoError(err)

		resp, err := impl.GetProblem(s.ctx, &tirion.GetProblemRequest{
			Id:        problem.ID.String(),
			WithTests: true,
		})
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().Equal(problem.ToProto(), resp.Problem)
		s.Require().NotNil(resp.Tests)
	})
}
