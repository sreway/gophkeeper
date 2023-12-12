package health

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/server/services/keeper"
)

type healthGRPC struct {
	keeper keeper.Service
	pb.UnimplementedHealthServiceServer
}

func Register(gRPCServer *grpc.Server, keeper keeper.Service) {
	pb.RegisterHealthServiceServer(gRPCServer, &healthGRPC{keeper: keeper})
}

func (k *healthGRPC) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}
