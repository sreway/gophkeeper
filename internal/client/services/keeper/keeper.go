package keeper

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/client/services"
	"github.com/sreway/gophkeeper/internal/domain/models"
	"github.com/sreway/gophkeeper/internal/lib/crypt"
)

const (
	keyLength   = 32
	threshold   = 2
	totalShares = 3
)

type (
	Service interface {
		GetClientID(ctx context.Context) (uuid.UUID, error)
		SetClientID(ctx context.Context, id uuid.UUID) error
		CreateKeySeals(key []byte, masterPassword string) (local, remote *models.Seal, err error)
		DecryptKeyShares(encryptedShares []byte, masterPassword string) (shares [][]byte, err error)
		DecryptMasterKey(seal *models.Seal, masterPassword string) (masterKey []byte, err error)
		RecoveryMasterKey(shares [][]byte) (masterKey []byte, err error)
		EncryptData(data, masterKey []byte) ([]byte, error)
		DecryptData(encryptedData, masterKey []byte) ([]byte, error)
		CreateProfile(ctx context.Context, user *models.User, seal *models.Seal) (*models.Profile, error)
		GetProfile(ctx context.Context, email string) (*models.Profile, error)
		GetActiveProfile(ctx context.Context) (*models.Profile, error)
		SetActiveProfile(ctx context.Context, clientID, profileID uuid.UUID) error
		RemoveActiveProfile(ctx context.Context, clientID uuid.UUID) error
		CreateSession(ctx context.Context, profileID uuid.UUID, encryptedToken []byte) error
		RemoveSession(ctx context.Context, id uuid.UUID) error
		UpdateSeal(ctx context.Context, seal *models.Seal) error
		AddSecret(ctx context.Context, secret *models.Secret) error
		GetSecret(ctx context.Context, id, ownerID uuid.UUID) (*models.Secret, error)
		UpdateSecret(ctx context.Context, secret *models.Secret) error
		DeleteSecret(ctx context.Context, id, ownerID uuid.UUID) error
		ListSecret(ctx context.Context, ownerID uuid.UUID) ([]*models.Secret, error)
		ListUpdatedSecrets(ctx context.Context, ownerID uuid.UUID, updatedAfter *time.Time) ([]*models.Secret, error)
		BatchUpdateSecrets(ctx context.Context, userID uuid.UUID, secrets []*models.Secret) error
		GetLastSync(ctx context.Context, userID uuid.UUID) (*time.Time, error)
		UpdateLastSync(ctx context.Context, userID uuid.UUID) (*time.Time, error)
	}

	service struct {
		profileProvider services.ProfileProvider
		clientProvider  services.ClientProvider
		secretProvider  services.SecretProvider
		syncProvider    services.SyncProvider
	}
)

func (s *service) CreateProfile(ctx context.Context, user *models.User, seal *models.Seal) (*models.Profile, error) {
	profile := &models.Profile{
		ID:   uuid.New(),
		User: user,
		Seal: seal,
	}
	return profile, s.profileProvider.AddProfile(ctx, profile)
}

func (s *service) GetProfile(ctx context.Context, email string) (*models.Profile, error) {
	return s.profileProvider.GetProfile(ctx, email)
}

func (s *service) GetActiveProfile(ctx context.Context) (*models.Profile, error) {
	return s.profileProvider.GetActiveProfile(ctx)
}

func (s *service) SetActiveProfile(ctx context.Context, clientID, profileID uuid.UUID) error {
	return s.profileProvider.SetActiveProfile(ctx, clientID, profileID)
}

func (s *service) RemoveActiveProfile(ctx context.Context, clientID uuid.UUID) error {
	return s.profileProvider.RemoveActiveProfile(ctx, clientID)
}

func (s *service) GetClientID(ctx context.Context) (uuid.UUID, error) {
	return s.clientProvider.GetClientID(ctx)
}

func (s *service) SetClientID(ctx context.Context, id uuid.UUID) error {
	return s.clientProvider.SetClientID(ctx, id)
}

