package app

import (
	"context"

	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateSubmissionStatus .
func (i *Implementation) UpdateSubmissionStatus(ctx context.Context, req *daenerys.UpdateSubmissionStatusRequest) (*daenerys.Submission, error) {
	if err := validateUpdateSubmissionStatusRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	submission, err := i.storage.UpdateSubmissionStatus(ctx, uuid.MustParse(req.Id), req.Status)
	if err != nil {
		return nil, utils.GRPCError(err)
	}

	return submission.ToProto(), nil
}

func validateUpdateSubmissionStatusRequest(req *daenerys.UpdateSubmissionStatusRequest) error {
	if err := validation.Validate(req, validation.Required.Error("request is required")); err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Id, validation.Required, is.UUID),
		validation.Field(&req.Status, validation.Required),
	)
}
