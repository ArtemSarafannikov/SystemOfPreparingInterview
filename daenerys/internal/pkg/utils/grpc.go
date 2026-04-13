package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCError .
func GRPCError(err error) error {
	switch {
	case err == nil:
		return nil
	case ErrorIsNotFound(err):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
