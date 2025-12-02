package arc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ViniTamanhao/arcadio/internal/crypto"
	"github.com/ViniTamanhao/arcadio/pkg/models"
	"github.com/google/uuid"
)

// AddDocument adds a file to an arc
func (m *Manager) AddDocument (arcID string, arc *models.Arc, key []byte, filePath string, tags []string) (*models.Document, error) {
	fmt.Printf("Reading file: %s\n", filePath)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	fmt.Println("Calculating content hash...")
	hash := sha256.Sum256(fileData)
	contentHash := hex.EncodeToString(hash[:])

	// TODO Optional: Compress the data for storage efficiency
	dataToEncrypt := fileData
	compressed := false

	fmt.Println("Encrypting document...")
	encryptedData, err := crypto.Encrypt(key, dataToEncrypt)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt document: %w", err)
	}

	doc := &models.Document{
		ID: uuid.New().String(),
		Filename: filepath.Base(filePath),
		AddedAt: time.Now(),
		ModifiedAt: time.Now(),
		Size: int64(len(fileData)),
		ContentHash: contentHash,
		Compressed: compressed,
	}

	docPath := m.GetDocumentPath(arcID, doc.ID)
	fmt.Printf("Saving encrypted document: %s\n", docPath)

	if err := os.WriteFile(docPath, encryptedData, 0600); err != nil {
		return nil, fmt.Errorf("failed to save encrypted document: %w", err)
	}

	arc.Documents[doc.ID] = doc

	if len(tags) > 0 {
		arc.Tags[doc.ID] = tags
		fmt.Printf("Added tags %v\n", tags)
	}

	if err := m.Update(arcID, arc, key); err != nil {
		return nil, fmt.Errorf("failed to update arc metadata: %w", err)
	}

	fmt.Println("Document added successfully")
	return doc, nil
}

// RemoveDocument removes a document from an arc
func (m *Manager) RemoveDocument(arcID string, arc *models.Arc, key []byte, docID string) error {
	doc, exists := arc.Documents[docID]
 	if !exists {
		return fmt.Errorf("document not foun: %s", docID)
	} 

	fmt.Printf("Removing document: %s\n", doc.Filename)

	docPath := m.GetDocumentPath(arcID, docID)
	if err := os.Remove(docPath); err != nil {
		return fmt.Errorf("failed to delete document file: %w", err)
	}

	delete(arc.Documents, docID)
	delete(arc.Tags, docID)

	if err := m.Update(arcID, arc, key); err != nil {
		return fmt.Errorf("failed to update arc metadata: %w", err)
	}

	fmt.Println("Document removed successfully")
	return nil
}

func (m *Manager) ExportDocument(arcID string, arc *models.Arc, key []byte, docID string, outputPath string) error {
	doc, exists := arc.Documents[docID]
	if !exists {
		return fmt.Errorf("document not found: %s", docID)
	}

	fmt.Printf("Exporting document: %s\n", doc.Filename)

	docPath := m.GetDocumentPath(arcID, docID)
	encryptedData, err := os.ReadFile(docPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted document: %w", err)
	}

	fmt.Println("Decrypting document...")
	decryptedData, err := crypto.Decrypt(key, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to decrypt document: %w", err)
	}

	hash := sha256.Sum256(decryptedData)
	contentHash := hex.EncodeToString(hash[:])

	if contentHash != doc.ContentHash {
		return fmt.Errorf("document integrity check failed - file may be corrupted")
	}

	fmt.Printf("Writing to: %s\n", outputPath)
	if err := os.WriteFile(outputPath, decryptedData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Println("Document exported successfully")
	return nil
}

// GetDocument retrieves and decrypts a document's content
func (m *Manager) GetDocument(arcID string, arc *models.Arc, key []byte, docID string) ([]byte, error) {
	_, exists := arc.Documents[docID]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", docID)
	}

	docPath := m.GetDocumentPath(arcID, docID)
	encryptedData, err := os.ReadFile(docPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted document: %w", err)
	}

	decryptedData, err := crypto.Decrypt(key, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt document: %w", err)
	}

	return decryptedData, nil
}

