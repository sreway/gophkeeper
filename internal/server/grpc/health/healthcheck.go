package health

import (
	"context"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
)

func (k *healthGRPC) HealthCheck(ctx context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	err := k.keeper.HealthCheck(ctx)
	if err != nil {
		return nil, err
	}

	response := &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}

	return response, nil
}
