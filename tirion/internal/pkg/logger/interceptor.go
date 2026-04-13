package logger

import (
	"context"
	"path"

	"google.golang.org/grpc"
)

// HandlerInterceptor is a gRPC unary server interceptor that extracts
// the method name from the gRPC call and stores it in the context
// so that logger functions can automatically include the handler field.
func HandlerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// info.FullMethod is like "/tirion.Tirion/GetProblem"
		// path.Base extracts "GetProblem"
		method := path.Base(info.FullMethod)
		ctx = WithHandler(ctx, method)
		return handler(ctx, req)
	}
}
