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

	"github.com/CodefriendOrg/tirion/internal/usecase"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/CodefriendOrg/tirion/internal/app"
	"github.com/CodefriendOrg/tirion/internal/config"
	"github.com/CodefriendOrg/tirion/internal/pb/github.com/CodefriendOrg/tirion/pkg/tirion"
	"github.com/CodefriendOrg/tirion/internal/pkg/logger"
	"github.com/CodefriendOrg/tirion/internal/pkg/store"
	"google.golang.org/grpc"
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

	// Clients
	clients, err := config.NewGrpcClientsContainer(ctx)
	if err != nil {
		logger.Fatalf(ctx, "NewGrpcClientsContainer", zap.Error(err))
	}

	service := usecase.NewService(storage, clients.GetKingsguardclient())
	impl := app.NewImplementation(storage, service)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// GRPC
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.HandlerInterceptor()),
	)
	tirion.RegisterTirionServer(grpcServer, impl)

	// GRPC Gateway
	httpServer := &http.Server{
		Addr:              addrGateway,
		ReadHeaderTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run grpc server
	go func() {
		lis, errListen := net.Listen("tcp", addrGPRC)
		if errListen != nil {
			logger.Fatalf(ctx, "Failed to listen", zap.Error(errListen))
		}

		if cfg.Env != config.Production {
			reflection.Register(grpcServer)
		}

		logger.Infof(ctx, "App started",
			zap.String("grpc_address", addrGPRC),
		)
		if errListen = grpcServer.Serve(lis); errListen != nil {
			logger.Fatalf(ctx, "Failed to serve grpc", zap.Error(errListen))
		}
	}()

	// Run grpc-gateway
	go func() {
		mainMux := http.NewServeMux()

		// GRPC Gateway
		gatewayMux := runtime.NewServeMux()
		if errGw := tirion.RegisterTirionHandlerFromEndpoint(ctx, gatewayMux, addrGPRC, opts); errGw != nil {
			logger.Fatalf(ctx, "Failed to register gateway", zap.Error(errGw))
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
		if errServe := httpServer.ListenAndServe(); errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
			logger.Fatalf(ctx, "Failed to serve grpc gateway", zap.Error(errServe))
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
	storage.Close()
	logger.Infof(ctx, "App stopped gracefully...")
}
