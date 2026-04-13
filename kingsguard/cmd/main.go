package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/CodefriendOrg/kingsguard/internal/usecase"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/reflection"

	"github.com/CodefriendOrg/kingsguard/internal/app"
	"github.com/CodefriendOrg/kingsguard/internal/config"
	kingsguard "github.com/CodefriendOrg/kingsguard/internal/pb/api"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/logger"
	"github.com/CodefriendOrg/kingsguard/internal/pkg/store"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoadFromFlag()
	if err := logger.Setup(cfg.Env); err != nil {
		panic(err)
	}

	defer logger.Logger.Sync() // nolint:errcheck
	logger.Infof(ctx, "Loaded config", zap.Any("env", cfg.Env))

	addrGPRC := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.GRPCPort))
	addrGateway := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.GatewayPort))

	storage, err := store.New(ctx, cfg.Database)
	if err != nil {
		logger.Fatalf(ctx, "Failed to connect database", zap.Error(err))
	}
	service := usecase.NewService(storage)

	impl := app.NewService(storage, service)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// GRPC
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.HandlerInterceptor()),
	)
	kingsguard.RegisterKingsguardServer(grpcServer, impl)

	// GRPC Gateway
	httpServer := &http.Server{
		Addr:              addrGateway,
		ReadHeaderTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run grpc server
	go func() {
		lis, err := net.Listen("tcp", addrGPRC)
		if err != nil {
			logger.Fatalf(ctx, "Failed to listen", zap.Error(err))
		}

		if cfg.Env != config.Production {
			reflection.Register(grpcServer)
		}

		logger.Infof(ctx, "App started",
			zap.String("grpc_address", addrGPRC),
		)
		if err = grpcServer.Serve(lis); err != nil {
			logger.Fatalf(ctx, "Failed to serve grpc", zap.Error(err))
		}
	}()

	// Run grpc-gateway
	go func() {
		mainMux := http.NewServeMux()

		// GRPC Gateway
		gatewayMux := runtime.NewServeMux()
		gwErr := kingsguard.RegisterKingsguardHandlerFromEndpoint(ctx, gatewayMux, addrGPRC, opts)
		if gwErr != nil {
			logger.Fatalf(ctx, "Failed to register gateway", zap.Error(gwErr))
		}
		mainMux.Handle("/", gatewayMux)

		// Swagger UI
		const swaggerPrefix = "/docs/"
		swaggerHandler := http.StripPrefix(swaggerPrefix, http.FileServer(http.Dir("api/swagger")))
		mainMux.Handle(swaggerPrefix, swaggerHandler)

		httpServer.Handler = mainMux

		logger.Infof(ctx, "Starting HTTP gateway and swagger",
			zap.String("gateway_address", addrGateway),
			zap.String("swagger_ui", addrGateway+swaggerPrefix),
		)
		if serveErr := httpServer.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Fatalf(ctx, "Failed to serve grpc gateway", zap.Error(serveErr))
		}
	}()

	<-stop

	logger.Infof(ctx, "Stopping app graceful shutdown...", zap.Duration("timeout", cfg.Server.GracefulShutdownTimeout))
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulShutdownTimeout)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		logger.Errorf(ctx, "HTTP server shutdown error", zap.Error(err))
	}

	go func() {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			grpcServer.Stop()
			logger.Errorf(ctx, "Deadline exceeded, app stopped forcibly")
		}
	}()

	grpcServer.GracefulStop()
	logger.Infof(ctx, "App stopped gracefully...")
}
