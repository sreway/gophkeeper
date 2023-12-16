package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/client/grpc/interceptors"
	"github.com/sreway/gophkeeper/internal/client/services"
	"github.com/sreway/gophkeeper/internal/client/services/keeper"
	"github.com/sreway/gophkeeper/internal/client/storage/sqlite"
	"github.com/sreway/gophkeeper/internal/config"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

type (
	Client interface {
		ID() uuid.UUID
		ActiveUser(ctx context.Context) (*models.User, error)
		SetMasterPassword(password string)
		IsMasterPasswordExists() bool
		ServerHealthCheck(ctx context.Context) error
		LoadActiveProfile(ctx context.Context) error
		Register(ctx context.Context, email, password string) error
		Login(ctx context.Context, email, password string) (*models.Profile, error)
		Logout(ctx context.Context) error
		ResetMasterPassword(ctx context.Context, email, password, masterPassword string) error
		CreateSecret(ctx context.Context, entry models.Entry) (*models.Secret, error)
		GetSecret(ctx context.Context, id uuid.UUID) (models.Entry, error)
		UpdateSecret(ctx context.Context, entry models.Entry) error
		DeleteSecret(ctx context.Context, id uuid.UUID) error
		ListSecret(ctx context.Context) ([]models.Entry, error)
		Sync(ctx context.Context, conflictCH chan []models.Entry, resolveCH chan int) error
	}

	client struct {
		id            uuid.UUID
		config        *config.Client
		masterKey     []byte
		tokenProvider services.TokenProvider
		keeper        keeper.Service
		keeperGRPC    pb.KeeperServiceClient
		authGRPC      pb.AuthServiceClient
		healthGRPC    pb.HealthServiceClient
	}
)

func (c *client) ID() uuid.UUID {
	return c.id
}

func (c *client) ActiveUser(ctx context.Context) (*models.User, error) {
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		return nil, err
	}
	return profile.User, nil
}

func (c *client) SetMasterPassword(password string) {
	c.config.MasterPassword = password
}

func (c *client) IsMasterPasswordExists() bool {
	return len(c.config.MasterPassword) > 0
}

func NewClient(ctx context.Context, clientConfig *config.Client) (*client, error) {
	var (
		err error
		id  uuid.UUID
	)

	tokenProvider := services.NewTokenProvider()
	authInterceptor := interceptors.NewAuthInterceptor(tokenProvider)

	grpcOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(authInterceptor.Unary()),
		grpc.WithChainStreamInterceptor(authInterceptor.Stream()),
	}

	conn, err := grpc.DialContext(ctx, clientConfig.Server, grpcOptions...)
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	storage, err := sqlite.New(ctx, clientConfig.Storage.SQLite)
	if err != nil {
		return nil, err
	}

	keeperService := keeper.New(storage)
	id, err = keeperService.GetClientID(ctx)
	switch {
	case errors.Is(err, models.ErrNotFound):
		id = uuid.New()
		err = keeperService.SetClientID(ctx, id)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	return &client{
		id:            id,
		config:        clientConfig,
		keeper:        keeperService,
		tokenProvider: tokenProvider,
		keeperGRPC:    pb.NewKeeperServiceClient(conn),
		authGRPC:      pb.NewAuthServiceClient(conn),
		healthGRPC:    pb.NewHealthServiceClient(conn),
	}, nil
}
