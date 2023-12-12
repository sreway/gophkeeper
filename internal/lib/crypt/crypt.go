package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, 10000, 32, sha256.New)
}

func EncryptGCM(data, password []byte) ([]byte, error) {
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	ciphertext = append(salt, ciphertext...) // nozero
	return ciphertext, nil
}

func DecryptGCM(ciphertext, password []byte) ([]byte, error) {
	salt := ciphertext[:8]
	ciphertext = ciphertext[8:]

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func GenerateMasterKey(keyLength int) ([]byte, error) {
	masterKey := make([]byte, keyLength)
	_, err := rand.Read(masterKey)
	if err != nil {
		return nil, err
	}

	return masterKey, nil
}

func HashSum(data []byte, salt []byte) string {
	hash := sha512.New()
	data = append(data, salt...)
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
