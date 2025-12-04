// Package auth deals with auth logic
package auth

import (
	"fmt"
	"syscall"
	"time"

	"github.com/ViniTamanhao/arcadio/internal/keyring"
	"golang.org/x/term"
)

type Manager struct {
	keyStore *keyring.KeyStore
	sessionCache *keyring.SessionCache
}

func NewManager(cacheDir string) *Manager {
	return &Manager{
		keyStore: keyring.NewKeyStore(),
		sessionCache: keyring.NewSessionCache(cacheDir, 15*time.Minute),
	}
}

func (m *Manager) GetPassword(arcID, arcName string, allowPrompt bool) (string, error) {
	if password, ok := m.sessionCache.Get(arcID); ok {
		return password, nil
	}

	if password, err := m.keyStore.GetPassword(arcID); err == nil {
		m.sessionCache.Set(arcID, password)
		return password, nil
	}

	if !allowPrompt {
		return "", fmt.Errorf("password not found and prompting disabled")
	}

	password, err := m.promptPassword(arcName)
	if err != nil {
		return "", err
	}

	if m.askSavePassword() {
		if err := m.keyStore.SavePassword(arcID, password); err != nil {
			fmt.Printf("Warning: Failed to save password to keyring: %v\n", err)
		} else {
			fmt.Println("Password saved to system keyring")
		}
	}

	m.sessionCache.Set(arcID, password)

	return password, nil
}

func (m *Manager) SavePassword(arcID, password string) error {
	return m.keyStore.SavePassword(arcID, password)
}

func (m *Manager) DeletePassword(arcID string) error {
	m.sessionCache.Clear(arcID)
	return m.keyStore.DeletePassword(arcID)
}

func (m *Manager) HasStoredPassword(arcID string) bool {
	return m.keyStore.HasPassword(arcID)
}

func (m *Manager) ClearSession() {
	m.sessionCache.ClearAll()
}

func (m *Manager) promptPassword(arcName string) (string, error) {
	fmt.Printf("Unlocking arc: %s\n", arcName)
	fmt.Print("Enter password: ")
	
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()

	return string(passwordBytes), nil
}

// askSavePassword asks user if they want to save the password
func (m *Manager) askSavePassword() bool {
	fmt.Print("Save password to system keyring? [y/N]: ")
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y" || response == "yes"
}
