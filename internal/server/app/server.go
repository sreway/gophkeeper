package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/sreway/gophkeeper/internal/domain/models"
	grpcauth "github.com/sreway/gophkeeper/internal/server/grpc/auth"
	grpchealth "github.com/sreway/gophkeeper/internal/server/grpc/health"
	"github.com/sreway/gophkeeper/internal/server/grpc/interceptors"
	grpckeeper "github.com/sreway/gophkeeper/internal/server/grpc/keeper"
	"github.com/sreway/gophkeeper/internal/server/services/auth"
	"github.com/sreway/gophkeeper/internal/server/storage/postgres"

	"github.com/sreway/gophkeeper/internal/config"
	"github.com/sreway/gophkeeper/internal/server/services/keeper"
)

type (
	Server interface {
		Run(ctx context.Context) error
	}

	server struct {
		config     *config.Server
		gRPCServer *grpc.Server
	}
)

func NewServer(ctx context.Context, serverConfig *config.Server) (*server, error) {
	storage, err := postgres.New(ctx, serverConfig.Storage.Postgres)
	if err != nil {
		return nil, err
	}

	jwtManager := auth.NewJWTManager(serverConfig.Secret, serverConfig.TokenTTL)

	authService := auth.New(storage)
	keeperService := keeper.New(storage)

	switch {
	case errors.Is(err, models.ErrNotFound):
		id := uuid.New()
		err = keeperService.SetServerID(ctx, id)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	authInterceptor := interceptors.NewAuthInterceptor(jwtManager)

	serverOptions := []grpc.ServerOption{
		grpc.ChainStreamInterceptor(authInterceptor.Stream()),
		grpc.ChainUnaryInterceptor(authInterceptor.Unary()),
	}

	gRPCServer := grpc.NewServer(serverOptions...)

	grpckeeper.Register(gRPCServer, keeperService)
	grpchealth.Register(gRPCServer, keeperService)
	grpcauth.Register(gRPCServer, jwtManager, authService)

	return &server{
		config:     serverConfig,
		gRPCServer: gRPCServer,
	}, nil
}