// AddTags adds tags to a document
func (m *Manager) AddTags(arcID string, arc *models.Arc, key []byte, docID string, tags []string) error {
	if _, exists := arc.Documents[docID]; !exists {
		return fmt.Errorf("document not found: %s", docID)
	}

	existingTags := arc.Tags[docID]
	if existingTags == nil {
		existingTags = []string{}
	}

	tagSet := make(map[string]bool)
	for _, tag := range existingTags {
		tagSet[tag] = true
	}
	for _, tag := range tags {
		tagSet[tag] = true
	}

	newTags := make([]string, 0,len(tagSet))
	for tag := range tagSet {
		newTags = append(newTags, tag)
	}

	arc.Tags[docID] = newTags
	return m.Update(arcID, arc, key)
}

// RemoveTags removes tags from a document
func (m *Manager) RemoveTags(arcID string, arc *models.Arc, key []byte, docID string, tags []string) error {
	if _, exists := arc.Documents[docID]; !exists {
		return fmt.Errorf("document not found: %s", docID)
	}

	existingTags := arc.Tags[docID]
	if existingTags == nil {
		return nil
	}

	tagSet := make(map[string]bool)
	for _, tag := range existingTags {
		tagSet[tag] = true
	}
	for _, tag := range tags {
		delete(tagSet, tag)
	}

	newTags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		newTags = append(newTags, tag)
	}

	arc.Tags[docID] = newTags
	return m.Update(arcID, arc, key)
}

// ListDocuments returns all documents in the arc
func (m *Manager) ListDocuments(arc *models.Arc) []*models.Document {
	docs := make([]*models.Document, 0, len(arc.Documents))
	for _, doc := range arc.Documents {
		docs = append(docs, doc)
	}
	return docs
}

// SearchDocuments performs fuzzy search on document filenames
func (m *Manager) SearchDocuments(arc *models.Arc, query string) []*models.Document {
	// Simple substring search for now
	// We can improve this with fuzzy matching library later
	var matches []*models.Document
	
	for _, doc := range arc.Documents {
		if contains(doc.Filename, query) {
			matches = append(matches, doc)
		}
	}
	
	return matches
}

// Helper function for case-insensitive substring matching
func contains(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return len(s) >= len(substr) && (s == substr || indexof(s, substr) >= 0)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c + 32
		}
		result[i] = c
	}
	return string(result)
}

func indexof(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// AddDocumentFromReader adds a document from an io.Reader
func (m *Manager) AddDocumentFromReader(arcID string, arc *models.Arc, key []byte, filename string, reader io.Reader, tags []string) (*models.Document, error) {
	fileData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	hash := sha256.Sum256(fileData)
	contentHash := hex.EncodeToString(hash[:])

	encryptedData, err := crypto.Encrypt(key, fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt document: %w", err)
	}

	doc := &models.Document{
		ID:          uuid.New().String(),
		Filename:    filename,
		AddedAt:     time.Now(),
		ModifiedAt:  time.Now(),
		Size:        int64(len(fileData)),
		ContentHash: contentHash,
		Compressed:  false,
	}

	docPath := m.GetDocumentPath(arcID, doc.ID)
	if err := os.WriteFile(docPath, encryptedData, 0600); err != nil {
		return nil, fmt.Errorf("failed to save encrypted document: %w", err)
	}

	arc.Documents[doc.ID] = doc
	if len(tags) > 0 {
		arc.Tags[doc.ID] = tags
	}

	arc.ModifiedAt = time.Now()
	if err := m.Update(arcID, arc, key); err != nil {
		return nil, fmt.Errorf("failed to update arc metadata: %w", err)
	}

	return doc, nil
}
