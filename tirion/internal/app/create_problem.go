package app

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/problemlimits"
	"github.com/CodefriendOrg/tirion/internal/pkg/utils"
	"github.com/CodefriendOrg/tirion/internal/usecase"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateProblem .
func (i *Implementation) CreateProblem(ctx context.Context, req *tirion.CreateProblemRequest) (*tirion.Problem, error) {
	err := validateCreateProblemRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	problem, err := i.service.CreateProblem(ctx, usecase.CreateProblemParams{
		Summary:       req.Summary,
		Description:   req.Description,
		AuthorID:      uuid.MustParse(req.AuthorId),
		TimeLimitMs:   req.TimeLimitMs,
		MemoryLimitKb: req.MemoryLimitKb,
	})
	if err != nil {
		return nil, utils.GRPCError(fmt.Errorf("i.service.CreateProblem: %w", err))
	}

	return problem, nil
}

func validateCreateProblemRequest(req *tirion.CreateProblemRequest) error {
	if err := validation.Validate(req, validation.Required.Error("request is required")); err != nil {
		return err
	}

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Summary, validation.Required),
		validation.Field(&req.Description, validation.Required),
		validation.Field(&req.AuthorId, validation.Required, is.UUID),
	); err != nil {
		return err
	}
	return problemlimits.Validate(req.TimeLimitMs, req.MemoryLimitKb)
}
