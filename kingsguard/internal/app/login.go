package app

import (
	"context"
	"time"

	kingsguard "github.com/CodefriendOrg/kingsguard/internal/pb/api"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/auth_token"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) Login(ctx context.Context, req *kingsguard.LoginRequest) (*kingsguard.LoginResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := i.storage.GetUserByUsername(ctx, req.Username)
	switch {
	case err == nil:
		// OK
	case utils.ErrorIsNotFound(err):
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	default:
		return nil, status.Errorf(codes.Internal, "storage.GetUserByUsername: %v", err)
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	claims := auth_token.JWTClaims{
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(time.Hour * 24),
	}

	token, err := claims.GenerateToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "claims.GenerateToken: %v", err)
	}

	return &kingsguard.LoginResponse{
		AccessToken: token,
		ExpiredAt:   timestamppb.New(claims.ExpiredAt),
	}, nil
}

func validateLoginRequest(req *kingsguard.LoginRequest) error {
	err := validation.Validate(req, validation.Required.Error("request is required"))
	if err != nil {
		return err
	}

	return validation.ValidateStruct(req,
		validation.Field(&req.Username, validation.Required),
		validation.Field(&req.Password, validation.Required),
	)
}
