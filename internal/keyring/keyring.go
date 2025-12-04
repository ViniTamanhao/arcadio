// Package keyring defines keyring methods
package keyring

import (
	"fmt"

	"github.com/zalando/go-keyring"
)


const (
	serviceName = "arcadio"
)

type KeyStore struct {}

func NewKeyStore() *KeyStore {
	return &KeyStore{}
}

func (ks *KeyStore) SavePassword(arcID, password string) error {
	return keyring.Set(serviceName, arcID, password)
}

func (ks *KeyStore) GetPassword(arcID string) (string, error) {
	password, err := keyring.Get(serviceName, arcID)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", fmt.Errorf("password not found in keyring")
		}
		return "", fmt.Errorf("failed to get password: %w", err)
	}
	return password, nil
}

func (ks *KeyStore) DeletePassword(arcID string) error {
	err := keyring.Delete(serviceName, arcID)
	if err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete password: %w", err)
	}
	return nil
}

func (ks *KeyStore) HasPassword(arcID string) bool {
	_, err := keyring.Get(serviceName, arcID)
	return err == nil
}

func (ks *KeyStore) ListStored() ([]string, error) {
	return nil, fmt.Errorf("list not supported by keyring backend")
}
