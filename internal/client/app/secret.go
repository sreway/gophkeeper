package app

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
	"github.com/sreway/gophkeeper/internal/lib/crypt"
)

func (c *client) entryToSecret(userID uuid.UUID, entry models.Entry) (*models.Secret, error) {
	var err error
	secret := new(models.Secret)

	secret.ID = entry.ID()

	entryBytes, err := entry.Bytes()
	if err != nil {
		return nil, err
	}

	secret.EncryptedValue, err = c.keeper.EncryptData(entryBytes, c.masterKey)
	if err != nil {
		return nil, err
	}

	secret.Owner = userID
	secret.Type = entry.Type()
	secret.Hash = crypt.HashSum(entryBytes, nil)

	return secret, nil
}

func (c *client) secretToEntry(secret *models.Secret) (models.Entry, error) {
	var (
		entry models.Entry
		err   error
	)

	entryBytes, err := c.keeper.DecryptData(secret.EncryptedValue, c.masterKey)
	if err != nil {
		return nil, err
	}

	switch secret.Type {
	case models.CredentialsType:
		entry = new(models.Credentials)
	case models.TextType:
		entry = new(models.Text)
	default:
		return nil, models.ErrUnknownSecretType
	}

	err = json.Unmarshal(entryBytes, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}
