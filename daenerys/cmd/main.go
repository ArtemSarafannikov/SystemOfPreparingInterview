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

	"github.com/CodefriendOrg/daenerys/internal/pkg/docker"
	"github.com/CodefriendOrg/daenerys/internal/pkg/workers"
	"github.com/CodefriendOrg/daenerys/internal/usecase"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/moby/moby/client"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/CodefriendOrg/daenerys/internal/app"
	"github.com/CodefriendOrg/daenerys/internal/config"
	"github.com/CodefriendOrg/daenerys/internal/pb/github.com/CodefriendOrg/daenerys/pkg/daenerys"
	"github.com/CodefriendOrg/daenerys/internal/pkg/logger"
	"github.com/CodefriendOrg/daenerys/internal/pkg/store"
	"go.uber.org/zap"
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

	// Clients
	clients, err := config.NewGrpcClientsContainer(ctx)
	if err != nil {
		logger.Fatalf(ctx, "NewGrpcClientsContainer", zap.Error(err))
	}

	dbPool, err := cfg.Database.GetConn(ctx)
	if err != nil {
		logger.Fatalf(ctx, "Failed to connect database", zap.Error(err))
	}

	dockerDaemon, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatalf(ctx, "failed to connect docker daemon", zap.Error(err))
	}
	dockerClient := docker.New(dockerDaemon)

	storage := store.New(dbPool)
	riverClient, err := workers.RegisterRiverClient(
		ctx,
		dbPool,
		dockerClient,
		storage,
		clients.GetTirionclient(),
		cfg.Judge,
	)
	if err != nil {
		logger.Fatalf(ctx, "Failed to create river queue", zap.Error(err))
	}
	logger.Infof(ctx, "River queue started")

	service := usecase.NewService(storage, riverClient)
	impl := app.NewImplementation(storage, service)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// GRPC
	grpcServer := grpc.NewServer()
	daenerys.RegisterDaenerysServer(grpcServer, impl)

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
		err = daenerys.RegisterDaenerysHandlerFromEndpoint(ctx, gatewayMux, addrGPRC, opts)
		if err != nil {
			logger.Fatalf(ctx, "Failed to register gateway", zap.Error(err))
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
		if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf(ctx, "Failed to serve grpc gateway", zap.Error(err))
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
	if err = riverClient.Stop(ctx); err != nil {
		logger.Errorf(ctx, "Error to stop river queue", zap.Error(err))
	}
	logger.Infof(ctx, "App stopped gracefully...")
}
