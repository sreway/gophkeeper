package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

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

	rows, err := s.pool.Query(ctx, query, args...)
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
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	for _, secret := range secrets {
		stmt := "INSERT INTO secrets (id, owner, encrypted_value, hash, secret_type, is_deleted, created_at, updated_at)" +
			" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (id) DO UPDATE SET encrypted_value = $3, hash = $4, " +
			"secret_type = $5, is_deleted = $6, updated_at = $8"
		_, err = tx.Exec(ctx, stmt, secret.ID, userID, secret.EncryptedValue,
			secret.Hash, secret.Type, secret.IsDeleted, time.Now(), time.Now())
		if err != nil {
			return errHandle(err)
		}
	}

	return tx.Commit(ctx)
}
