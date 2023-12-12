package sqlite

import (
	"context"

	"github.com/google/uuid"
)

func (s *storage) GetClientID(ctx context.Context) (uuid.UUID, error) {
	var id uuid.UUID

	query := "SELECT id FROM client"

	err := s.db.QueryRowContext(ctx, query).Scan(&id)
	if err != nil {
		return id, errHandle(err)
	}

	return id, err
}

func (s *storage) SetClientID(ctx context.Context, id uuid.UUID) error {
	stmt := "INSERT INTO client (id) VALUES ($1)"
	_, err := s.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return errHandle(err)
	}
	return nil
}
