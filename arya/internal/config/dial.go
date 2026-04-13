package config

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	keepAliveTime    = 10 * time.Second
	keepAliveTimeout = 1 * time.Minute
)

func dialService(_ context.Context, url string, timeout time.Duration, interns ...grpc.UnaryClientInterceptor) (grpc.ClientConnInterface, error) {
	interceptor := []grpc.UnaryClientInterceptor{
		TimeoutInterceptor(url, timeout),
	}

	if interns != nil {
		interceptor = append(interceptor, interns...)
	}

	conn, err := grpc.NewClient(url, getGRPCOptions(interceptor)...)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}
	return conn, nil
}

func dialContainerService(ctx context.Context, url string, timeoutString string) (grpc.ClientConnInterface, error) {
	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		return nil, fmt.Errorf("time.ParseDuration: %w", err)
	}

	return dialService(ctx, url, timeout)
}

func getGRPCOptions(interceptors []grpc.UnaryClientInterceptor, opts ...grpc.DialOption) []grpc.DialOption {
	dialOptions := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			interceptors...,
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepAliveTime,
			Timeout:             keepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	dialOptions = append(dialOptions, opts...)

	return dialOptions
}
