package app

import (
	desc "github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/model"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *testSuite) Test_GetSubmission() {
	s.Run("validation error", func() {
		api := s.newMockedImplementation()
		resp, err := api.GetSubmission(s.ctx, nil)
		s.Require().Nil(resp)
		s.Require().Error(err)
		s.Require().Equal(status.Code(err), codes.InvalidArgument)
	})

	s.Run("success", func() {
		submission, err := s.storage.CreateSubmission(s.ctx, &model.Submission{
			ProblemID: uuid.New(),
			UserID:    uuid.New(),
			Code:      "123",
			Language:  desc.ProgrammingLanguage_LANGUAGE_UNKNOWN,
		})
		s.Require().NoError(err)

		api := s.newMockedImplementation()
		resp, err := api.GetSubmission(s.ctx, &desc.GetSubmissionRequest{
			Id: submission.ID.String(),
		})
		s.Require().NotNil(resp)
		s.Require().NoError(err)
		s.Require().Equal(submission.ToProto(), resp)
	})

	s.Run("not found", func() {
		api := s.newMockedImplementation()
		resp, err := api.GetSubmission(s.ctx, &desc.GetSubmissionRequest{
			Id: uuid.NewString(),
		})
		s.Require().Nil(resp)
		s.Require().Error(err)
		s.Require().Equal(status.Code(err), codes.NotFound)
	})
}

func (s *testSuite) Test_validateGetSubmissionRequest() {
	testCases := []struct {
		name    string
		req     *desc.GetSubmissionRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "request is required",
		},
		{
			name:    "id is empty",
			req:     &desc.GetSubmissionRequest{},
			wantErr: "id: cannot be blank.",
		},
		{
			name: "id is not uuid",
			req: &desc.GetSubmissionRequest{
				Id: "1",
			},
			wantErr: "id: must be a valid UUID.",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := validateGetSubmissionRequest(tc.req)
			s.Require().EqualError(err, tc.wantErr)
		})
	}
}
