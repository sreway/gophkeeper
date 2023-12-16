package auth

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/server/services/auth"
)

type authGRPC struct {
	auth auth.Service
	jwt  auth.JWTManager
	pb.UnimplementedAuthServiceServer
}

func (a *authGRPC) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}

func Register(gRPCServer *grpc.Server, jwt auth.JWTManager, auth auth.Service) {
	pb.RegisterAuthServiceServer(gRPCServer, &authGRPC{auth: auth, jwt: jwt})
}
