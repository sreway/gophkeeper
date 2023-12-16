package sqlite

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (s *storage) AddSecret(ctx context.Context, secret *models.Secret) error {
	stmt := "INSERT INTO secrets (id, owner, encrypted_value, hash, secret_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := s.db.ExecContext(ctx, stmt, secret.ID, secret.Owner, secret.EncryptedValue, secret.Hash,
		secret.Type, time.Now(), time.Now())
	if err != nil {
		return errHandle(err)
	}
	return nil
}

func (s *storage) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	stmt := "UPDATE secrets SET encrypted_value = $1, hash = $2, secret_type =  $3, is_deleted = $4, " +
		"updated_at = $5 WHERE id = $6"
	_, err := s.db.ExecContext(ctx, stmt, secret.EncryptedValue, secret.Hash, secret.Type,
		secret.IsDeleted, time.Now(), secret.ID)
	if err != nil {
		return errHandle(err)
	}
	return nil
}

func (s *storage) DeleteSecret(ctx context.Context, id, ownerID uuid.UUID) error {
	stmt := "DELETE FROM secrets WHERE id = $1 AND owner = $2"
	_, err := s.db.ExecContext(ctx, stmt, id, ownerID)
	if err != nil {
		return errHandle(err)
	}
	return nil
}

func (s *storage) GetSecret(ctx context.Context, id, ownerID uuid.UUID) (*models.Secret, error) {
	secret := &models.Secret{
		ID:    id,
		Owner: ownerID,
	}
	query := "SELECT encrypted_value, hash, secret_type, created_at, updated_at FROM secrets " +
		"WHERE id = $1 AND owner = $2 AND is_deleted = false"
	err := s.db.QueryRowContext(ctx, query, id, ownerID).Scan(&secret.EncryptedValue, &secret.Hash, &secret.Type,
		&secret.CreatedAt, &secret.UpdatedAt)
	if err != nil {
		return nil, errHandle(err)
	}

	return secret, nil
}

func (s *storage) ListSecret(ctx context.Context, ownerID uuid.UUID) ([]*models.Secret, error) {
	secrets := make([]*models.Secret, 0)
	query := "SELECT id, encrypted_value, hash, secret_type, created_at, updated_at FROM secrets " +
		"WHERE owner = $1 AND is_deleted = false"
	rows, err := s.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		secret := models.Secret{
			Owner: ownerID,
		}

		if err = rows.Scan(&secret.ID, &secret.EncryptedValue, &secret.Hash,
			&secret.Type, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
			return nil, err
		}

		secrets = append(secrets, &secret)
	}

	return secrets, nil
}

func (s *storage) ListUpdatedSecrets(ctx context.Context, ownerID uuid.UUID, updatedAfter *time.Time) ([]*models.Secret, error) {
	secrets := make([]*models.Secret, 0)
	query := "SELECT id, encrypted_value, hash, secret_type, is_deleted, created_at, updated_at FROM secrets " +
		"WHERE owner = $1"
	var args []interface{}
	args = append(args, ownerID)

	if updatedAfter != nil {
		query += " AND updated_at > $2"
		args = append(args, updatedAfter)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		secret := models.Secret{
			Owner: ownerID,
		}

		if err = rows.Scan(&secret.ID, &secret.EncryptedValue, &secret.Hash,
			&secret.Type, &secret.IsDeleted, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
			return nil, err
		}

		secrets = append(secrets, &secret)
	}

	return secrets, nil
}

func (s *storage) BatchUpdateSecrets(ctx context.Context, userID uuid.UUID, secrets []*models.Secret) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	deleteItems := make([]*models.Secret, 0)
	updateItems := make([]*models.Secret, 0)

	for _, secret := range secrets {
		if secret.IsDeleted {
			deleteItems = append(deleteItems, secret)
		} else {
			updateItems = append(updateItems, secret)
		}
	}

	for _, secret := range updateItems {
		stmt := "INSERT INTO secrets (id, owner, encrypted_value, hash, secret_type, created_at, updated_at) VALUES " +
			"($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id) DO UPDATE SET encrypted_value = $3, hash = $4, " +
			"secret_type = $5, updated_at = $7"
		_, err = tx.ExecContext(ctx, stmt, secret.ID, userID, secret.EncryptedValue,
			secret.Hash, secret.Type, time.Now(), time.Now())
		if err != nil {
			return errHandle(err)
		}
	}

	for _, secret := range deleteItems {
		stmt := "DELETE FROM secrets WHERE id = $1 AND owner = $2"
		_, err = tx.ExecContext(ctx, stmt, secret.ID, userID)
		if err != nil {
			return errHandle(err)
		}
	}

	return tx.Commit()
}
