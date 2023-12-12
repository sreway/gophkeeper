package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (s *storage) AddSeal(ctx context.Context, userID uuid.UUID, seal *models.Seal) error {
	stmt := "INSERT INTO seals (id, user_id, encrypted_shares, recovery_share, total_shares, required_shares," +
		" hash_master_password, hash_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := s.pool.Exec(ctx, stmt, seal.ID, userID, seal.EncryptedShares, seal.RecoveryShare, seal.TotalShares,
		seal.RequiredShares, seal.HashMasterPassword, seal.HashKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) GetSeal(ctx context.Context, userID uuid.UUID) (*models.Seal, error) {
	query := "SELECT id, encrypted_shares, recovery_share, total_shares, required_shares, hash_master_password, " +
		"hash_key, created_at, updated_at FROM seals WHERE user_id = $1"

	seal := new(models.Seal)

	err := s.pool.QueryRow(ctx, query, userID).Scan(&seal.ID, &seal.EncryptedShares,
		&seal.RecoveryShare, &seal.TotalShares, &seal.RequiredShares, &seal.HashMasterPassword, &seal.HashKey,
		&seal.CreatedAt, &seal.UpdatedAt)
	if err != nil {
		return nil, errHandle(err)
	}

	return seal, nil
}

func (s *storage) UpdateSeal(ctx context.Context, userID uuid.UUID, seal *models.Seal) error {
	stmt := "UPDATE seals SET encrypted_shares = $2, recovery_share = $3, required_shares = $4,  total_shares = $5," +
		"  hash_master_password = $6, hash_key = $7, updated_at = $8 WHERE user_id = $1"

	_, err := s.pool.Exec(ctx, stmt, userID, seal.EncryptedShares, seal.RecoveryShare, seal.RequiredShares,
		seal.TotalShares, seal.HashMasterPassword, seal.HashKey, time.Now())
	if err != nil {
		return err
	}
	return nil
}
