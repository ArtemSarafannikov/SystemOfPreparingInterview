package usecase

import (
	"github.com/google/uuid"
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/kingsguard/pkg/kingsguard"
)

func (s *testSuite) Test_CreateProblem() {
	s.Run("success", func() {
		userID := uuid.New()
		user := &kingsguard.User{
			Id: userID.String(),
		}

		service := s.newMockedService()

		problem, err := service.CreateProblem(s.ctx, CreateProblemParams{
			Summary:       "test",
			Description:   "test",
			AuthorID:      userID,
			TimeLimitMs:   5000,
			MemoryLimitKb: 65536,
		})
		s.Require().NoError(err)
		s.Require().NotNil(problem)
		s.Require().Equal("test", problem.Summary)
		s.Require().Equal("test", problem.Description)
		s.Require().Equal(user.Id, problem.AuthorId)
	})
}
