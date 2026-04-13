package app

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
	"github.com/CodefriendOrg/tirion/internal/pkg/utils"
	"github.com/CodefriendOrg/tirion/internal/pkg/validator"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListProblems .
func (i *Implementation) ListProblems(ctx context.Context, req *tirion.ListProblemsRequest) (*tirion.ListProblemsResponse, error) {
	err := validateListProblemsRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	problems, count, err := i.service.ListProblems(ctx, req)
	if err != nil {
		return nil, utils.GRPCError(fmt.Errorf("i.service.ListProblems: %w", err))
	}

	return &tirion.ListProblemsResponse{
		Problems: lo.Map(problems, func(item *model.Problem, _ int) *tirion.Problem {
			return item.ToProto()
		}),
		TotalItems: count,
	}, nil
}

func validateListProblemsRequest(req *tirion.ListProblemsRequest) error {
	err := validation.Validate(req, validation.Required.Error("request is required"))
	if err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Filter, validation.Required),
		validation.Field(&req.Pagination, validation.Required, validation.By(validator.ValidatePagination)),
	)
}