func (s *service) CreateKeySeals(key []byte, masterPassword string) (local, remote *models.Seal, err error) {
	if key == nil {
		key, err = crypt.GenerateMasterKey(keyLength)
		if err != nil {
			return nil, nil, err
		}
	}

	shares, err := crypt.CreateKeyShares(key, threshold, totalShares)
	if err != nil {
		return nil, nil, err
	}

	masterPasswordHash := crypt.HashSum([]byte(masterPassword), nil)

	sealShares, err := json.Marshal(shares[:2])
	if err != nil {
		return nil, nil, err
	}

	encryptedShares, err := crypt.EncryptGCM(sealShares, []byte(masterPassword))
	if err != nil {
		return nil, nil, err
	}

	local = &models.Seal{
		ID:                 uuid.New(),
		EncryptedShares:    encryptedShares,
		RecoveryShare:      shares[0],
		TotalShares:        totalShares,
		RequiredShares:     threshold,
		HashKey:            crypt.HashSum(key, nil),
		HashMasterPassword: masterPasswordHash,
	}

	remote = &models.Seal{
		ID:                 uuid.New(),
		EncryptedShares:    encryptedShares,
		RecoveryShare:      shares[len(shares)-1],
		TotalShares:        totalShares,
		RequiredShares:     threshold,
		HashKey:            crypt.HashSum(key, nil),
		HashMasterPassword: masterPasswordHash,
	}

	return local, remote, nil
}

func (s *service) DecryptKeyShares(encryptedShares []byte, masterPassword string) (shares [][]byte, err error) {
	decryptedShares, err := crypt.DecryptGCM(encryptedShares, []byte(masterPassword))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(decryptedShares, &shares)
	if err != nil {
		return nil, err
	}

	return shares, nil
}

func (s *service) DecryptMasterKey(seal *models.Seal, masterPassword string) (masterKey []byte, err error) {
	shares, err := s.DecryptKeyShares(seal.EncryptedShares, masterPassword)
	if err != nil {
		return nil, err
	}
	hashMasterPassword := crypt.HashSum([]byte(masterPassword), nil)

	if hashMasterPassword != seal.HashMasterPassword {
		return nil, models.ErrInvalidMasterPassword
	}
	masterKey, err = crypt.RecoveryKeyShares(shares)
	if err != nil {
		return nil, err
	}

	return
}

func (s *service) RecoveryMasterKey(shares [][]byte) (masterKey []byte, err error) {
	return crypt.RecoveryKeyShares(shares)
}

func (s *service) CreateSession(ctx context.Context, profileID uuid.UUID, encryptedToken []byte) error {
	session := &models.Session{
		ID:             uuid.New(),
		EncryptedToken: encryptedToken,
	}

	return s.profileProvider.AddSession(ctx, profileID, session)
}

func (s *service) RemoveSession(ctx context.Context, id uuid.UUID) error {
	return s.profileProvider.RemoveSession(ctx, id)
}

func (s *service) UpdateSeal(ctx context.Context, seal *models.Seal) error {
	return s.profileProvider.UpdateSeal(ctx, seal)
}

func (s *service) EncryptData(data, masterKey []byte) ([]byte, error) {
	return crypt.EncryptGCM(data, masterKey)
}

func (s *service) DecryptData(encryptedData, masterKey []byte) ([]byte, error) {
	return crypt.DecryptGCM(encryptedData, masterKey)
}

func (s *service) AddSecret(ctx context.Context, secret *models.Secret) error {
	return s.secretProvider.AddSecret(ctx, secret)
}

func (s *service) GetSecret(ctx context.Context, id, ownerID uuid.UUID) (*models.Secret, error) {
	return s.secretProvider.GetSecret(ctx, id, ownerID)
}

func (s *service) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	return s.secretProvider.UpdateSecret(ctx, secret)
}

func (s *service) DeleteSecret(ctx context.Context, id, ownerID uuid.UUID) error {
	return s.secretProvider.DeleteSecret(ctx, id, ownerID)
}

func (s *service) ListSecret(ctx context.Context, ownerID uuid.UUID) ([]*models.Secret, error) {
	return s.secretProvider.ListSecret(ctx, ownerID)
}

func (s *service) GetLastSync(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	return s.syncProvider.GetLastSync(ctx, userID)
}

func (s *service) ListUpdatedSecrets(ctx context.Context, ownerID uuid.UUID, updatedAfter *time.Time) ([]*models.Secret, error) {
	return s.secretProvider.ListUpdatedSecrets(ctx, ownerID, updatedAfter)
}

func (s *service) BatchUpdateSecrets(ctx context.Context, userID uuid.UUID, secrets []*models.Secret) error {
	return s.secretProvider.BatchUpdateSecrets(ctx, userID, secrets)
}

func (s *service) UpdateLastSync(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	return s.syncProvider.UpdateLastSync(ctx, userID)
}

func New(storage services.KeeperProvider) *service {
	return &service{
		profileProvider: storage,
		clientProvider:  storage,
		secretProvider:  storage,
		syncProvider:    storage,
	}
}
