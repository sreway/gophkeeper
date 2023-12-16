package sqlite

import (
	"context"
	"time"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (s *storage) UpdateSeal(ctx context.Context, seal *models.Seal) error {
	stmt := "UPDATE seals SET encrypted_shares = $1, recovery_share = $2, required_shares = $3, total_shares = $4," +
		" hash_master_password = $5, hash_key = $6, updated_at = $7 WHERE id = $8"

	_, err := s.db.ExecContext(ctx, stmt, seal.EncryptedShares, seal.RecoveryShare, seal.RequiredShares,
		seal.TotalShares, seal.HashMasterPassword, seal.HashKey, time.Now(), seal.ID)
	if err != nil {
		return errHandle(err)
	}

	return nil
}
