package app

import (
	"context"

	desc "github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetSubmission .
func (i *Implementation) GetSubmission(ctx context.Context, req *desc.GetSubmissionRequest) (*desc.Submission, error) {
	err := validateGetSubmissionRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	submission, err := i.storage.GetSubmissionByID(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return nil, utils.GRPCError(err)
	}

	return submission.ToProto(), nil
}

func validateGetSubmissionRequest(req *desc.GetSubmissionRequest) error {
	err := validation.Validate(req, validation.Required.Error("request is required"))
	if err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Id, validation.Required, is.UUID),
	)
}
