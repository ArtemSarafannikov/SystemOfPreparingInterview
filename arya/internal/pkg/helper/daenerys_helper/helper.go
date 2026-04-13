package daenerys_helper

import (
	"context"

	"github.com/CodefriendOrg/arya/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/arya/internal/pkg/auth"
	"github.com/CodefriendOrg/arya/internal/pkg/model"
	"github.com/CodefriendOrg/arya/internal/pkg/user_error"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service .
type Service struct {
	daenerysClient daenerys.DaenerysClient
}

// NewService .
func NewService(daenerysClient daenerys.DaenerysClient) *Service {
	return &Service{
		daenerysClient: daenerysClient,
	}
}

// SendSubmission .
func (s *Service) SendSubmission(ctx context.Context, problemID uuid.UUID, code string, language model.ProgrammingLanguage) (*daenerys.Submission, error) {
	user := auth.GetUserFromContext(ctx)

	resp, err := s.daenerysClient.SendSubmission(ctx, &daenerys.SendSubmissionRequest{
		ProblemId: problemID.String(),
		UserId:    user.ID,
		Code:      code,
		Language:  language.ToProto(),
	})
	switch status.Code(err) {
	case codes.OK:
		return resp, nil
	default:
		return nil, user_error.New(user_error.InternalError, "s.daenerysClient.SendSubmission: %v", err)
	}
}

// GetSubmissionByID .
func (s *Service) GetSubmissionByID(ctx context.Context, submissionID uuid.UUID) (*daenerys.Submission, error) {
	user := auth.GetUserFromContext(ctx)

	resp, err := s.daenerysClient.GetSubmission(ctx, &daenerys.GetSubmissionRequest{
		Id: submissionID.String(),
	})
	switch status.Code(err) {
	case codes.OK:
		if resp.UserId != user.ID {
			return nil, user_error.WithoutLoggerMessage(user_error.PermissionDenied)
		}
		return resp, nil
	case codes.NotFound:
		return nil, user_error.WithoutLoggerMessage(user_error.NotFound)
	default:
		return nil, user_error.New(user_error.InternalError, "s.daenerysClient.GetSubmission: %v", err)
	}
}

// GetSubmissionsByProblem .
func (s *Service) GetSubmissionsByProblem(ctx context.Context, problemID uuid.UUID, pagination model.Pagination) ([]*daenerys.Submission, uint64, error) {
	user := auth.GetUserFromContext(ctx)

	resp, err := s.daenerysClient.ListSubmissions(ctx, &daenerys.ListSubmissionsRequest{
		Filter: &daenerys.ListSubmissionsRequest_Filter{
			ProblemIdEq: lo.ToPtr(problemID.String()),
			UserIdEq:    &user.ID,
		},
		OrderBy:        daenerys.ListSubmissionsRequest_CREATED_AT,
		OrderDirection: daenerys.OrderDirection_DESC,
		Pagination: &daenerys.Pagination{
			Page:    uint64(pagination.Page),    //nolint:gosec
			PerPage: uint64(pagination.PerPage), //nolint:gosec
		},
	})
	if err != nil {
		return nil, 0, user_error.New(user_error.InternalError, "s.daenerysClient.ListSubmissions: %v", err)
	}

	return resp.Submissions, resp.TotalItems, nil
}
