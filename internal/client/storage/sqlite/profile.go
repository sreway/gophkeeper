package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (s *storage) AddProfile(ctx context.Context, profile *models.Profile) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	userStmt, err := tx.Prepare("INSERT INTO users (id, email, hash_password) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer userStmt.Close()

	_, err = userStmt.ExecContext(ctx, profile.User.ID, profile.User.Email, profile.User.HashPassword)
	if err != nil {
		return errHandle(err)
	}

	sealStmt, err := tx.Prepare(
		"INSERT INTO seals (id, encrypted_shares, recovery_share, required_shares, total_shares, " +
			"hash_master_password, hash_key) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	defer sealStmt.Close()

	_, err = sealStmt.ExecContext(ctx,
		profile.Seal.ID, profile.Seal.EncryptedShares, profile.Seal.RecoveryShare, profile.Seal.RequiredShares,
		profile.Seal.TotalShares, profile.Seal.HashMasterPassword, profile.Seal.HashKey)
	if err != nil {
		return errHandle(err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	profileStmt, err := s.db.Prepare("INSERT INTO profiles (id, user_id, seal_id) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer profileStmt.Close()

	_, err = profileStmt.ExecContext(ctx, profile.ID, profile.User.ID, profile.Seal.ID)
	if err != nil {
		return errHandle(err)
	}

	return nil
}

func (s *storage) GetProfile(ctx context.Context, email string) (*models.Profile, error) {
	profile := &models.Profile{
		User:    new(models.User),
		Seal:    new(models.Seal),
		Session: new(models.Session),
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	profileStmt, err := tx.Prepare("SELECT profiles.id, user_id, seal_id, session_id FROM profiles " +
		"LEFT JOIN users ON profiles.user_id = users.id WHERE email = $1")
	if err != nil {
		return nil, err
	}
	defer profileStmt.Close()

	err = profileStmt.QueryRowContext(ctx, email).Scan(&profile.ID, &profile.User.ID, &profile.Seal.ID, &profile.Session.ID)
	if err != nil {
		return nil, errHandle(err)
	}

	userStmt, err := tx.Prepare("SELECT hash_password, created_at, updated_at FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer userStmt.Close()

	err = userStmt.QueryRowContext(ctx, profile.User.ID).Scan(
		&profile.User.HashPassword, &profile.User.CreatedAt, &profile.User.UpdatedAt)
	if err != nil {
		return nil, err
	}

	sealStmt, err := tx.Prepare("SELECT encrypted_shares, recovery_share, required_shares, total_shares, " +
		"hash_master_password, hash_key, created_at, updated_at FROM seals WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer sealStmt.Close()

	err = sealStmt.QueryRowContext(ctx, profile.Seal.ID).Scan(
		&profile.Seal.EncryptedShares,
		&profile.Seal.RecoveryShare,
		&profile.Seal.RequiredShares,
		&profile.Seal.TotalShares,
		&profile.Seal.HashMasterPassword,
		&profile.Seal.HashKey,
		&profile.Seal.CreatedAt,
		&profile.Seal.UpdatedAt,
	)
	if err != nil {
		return nil, errHandle(err)
	}

	sessionStmt, err := tx.Prepare("SELECT encrypted_token FROM sessions WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer sessionStmt.Close()

	err = sessionStmt.QueryRowContext(ctx, profile.Session.ID).Scan(
		&profile.Session.EncryptedToken,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errHandle(err)
	}

	return profile, tx.Commit()
}

func (s *storage) GetActiveProfile(ctx context.Context) (*models.Profile, error) {
	profile := &models.Profile{
		User:    new(models.User),
		Seal:    new(models.Seal),
		Session: new(models.Session),
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	profileStmt, err := tx.Prepare("SELECT profiles.id, user_id, seal_id, session_id FROM profiles " +
		"LEFT JOIN users ON profiles.user_id = users.id LEFT JOIN client on profiles.id = client.active_profile " +
		"WHERE profiles.id = client.active_profile")
	if err != nil {
		return nil, err
	}
	defer profileStmt.Close()

	err = profileStmt.QueryRowContext(ctx).Scan(&profile.ID, &profile.User.ID, &profile.Seal.ID, &profile.Session.ID)
	if err != nil {
		return nil, errHandle(err)
	}

	userStmt, err := tx.Prepare("SELECT email, hash_password, created_at, updated_at FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer userStmt.Close()

	err = userStmt.QueryRowContext(ctx, profile.User.ID).Scan(
		&profile.User.Email, &profile.User.HashPassword, &profile.User.CreatedAt, &profile.User.UpdatedAt)
	if err != nil {
		return nil, err
	}

	sealStmt, err := tx.Prepare("SELECT encrypted_shares, recovery_share, required_shares, total_shares, " +
		"hash_master_password, hash_key, created_at, updated_at FROM seals WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer sealStmt.Close()

	err = sealStmt.QueryRowContext(ctx, profile.Seal.ID).Scan(
		&profile.Seal.EncryptedShares,
		&profile.Seal.RecoveryShare,
		&profile.Seal.RequiredShares,
		&profile.Seal.TotalShares,
		&profile.Seal.HashMasterPassword,
		&profile.Seal.HashKey,
		&profile.Seal.CreatedAt,
		&profile.Seal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	sessionStmt, err := tx.Prepare("SELECT encrypted_token FROM sessions WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer sessionStmt.Close()

	err = sessionStmt.QueryRowContext(ctx, profile.Session.ID).Scan(
		&profile.Session.EncryptedToken,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errHandle(err)
	}

	return profile, tx.Commit()
}

func (s *storage) SetActiveProfile(ctx context.Context, clientID, profileID uuid.UUID) error {
	stmt := "UPDATE client SET active_profile = $1 WHERE id = $2"
	_, err := s.db.ExecContext(ctx, stmt, profileID, clientID)
	if err != nil {
		return errHandle(err)
	}

	return nil
}

func (s *storage) RemoveActiveProfile(ctx context.Context, id uuid.UUID) error {
	stmt := "UPDATE client SET active_profile = NULL WHERE id = $1"
	_, err := s.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return errHandle(err)
	}
	return nil
}

func (s *storage) AddSession(ctx context.Context, profileID uuid.UUID, session *models.Session) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	deleteStmt, err := tx.Prepare("DELETE FROM sessions WHERE id IN (SELECT session_id FROM profiles WHERE id = $1)")
	if err != nil {
		return err
	}
	defer deleteStmt.Close()

	_, err = deleteStmt.ExecContext(ctx, profileID)
	if err != nil {
		return errHandle(err)
	}

	sessionStmt, err := tx.Prepare("INSERT INTO sessions (id, encrypted_token) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer sessionStmt.Close()

	_, err = sessionStmt.ExecContext(ctx, session.ID, session.EncryptedToken)
	if err != nil {
		return errHandle(err)
	}

	profileStmt, err := tx.Prepare("UPDATE profiles SET session_id = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	defer profileStmt.Close()

	_, err = profileStmt.ExecContext(ctx, session.ID, profileID)
	if err != nil {
		return errHandle(err)
	}

	return tx.Commit()
}

func (s *storage) RemoveSession(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	profileStmt, err := tx.Prepare("UPDATE profiles SET session_id = NULL WHERE session_id = $1")
	if err != nil {
		return err
	}
	defer profileStmt.Close()

	_, err = profileStmt.ExecContext(ctx, id)
	if err != nil {
		return errHandle(err)
	}

	sessionStmt, err := tx.Prepare("DELETE FROM sessions WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = sessionStmt.ExecContext(ctx, id)
	if err != nil {
		return errHandle(err)
	}
	return tx.Commit()
}
