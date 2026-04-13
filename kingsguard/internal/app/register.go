package app

import (
	"context"

	"github.com/CodefriendOrg/kingsguard/internal/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kingsguard "github.com/CodefriendOrg/kingsguard/internal/pb/api"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/serializer"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/validator"
)

func (i *Implementation) Register(ctx context.Context, req *kingsguard.RegisterRequest) (*kingsguard.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "utils.HashPassword: %v", err)
	}

	_, err = i.storage.GetUserByUsername(ctx, req.Username)
	switch {
	case err == nil:
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	case utils.ErrorIsNotFound(err):
		// ОК
	default:
		return nil, status.Errorf(codes.Internal, "storage.GetUserByUsername: %v", err)
	}

	user, err := i.storage.CreateUser(ctx, req.Username, passwordHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "storage.CreateUser: %v", err)
	}

	return &kingsguard.RegisterResponse{
		User: serializer.BuildUser(user),
	}, nil
}

func validateRegisterRequest(req *kingsguard.RegisterRequest) error {
	err := validation.Validate(req, validation.Required.Error("request is required"))
	if err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Username, validation.Required, validation.Length(5, 64)),
		validation.Field(&req.Password, validation.Required, validator.Password()),
	)
}
