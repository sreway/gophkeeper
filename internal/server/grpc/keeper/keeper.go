package keeper

import (
	"google.golang.org/grpc"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/server/services/keeper"
)

type keeperGRPC struct {
	keeper keeper.Service
	pb.UnimplementedKeeperServiceServer
}

func Register(gRPCServer *grpc.Server, keeper keeper.Service) {
	pb.RegisterKeeperServiceServer(gRPCServer, &keeperGRPC{keeper: keeper})
}
