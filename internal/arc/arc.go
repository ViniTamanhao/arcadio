// Package arc defines all the methods and managers for arcs
package arc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ViniTamanhao/arcadio/internal/crypto"
	"github.com/ViniTamanhao/arcadio/pkg/models"
	"github.com/google/uuid"
)

type Manager struct {
	baseDir string
	registry *Registry
}

// NewManager creates a NewManager instance
func NewManager(baseDir string) (*Manager, error) {
	registry, err := NewRegistry(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load registry: %w", err)
	}

	return &Manager{baseDir: baseDir, registry: registry}, nil
}

// Create creates a new arc
func (m *Manager) Create(name, password, securityQuestion, securityAnswer string) (*models.Arc, error) {
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key := crypto.DeriveKey(password, salt)
	passwordHash := crypto.HashAnswer(password)
	answerHash := crypto.HashAnswer(securityAnswer)

	arc := &models.Arc{
		ID: uuid.New().String(),
		Name: name,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
		Documents: make(map[string]*models.Document),
		Tags: make(map[string][]string),
		EncryptionVersion: "v1",
	}

	secConfig := &models.SecurityConfig{
		Salt: salt,
		PasswordHash: passwordHash,
		SecurityQuestion: securityQuestion,
		AnswerHash: answerHash,
		KeyDerivation: "argon2id",
	}

	arcDir := filepath.Join(m.baseDir, arc.ID)
	if err := os.MkdirAll(arcDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create arc directory: %w", err)
	}

	docsDir := filepath.Join(arcDir, "documents")
	if err := os.MkdirAll(docsDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create documents directory: %w", err)
	}

	if err := m.saveSecurityConfig(arcDir, secConfig); err != nil {
		return nil, fmt.Errorf("failed to save security confi: %w", err)
	}

	if err := m.saveArcMetadata(arcDir, arc, key); err != nil {
		return nil, fmt.Errorf("failed to save arc metadata: %w", err)
	}

	fmt.Println("Registering arc...")
	if err := m.registry.Register(arc.ID, arc.Name, arc.CreatedAt); err != nil {
		return nil,fmt.Errorf("failed to register arc: %w", err)
	}

	return arc, nil
}

// Delete removes an arc completely
func (m *Manager) Delete(idOrName string) error {
	entry, err := m.registry.FindArc(idOrName)
	if err != nil {
		return err
	}

	arcDir := filepath.Join(m.baseDir, entry.ID)
	if err := os.RemoveAll(arcDir); err != nil {
		return fmt.Errorf("failed to delete arc directory: %w", err)
	}

	if err := m.registry.Unregister(entry.ID); err != nil {
		return fmt.Errorf("failed to unregister arc: %w", err)
	}

	return nil
}

// FindArc returns arc entry by ID or name
func (m *Manager) FindArc(idOrName string) (*ArcEntry, error) {
	return m.registry.FindArc(idOrName)
}

// ListArcs returns all registered arcs
func (m *Manager) ListArcs() []*ArcEntry {
	return m.registry.ListAll()
}

// Unlock verifies password and loads arc into memory
func (m *Manager) Unlock(idOrName, password string) (*models.Arc, []byte, error) {
	entry, err := m.registry.FindArc(idOrName)
	if err != nil {
		return nil, nil, err
	}

	arcDir := filepath.Join(m.baseDir, entry.ID)

	fmt.Println("Loading security configuration...")
	secConfig, err := m.loadSecurityConfig(arcDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load security config: %w", err)
	}

	// Derive key from provided password
	fmt.Println("Deriving encryption key...")
	key := crypto.DeriveKey(password, secConfig.Salt)
	
	// Verify password by comparing hashes
	fmt.Println("Verifying password...")
	passwordHash := crypto.HashAnswer(password)
	if !bytesEqual(passwordHash, secConfig.PasswordHash) {
		return nil, nil, fmt.Errorf("invalid password")
	}

	// Load and decrypt arc metadata
	fmt.Println("Decrypting arc metadata...")
	arc, err := m.loadArcMetadata(arcDir, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load arc: %w", err)
	}

	fmt.Println("Arc unlocked successfully")
	return arc, key, nil
}

// GetDocumentPath gets the path of a certain document
func (m *Manager) GetDocumentPath(arcID, docID string) string {
	return filepath.Join(m.baseDir, arcID, "documents", docID+".bin")
}

// saveSecurityConfig saves the security config of a certain arc
func (m *Manager) saveSecurityConfig(arcDir string, config *models.SecurityConfig) error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	path := filepath.Join(arcDir, "arc.sec")
	return os.WriteFile(path, data, 0600)
}

// loadSecurityConfig loads the security config of a certain arc
func (m *Manager) loadSecurityConfig(arcDir string) (*models.SecurityConfig, error) {
	path := filepath.Join(arcDir, "arc.sec")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config models.SecurityConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// saveArcMetadata saves the metadata of a certain arc
func (m *Manager) saveArcMetadata(arcDir string, arc *models.Arc, key []byte) error {
	data, err := json.MarshalIndent(arc, "", " ")
	if err != nil {
		return err
	}

	encrypted, err := crypto.Encrypt(key, data)
	if err != nil {
		return err
	}

	path := filepath.Join(arcDir, "arc.meta")
	return os.WriteFile(path, encrypted, 0600)
}

// loadArcMetadata loads the metadata of a certain arc
func (m *Manager) loadArcMetadata(arcDir string, key []byte) (*models.Arc, error) {
	path := filepath.Join(arcDir, "arc.meta")
	encrypted, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data, err := crypto.Decrypt(key, encrypted)
	if err != nil {
		return nil, err
	}

	var arc models.Arc
	if err := json.Unmarshal(data, &arc); err != nil {
		return nil, err
	}

	return &arc, nil
}

// Update updates the metadata of an arc
func (m *Manager) Update(arcID string, arc *models.Arc, key []byte) error {
	arcDir := filepath.Join(m.baseDir, arcID)
	arc.ModifiedAt = time.Now()
	return m.saveArcMetadata(arcDir, arc, key)
}

// bytesEqual compares two bytes
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
