package usecase

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/logger"
	"github.com/CodefriendOrg/daenerys/internal/pkg/model"
	"github.com/CodefriendOrg/daenerys/internal/pkg/workers/python"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// SendSubmissionParams .
type SendSubmissionParams struct {
	ProblemID uuid.UUID
	UserID    uuid.UUID
	Code      string
	Language  daenerys.ProgrammingLanguage
}

// SendSubmission .
func (s *Service) SendSubmission(ctx context.Context, params SendSubmissionParams) (*model.Submission, error) {
	var submission *model.Submission
	err := s.storage.WithTransaction(ctx, func() error {
		var errTx error

		submission, errTx = s.storage.CreateSubmission(ctx, &model.Submission{
			ProblemID: params.ProblemID,
			UserID:    params.UserID,
			Code:      params.Code,
			Language:  params.Language,
		})
		if errTx != nil {
			return fmt.Errorf("s.storage.CreateSubmission: %w", errTx)
		}

		args := getArgsByParams(submission, params)
		if args == nil {
			logger.Errorf(ctx, "Send unsupported language", zap.String("language", string(params.Language)))
			_, errTx = s.storage.UpdateSubmissionStatus(ctx, submission.ID, daenerys.SubmissionStatus_STATUS_INTERNAL_ERROR)
			return errTx
		}

		_, errTx = s.riverClient.Insert(ctx, args, nil)
		if errTx != nil {
			return fmt.Errorf("s.riverClient.Insert: %w", errTx)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("storage.WithTransaction: %w", err)
	}

	return submission, nil
}

func getArgsByParams(submission *model.Submission, params SendSubmissionParams) river.JobArgs {
	switch params.Language {
	case daenerys.ProgrammingLanguage_LANGUAGE_PYTHON_3_12:
		return python.Args{
			SubmissionID:  submission.ID,
			ProblemID:     submission.ProblemID,
			Code:          submission.Code,
			PythonVersion: "3.12",
		}
	default:
		return nil
	}
}
