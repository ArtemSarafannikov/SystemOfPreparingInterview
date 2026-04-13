package app

import (
	"context"

	kingsguard "github.com/CodefriendOrg/kingsguard/internal/pb/api"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/serializer"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetUser(ctx context.Context, req *kingsguard.GetUserRequest) (*kingsguard.GetUserResponse, error) {
	err := validateGetUserRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := i.storage.GetUserByID(ctx, uuid.MustParse(req.UserId))
	switch {
	case err == nil:
	// OK
	case utils.ErrorIsNotFound(err):
		return nil, status.Error(codes.NotFound, "user not found")
	default:
		return nil, status.Errorf(codes.Internal, "i.storage.GetUserByID: %v", err)
	}

	return &kingsguard.GetUserResponse{
		User: serializer.BuildUser(user),
	}, nil
}

func validateGetUserRequest(req *kingsguard.GetUserRequest) error {
	if err := validation.Validate(req, validation.Required.Error("request is required")); err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.UserId, validation.Required, is.UUID),
	)
}
