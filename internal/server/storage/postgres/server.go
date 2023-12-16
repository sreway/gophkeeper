package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (s *storage) GetServerID(ctx context.Context) (uuid.UUID, error) {
	var id uuid.UUID

	query := "SELECT id FROM server"

	err := s.pool.QueryRow(ctx, query).Scan(&id)
	if err != nil {
		return id, errHandle(err)
	}

	return id, err
}

func (s *storage) SetServerID(ctx context.Context, id uuid.UUID) error {
	stmt := "INSERT INTO server (id) VALUES ($1)"
	_, err := s.pool.Exec(ctx, stmt, id)
	if err != nil {
		return errHandle(err)
	}
	return nil
}
