package app

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/utils"
	"github.com/CodefriendOrg/daenerys/internal/usecase"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SendSubmission .
func (i *Implementation) SendSubmission(ctx context.Context, req *daenerys.SendSubmissionRequest) (*daenerys.Submission, error) {
	if err := validateSendSubmissionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	submission, err := i.service.SendSubmission(ctx, usecase.SendSubmissionParams{
		ProblemID: uuid.MustParse(req.ProblemId),
		UserID:    uuid.MustParse(req.UserId),
		Code:      req.Code,
		Language:  req.Language,
	})
	if err != nil {
		return nil, utils.GRPCError(fmt.Errorf("i.service.SendSubmission: %w", err))
	}

	return submission.ToProto(), nil
}

func validateSendSubmissionRequest(req *daenerys.SendSubmissionRequest) error {
	if err := validation.Validate(req, validation.Required.Error("request is required")); err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.ProblemId, validation.Required, is.UUID),
		validation.Field(&req.UserId, validation.Required, is.UUID),
		validation.Field(&req.Code, validation.Required),
		validation.Field(&req.Language, validation.Required),
	)
}
