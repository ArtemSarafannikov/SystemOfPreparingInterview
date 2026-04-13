package app

import (
	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/model"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *testSuite) Test_validateUpdateSubmissionStatusRequest() {
	testCases := []struct {
		name    string
		req     *daenerys.UpdateSubmissionStatusRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "request is required",
		},
		{
			name: "id empty",
			req: &daenerys.UpdateSubmissionStatusRequest{
				Status: daenerys.SubmissionStatus_STATUS_NEW,
			},
			wantErr: "id: cannot be blank.",
		},
		{
			name: "id is not uuid",
			req: &daenerys.UpdateSubmissionStatusRequest{
				Id:     "123",
				Status: daenerys.SubmissionStatus_STATUS_NEW,
			},
			wantErr: "id: must be a valid UUID.",
		},
		{
			name: "empty status",
			req: &daenerys.UpdateSubmissionStatusRequest{
				Id: uuid.NewString(),
			},
			wantErr: "status: cannot be blank.",
		},
		{
			name: "success",
			req: &daenerys.UpdateSubmissionStatusRequest{
				Id:     uuid.NewString(),
				Status: daenerys.SubmissionStatus_STATUS_NEW,
			},
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			err := validateUpdateSubmissionStatusRequest(tt.req)
			if tt.wantErr != "" {
				s.Require().EqualError(err, tt.wantErr)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *testSuite) Test_UpdateSubmissionStatus() {
	s.Run("success", func() {
		submission, err := s.storage.CreateSubmission(s.ctx, &model.Submission{
			ProblemID: uuid.New(),
			UserID:    uuid.New(),
			Code:      "test",
		})
		s.Require().NoError(err)
		s.Require().NotNil(submission)
		s.Require().Equal(daenerys.SubmissionStatus_STATUS_NEW, submission.Status)

		impl := s.newMockedImplementation()
		resp, err := impl.UpdateSubmissionStatus(s.ctx, &daenerys.UpdateSubmissionStatusRequest{
			Id:     submission.ID.String(),
			Status: daenerys.SubmissionStatus_STATUS_JUDGING,
		})
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().Equal(daenerys.SubmissionStatus_STATUS_JUDGING, resp.Status)
	})

	s.Run("invalid argument", func() {
		impl := s.newMockedImplementation()
		resp, err := impl.UpdateSubmissionStatus(s.ctx, nil)
		s.Require().Nil(resp)

		stErr, ok := status.FromError(err)
		s.Require().True(ok)
		s.Require().Equal(codes.InvalidArgument, stErr.Code())
		s.Require().Equal("request is required", stErr.Message())
	})

	s.Run("not found", func() {
		impl := s.newMockedImplementation()
		resp, err := impl.UpdateSubmissionStatus(s.ctx, &daenerys.UpdateSubmissionStatusRequest{
			Id:     uuid.NewString(),
			Status: daenerys.SubmissionStatus_STATUS_JUDGING,
		})
		s.Require().Nil(resp)

		stErr, ok := status.FromError(err)
		s.Require().True(ok)
		s.Require().Equal(codes.NotFound, stErr.Code())
		s.Require().Equal("no rows in result set", stErr.Message())
	})
}
