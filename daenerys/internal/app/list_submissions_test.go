package app

import (
	desc "github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/model"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *testSuite) Test_ListSubmissions() {
	s.Run("validation error", func() {
		impl := s.newMockedImplementation()

		resp, err := impl.ListSubmissions(s.ctx, nil)
		s.Require().Nil(resp)

		errSt, ok := status.FromError(err)
		s.Require().True(ok)
		s.Require().Equal(codes.InvalidArgument, errSt.Code())
		s.Require().Equal("request is required", errSt.Message())
	})

	s.Run("success by id", func() {
		impl := s.newMockedImplementation()

		submission, err := s.storage.CreateSubmission(s.ctx, &model.Submission{
			ProblemID: uuid.New(),
			UserID:    uuid.New(),
			Code:      "123",
			Language:  desc.ProgrammingLanguage_LANGUAGE_UNKNOWN,
		})
		s.Require().NoError(err)

		resp, err := impl.ListSubmissions(s.ctx, &desc.ListSubmissionsRequest{
			Filter: &desc.ListSubmissionsRequest_Filter{
				SubmissionIdIn: []string{submission.ID.String()},
			},
			Pagination: &desc.Pagination{
				Page:    1,
				PerPage: 10,
			},
		})
		s.Require().NoError(err)
		s.Require().Len(resp.Submissions, 1)
		s.Require().Equal(uint64(1), resp.TotalItems)
		s.Require().Equal(submission.ToProto(), resp.Submissions[0])
	})

	s.Run("success by problem_id", func() {
		impl := s.newMockedImplementation()

		submission, err := s.storage.CreateSubmission(s.ctx, &model.Submission{
			ProblemID: uuid.New(),
			UserID:    uuid.New(),
			Code:      "123",
			Language:  desc.ProgrammingLanguage_LANGUAGE_UNKNOWN,
		})
		s.Require().NoError(err)

		resp, err := impl.ListSubmissions(s.ctx, &desc.ListSubmissionsRequest{
			Filter: &desc.ListSubmissionsRequest_Filter{
				ProblemIdEq: lo.ToPtr(submission.ProblemID.String()),
			},
			Pagination: &desc.Pagination{
				Page:    1,
				PerPage: 10,
			},
		})
		s.Require().NoError(err)
		s.Require().Len(resp.Submissions, 1)
		s.Require().Equal(uint64(1), resp.TotalItems)
		s.Require().Equal(submission.ToProto(), resp.Submissions[0])
	})

	s.Run("success by user_id", func() {
		impl := s.newMockedImplementation()

		submission, err := s.storage.CreateSubmission(s.ctx, &model.Submission{
			ProblemID: uuid.New(),
			UserID:    uuid.New(),
			Code:      "123",
			Language:  desc.ProgrammingLanguage_LANGUAGE_UNKNOWN,
		})
		s.Require().NoError(err)

		resp, err := impl.ListSubmissions(s.ctx, &desc.ListSubmissionsRequest{
			Filter: &desc.ListSubmissionsRequest_Filter{
				UserIdEq: lo.ToPtr(submission.UserID.String()),
			},
			Pagination: &desc.Pagination{
				Page:    1,
				PerPage: 10,
			},
		})
		s.Require().NoError(err)
		s.Require().Len(resp.Submissions, 1)
		s.Require().Equal(uint64(1), resp.TotalItems)
		s.Require().Equal(submission.ToProto(), resp.Submissions[0])
	})
}

func (s *testSuite) Test_validateListSubmissionsRequest() {
	testCases := []struct {
		name      string
		req       *desc.ListSubmissionsRequest
		wantedErr string
	}{
		{
			name:      "nil request",
			req:       nil,
			wantedErr: "request is required",
		},
		{
			name: "nil filter",
			req: &desc.ListSubmissionsRequest{
				Filter: nil,
				Pagination: &desc.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "filter: cannot be blank.",
		},
		{
			name: "filter is empty",
			req: &desc.ListSubmissionsRequest{
				Filter: &desc.ListSubmissionsRequest_Filter{},
				Pagination: &desc.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "filter: must have at least one field",
		},
		{
			name: "submission_id_in has empty value",
			req: &desc.ListSubmissionsRequest{
				Filter: &desc.ListSubmissionsRequest_Filter{
					SubmissionIdIn: []string{""},
				},
				Pagination: &desc.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "submission_id_in: (0: cannot be blank.).",
		},
		{
			name: "submission_id_in has not uuid value",
			req: &desc.ListSubmissionsRequest{
				Filter: &desc.ListSubmissionsRequest_Filter{
					SubmissionIdIn: []string{"1"},
				},
				Pagination: &desc.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "submission_id_in: (0: must be a valid UUID.).",
		},
		{
			name: "problem_id_eq is empty",
			req: &desc.ListSubmissionsRequest{
				Filter: &desc.ListSubmissionsRequest_Filter{
					ProblemIdEq: lo.ToPtr(""),
				},
				Pagination: &desc.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "problem_id_eq: cannot be blank.",
		},
		{
			name: "problem_id_eq is not uuid",
			req: &desc.ListSubmissionsRequest{
				Filter: &desc.ListSubmissionsRequest_Filter{
					ProblemIdEq: lo.ToPtr("1"),
				},
				Pagination: &desc.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "problem_id_eq: must be a valid UUID.",
		},
		{
			name: "user_id_eq is empty",
			req: &desc.ListSubmissionsRequest{
				Filter: &desc.ListSubmissionsRequest_Filter{
					UserIdEq: lo.ToPtr(""),
				},
				Pagination: &desc.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "user_id_eq: cannot be blank.",
		},
		{
			name: "user_id_eq is not uuid",
			req: &desc.ListSubmissionsRequest{
				Filter: &desc.ListSubmissionsRequest_Filter{
					UserIdEq: lo.ToPtr("1"),
				},
				Pagination: &desc.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "user_id_eq: must be a valid UUID.",
		},
		{
			name: "pagination is nil",
			req: &desc.ListSubmissionsRequest{
				Filter:     &desc.ListSubmissionsRequest_Filter{},
				Pagination: nil,
			},
			wantedErr: "pagination: cannot be blank.",
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			err := validateListSubmissionsRequest(tt.req)
			s.EqualError(err, tt.wantedErr)
		})
	}
}
