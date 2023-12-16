package app

import (
	"context"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) ServerHealthCheck(ctx context.Context) error {
	_, err := c.healthGRPC.HealthCheck(ctx, new(pb.HealthCheckRequest))
	if err != nil {
		return models.ErrServiceUnavailable
	}
	return nil
}
