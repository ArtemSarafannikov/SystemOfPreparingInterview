package app

import (
	"github.com/google/uuid"
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *testSuite) Test_validateListTestsRequest() {
	testCases := []struct {
		name      string
		req       *tirion.ListTestsRequest
		wantedErr string
	}{
		{
			name:      "nil request",
			req:       nil,
			wantedErr: "request is required",
		},
		{
			name: "nil filter",
			req: &tirion.ListTestsRequest{
				Filter: nil,
				Pagination: &tirion.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "filter: cannot be blank.",
		},
		{
			name: "problem_id_eq is empty",
			req: &tirion.ListTestsRequest{
				Filter: &tirion.ListTestsRequest_Filter{},
				Pagination: &tirion.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "problem_id_eq: cannot be blank.",
		},
		{
			name: "problem_id_eq is not uuid",
			req: &tirion.ListTestsRequest{
				Filter: &tirion.ListTestsRequest_Filter{
					ProblemIdEq: "123",
				},
				Pagination: &tirion.Pagination{
					Page:    1,
					PerPage: 10,
				},
			},
			wantedErr: "problem_id_eq: must be a valid UUID.",
		},
		{
			name: "pagination is nil",
			req: &tirion.ListTestsRequest{
				Filter: &tirion.ListTestsRequest_Filter{
					ProblemIdEq: uuid.NewString(),
				},
				Pagination: nil,
			},
			wantedErr: "pagination: cannot be blank.",
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			err := validateListTestsRequest(tt.req)
			s.EqualError(err, tt.wantedErr)
		})
	}
}

func (s *testSuite) Test_ListTests() {
	s.Run("validation error", func() {
		impl := s.newMockedImplementation()

		resp, err := impl.ListTests(s.ctx, nil)
		s.Require().Nil(resp)

		errSt, ok := status.FromError(err)
		s.Require().True(ok)
		s.Require().Equal(codes.InvalidArgument, errSt.Code())
		s.Require().Equal("request is required", errSt.Message())
	})

	s.Run("success", func() {
		impl := s.newMockedImplementation()

		resp, err := impl.ListTests(s.ctx, &tirion.ListTestsRequest{
			Filter: &tirion.ListTestsRequest_Filter{
				ProblemIdEq: uuid.NewString(),
			},
			Pagination: &tirion.Pagination{
				Page:    1,
				PerPage: 10,
			},
		})
		s.Require().NoError(err)
		s.Require().Len(resp.Tests, 0)
		s.Require().Equal(uint64(0), resp.TotalItems)
	})
}
