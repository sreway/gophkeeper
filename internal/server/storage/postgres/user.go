package postgres

import (
	"context"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (s *storage) AddUser(ctx context.Context, user *models.User) error {
	stmt := "INSERT INTO users (id, email, hash_password) VALUES ($1, $2, $3)"
	_, err := s.pool.Exec(ctx, stmt, user.ID, user.Email, user.HashPassword)
	if err != nil {
		return errHandle(err)
	}

	return nil
}

func (s *storage) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, hash_password, created_at, updated_at FROM users WHERE email = $1"

	user := &models.User{
		Email: email,
	}

	err := s.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.HashPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, errHandle(err)
	}

	return user, nil
}
