package postgres

import (
	"context"
)

func (s *storage) HealthCheck(ctx context.Context) error {
	return s.pool.Ping(ctx)
}
