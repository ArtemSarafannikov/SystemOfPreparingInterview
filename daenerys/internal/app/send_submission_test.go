package app

import (
	desc "github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *testSuite) Test_SendSubmission() {
	s.Run("validation error", func() {
		api := s.newMockedImplementation()
		resp, err := api.SendSubmission(s.ctx, nil)
		s.Require().Nil(resp)
		s.Require().Error(err)
		s.Require().Equal(status.Code(err), codes.InvalidArgument)
	})
}

func (s *testSuite) Test_validateSendSubmissionRequest() {
	testCases := []struct {
		name    string
		req     *desc.SendSubmissionRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "request is required",
		},
		{
			name: "problem_id is empty",
			req: &desc.SendSubmissionRequest{
				UserId:   uuid.NewString(),
				Code:     "123",
				Language: desc.ProgrammingLanguage_LANGUAGE_PYTHON_3_12,
			},
			wantErr: "problem_id: cannot be blank.",
		},
		{
			name: "problem_id is not uuid",
			req: &desc.SendSubmissionRequest{
				ProblemId: "1",
				UserId:    uuid.NewString(),
				Code:      "123",
				Language:  desc.ProgrammingLanguage_LANGUAGE_PYTHON_3_12,
			},
			wantErr: "problem_id: must be a valid UUID.",
		},
		{
			name: "user_id is empty",
			req: &desc.SendSubmissionRequest{
				ProblemId: uuid.NewString(),
				Code:      "123",
				Language:  desc.ProgrammingLanguage_LANGUAGE_PYTHON_3_12,
			},
			wantErr: "user_id: cannot be blank.",
		},
		{
			name: "user_id is not uuid",
			req: &desc.SendSubmissionRequest{
				ProblemId: uuid.NewString(),
				UserId:    "1",
				Code:      "123",
				Language:  desc.ProgrammingLanguage_LANGUAGE_PYTHON_3_12,
			},
			wantErr: "user_id: must be a valid UUID.",
		},
		{
			name: "code is empty",
			req: &desc.SendSubmissionRequest{
				ProblemId: uuid.NewString(),
				UserId:    uuid.NewString(),
				Language:  desc.ProgrammingLanguage_LANGUAGE_PYTHON_3_12,
			},
			wantErr: "code: cannot be blank.",
		},
		{
			name: "language is empty",
			req: &desc.SendSubmissionRequest{
				ProblemId: uuid.NewString(),
				UserId:    uuid.NewString(),
				Code:      "1",
			},
			wantErr: "language: cannot be blank.",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := validateSendSubmissionRequest(tc.req)
			s.Require().EqualError(err, tc.wantErr)
		})
	}
}
