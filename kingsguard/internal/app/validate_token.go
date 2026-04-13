package app

import (
	"context"

	kingsguard "github.com/CodefriendOrg/kingsguard/internal/pb/api"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/auth_token"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/serializer"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ValidateToken(ctx context.Context, req *kingsguard.ValidateTokenRequest) (*kingsguard.ValidateTokenResponse, error) {
	if err := validateValidateTokenRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims, err := auth_token.GetClaimsFromJWTToken(req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid access token")
	}

	user, err := i.storage.GetUserByID(ctx, claims.UserID)
	switch {
	case err == nil:
		// OK
	case utils.ErrorIsNotFound(err):
		return nil, status.Errorf(codes.Unauthenticated, "invalid access token")
	default:
		return nil, status.Errorf(codes.Internal, "storage.GetUserByID: %v", err)
	}

	if user.DeletedAt != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid access token")
	}

	return &kingsguard.ValidateTokenResponse{
		User: serializer.BuildUser(user),
	}, nil
}

func validateValidateTokenRequest(req *kingsguard.ValidateTokenRequest) error {
	err := validation.Validate(req, validation.Required.Error("request is required"))
	if err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.AccessToken, validation.Required),
	)
}
