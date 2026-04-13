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

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/CodefriendOrg/arya/internal/app/graph"
	"github.com/CodefriendOrg/arya/internal/app/graph/directives"
	"github.com/CodefriendOrg/arya/internal/app/graph/generated"
	"github.com/CodefriendOrg/arya/internal/config"
	"github.com/CodefriendOrg/arya/internal/middleware"
	"github.com/CodefriendOrg/arya/internal/pkg/helper/daenerys_helper"
	"github.com/CodefriendOrg/arya/internal/pkg/helper/kingsguard_helper"
	"github.com/CodefriendOrg/arya/internal/pkg/helper/tirion_helper"
	"github.com/CodefriendOrg/arya/internal/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoadFromFlag()
	if err := logger.Setup(cfg.Env); err != nil {
		panic(err)
	}

	addr := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.GraphQLPort))

	// Clients
	clients, err := config.NewGrpcClientsContainer(ctx)
	if err != nil {
		logger.Fatalf(ctx, "NewGrpcClientsContainer", zap.Error(err))
	}

	// ** Helpers **
	kingsguardHelper := kingsguard_helper.NewService(clients.GetKingsguardclient())
	daenerysHelper := daenerys_helper.NewService(clients.GetDaenerysclient())
	tirionHelper := tirion_helper.NewService(clients.GetTirionclient())

	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: graph.NewResolver(
			kingsguardHelper,
			daenerysHelper,
			tirionHelper,
		),
		Directives: generated.DirectiveRoot{
			Auth: directives.Auth,
		},
	}))

	allowedWSHosts := cfg.Server.WebsocketAllowedHosts
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				for _, h := range allowedWSHosts {
					if r.Host == h {
						return true
					}
				}
				return false
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	srv.SetErrorPresenter(middleware.ErrorMiddleware)
	srv.SetRecoverFunc(middleware.PanicMiddleware)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	httpServer := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		mux := http.NewServeMux()

		mux.Handle("/query", middleware.AuthMiddleware(ctx, kingsguardHelper, srv))

		if cfg.Env != config.Production {
			mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
			logger.Infof(ctx, "running GraphQL playground", zap.String("address", addr))
		}

		c := cors.New(cors.Options{
			AllowedOrigins:   cfg.Server.CORSAllowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept"},
			AllowCredentials: true,
			Debug:            true,
		})

		logger.Infof(ctx, "App started", zap.String("graphql_address", addr))

		httpServer.Handler = c.Handler(mux)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf(ctx, "Failed to serve graphql", zap.Error(err))
		}
	}()

	<-stop

	logger.Infof(ctx, "Stopping app graceful shutdown...", zap.Duration("timeout", cfg.Server.GracefulShutdownTimeout))
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf(ctx, "Server shutdown error", zap.Error(err))
	}

	logger.Infof(ctx, "App stopped gracefully...")
}
