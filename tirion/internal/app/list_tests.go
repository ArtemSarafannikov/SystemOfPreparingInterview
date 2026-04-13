package app

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
	"github.com/CodefriendOrg/tirion/internal/pkg/store"
	"github.com/CodefriendOrg/tirion/internal/pkg/utils"
	"github.com/CodefriendOrg/tirion/internal/pkg/validator"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListTests .
func (i *Implementation) ListTests(ctx context.Context, req *tirion.ListTestsRequest) (*tirion.ListTestsResponse, error) {
	if err := validateListTestsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	params := store.ListTestsParams{
		ProblemIDEq: lo.ToPtr(uuid.MustParse(req.Filter.ProblemIdEq)),
		WithHidden:  req.Filter.WithHidden,
	}
	params.Limit, params.Offset = store.GetLimitOffsetFromProtoPagination(req.Pagination)
	switch req.OrderBy {
	case tirion.ListTestsRequest_CREATED_AT:
		params.OrderBy = "created_at"
	}
	params.OrderDirection = req.OrderDirection.String()

	tests, count, err := i.storage.ListTests(ctx, params)
	if err != nil {
		return nil, utils.GRPCError(fmt.Errorf("i.storage.ListTests: %w", err))
	}
	return &tirion.ListTestsResponse{
		Tests: lo.Map(tests, func(item *model.Test, _ int) *tirion.Test {
			return item.ToProto()
		}),
		TotalItems: count,
	}, nil
}

func validateListTestsRequest(req *tirion.ListTestsRequest) error {
	if err := validation.Validate(req, validation.Required.Error("request is required")); err != nil {
		return err
	}

	err := validation.ValidateStruct(req,
		validation.Field(&req.Filter, validation.Required),
		validation.Field(&req.Pagination, validation.Required, validation.By(validator.ValidatePagination)),
	)
	if err != nil {
		return err
	}

	return validation.ValidateStruct(req.Filter,
		validation.Field(&req.Filter.ProblemIdEq, validation.Required, is.UUID),
	)
}
