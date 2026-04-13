package kingsguard_helper

import (
	"context"
	"fmt"

	"github.com/CodefriendOrg/arya/internal/pb/github.com/CodefriendOrg/kingsguard/pkg/kingsguard"
	"github.com/CodefriendOrg/arya/internal/pkg/logger"
	"github.com/CodefriendOrg/arya/internal/pkg/user_error"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service .
type Service struct {
	kingsguardClient kingsguard.KingsguardClient
}

// NewService .
func NewService(kingsguardClient kingsguard.KingsguardClient) *Service {
	return &Service{
		kingsguardClient: kingsguardClient,
	}
}

// GetUserFromAccessTokenWoError .
func (s *Service) GetUserFromAccessTokenWoError(ctx context.Context, token string) *kingsguard.User {
	resp, err := s.kingsguardClient.ValidateToken(ctx, &kingsguard.ValidateTokenRequest{
		AccessToken: token,
	})
	switch status.Code(err) {
	case codes.OK:
		return resp.User
	case codes.Unauthenticated:
		return nil
	default:
		logger.Errorf(ctx, "GetUserFromAccessToken: s.kingsguardClient.ValidateToken", zap.Error(err))
		return nil
	}
}

// Register .
func (s *Service) Register(ctx context.Context, username string, password string) (*kingsguard.User, error) {
	resp, err := s.kingsguardClient.Register(ctx, &kingsguard.RegisterRequest{
		Username: username,
		Password: password,
	})
	switch status.Code(err) {
	case codes.OK:
		return resp.User, nil
	case codes.AlreadyExists:
		return nil, user_error.WithoutLoggerMessage(user_error.UserAlreadyExists)
	default:
		return nil, user_error.New(user_error.InternalError, "s.kingsguardClient.Register: %v", err)
	}
}

// Login .
func (s *Service) Login(ctx context.Context, username string, password string) (string, error) {
	resp, err := s.kingsguardClient.Login(ctx, &kingsguard.LoginRequest{
		Username: username,
		Password: password,
	})
	switch status.Code(err) {
	case codes.OK:
		return resp.AccessToken, nil
	case codes.Unauthenticated:
		return "", user_error.New(user_error.WrongPassword, "неверный логин или пароль")
	default:
		return "", fmt.Errorf("s.kingsguardClient.Login: %w", err)
	}
}
