package sqlite

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *storage) GetLastSync(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	var lastSync time.Time

	query := "SELECT last_sync FROM sync WHERE user_id = $1"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&lastSync)
	if err != nil {
		return nil, errHandle(err)
	}
	return &lastSync, nil
}

func (s *storage) UpdateLastSync(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	lastSync := time.Now()
	stmt := "INSERT INTO sync (user_id, last_sync) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET last_sync = $2"
	_, err := s.db.ExecContext(ctx, stmt, userID, lastSync)
	if err != nil {
		return nil, errHandle(err)
	}

	return &lastSync, nil
}
