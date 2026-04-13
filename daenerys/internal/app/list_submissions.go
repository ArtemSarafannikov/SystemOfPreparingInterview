package app

import (
	"context"
	"errors"
	"fmt"

	desc "github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/model"
	"github.com/CodefriendOrg/daenerys/internal/pkg/store"
	"github.com/CodefriendOrg/daenerys/internal/pkg/utils"
	"github.com/CodefriendOrg/daenerys/internal/pkg/validator"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListSubmissions .
func (i *Implementation) ListSubmissions(ctx context.Context, req *desc.ListSubmissionsRequest) (*desc.ListSubmissionsResponse, error) {
	err := validateListSubmissionsRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	params := store.ListSubmissionsParams{
		IDIn: lo.Map(req.Filter.SubmissionIdIn, func(id string, _ int) uuid.UUID {
			return uuid.MustParse(id)
		}),
	}
	if req.Filter.UserIdEq != nil {
		params.UserIDEq = lo.ToPtr(uuid.MustParse(*req.Filter.UserIdEq))
	}
	if req.Filter.ProblemIdEq != nil {
		params.ProblemIDEq = lo.ToPtr(uuid.MustParse(*req.Filter.ProblemIdEq))
	}
	params.Limit, params.Offset = store.GetLimitOffsetFromProtoPagination(req.Pagination)

	switch req.OrderBy {
	case desc.ListSubmissionsRequest_CREATED_AT:
		params.OrderBy = "created_at"
	}
	params.OrderDirection = req.OrderDirection.String()

	submissions, count, err := i.storage.ListSubmissions(ctx, params)
	if err != nil {
		return nil, utils.GRPCError(fmt.Errorf("i.storage.ListSubmissions: %w", err))
	}

	return &desc.ListSubmissionsResponse{
		Submissions: lo.Map(submissions, func(item *model.Submission, _ int) *desc.Submission {
			return item.ToProto()
		}),
		TotalItems: count,
	}, nil
}

func validateListSubmissionsRequest(req *desc.ListSubmissionsRequest) error {
	err := validation.Validate(req, validation.Required.Error("request is required"))
	if err != nil {
		return err
	}

	err = validation.ValidateStruct(req,
		validation.Field(&req.Filter, validation.Required),
		validation.Field(&req.Pagination, validation.Required, validation.By(validator.ValidatePagination)),
	)
	if err != nil {
		return err
	}

	filter := req.Filter
	err = validation.ValidateStruct(filter,
		validation.Field(&filter.SubmissionIdIn, validation.Each(validation.Required, is.UUID)),
		validation.Field(&filter.ProblemIdEq, validation.NilOrNotEmpty, is.UUID),
		validation.Field(&filter.UserIdEq, validation.NilOrNotEmpty, is.UUID),
	)
	if err != nil {
		return err
	}

	if len(filter.SubmissionIdIn) == 0 &&
		filter.ProblemIdEq == nil &&
		filter.UserIdEq == nil {
		return errors.New("filter: must have at least one field")
	}

	return nil
}
