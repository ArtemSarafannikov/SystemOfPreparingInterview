package app

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/model"
	"github.com/CodefriendOrg/tirion/internal/pkg/store"
	"github.com/CodefriendOrg/tirion/internal/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetProblem .
func (i *Implementation) GetProblem(ctx context.Context, req *tirion.GetProblemRequest) (*tirion.GetProblemResponse, error) {
	err := validateGetProblemRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	problem, err := i.storage.GetProblemByID(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return nil, utils.GRPCError(fmt.Errorf("i.storage.GetProblemByID: %w", err))
	}

	var testsProto []*tirion.Test
	if req.WithTests {
		tests, _, errTest := i.storage.ListTests(ctx, store.ListTestsParams{
			ProblemIDEq:    lo.ToPtr(uuid.MustParse(req.Id)),
			WithHidden:     true,
			OrderBy:        "created_at",
			OrderDirection: "ASC",
		})
		if errTest != nil {
			return nil, utils.GRPCError(fmt.Errorf("i.storage.ListTests: %w", errTest))
		}

		testsProto = lo.Map(tests, func(item *model.Test, _ int) *tirion.Test {
			return item.ToProto()
		})
	}

	return &tirion.GetProblemResponse{
		Problem: problem.ToProto(),
		Tests:   testsProto,
	}, nil
}

func validateGetProblemRequest(req *tirion.GetProblemRequest) error {
	if err := validation.Validate(req, validation.Required.Error("request is required")); err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Id, validation.Required, is.UUID),
	)
}
