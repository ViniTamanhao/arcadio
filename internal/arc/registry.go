package arc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)
 
type ArcEntry struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Registry struct {
	configPath string
	Arcs       map[string]*ArcEntry `json:"arcs"` // ID -> ArcEntry
}

// NewRegistry creates or loads the arc registry
func NewRegistry(baseDir string) (*Registry, error) {
	configPath := filepath.Join(baseDir, "..", "registry.json")
	
	registry := &Registry{
		configPath: configPath,
		Arcs:       make(map[string]*ArcEntry),
	}

	if err := registry.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load registry: %w", err)
		}
	}

	return registry, nil
}

// Register adds or updates an arc entry
func (r *Registry) Register(id, name string, createdAt time.Time) error {
	r.Arcs[id] = &ArcEntry{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
	}
	return r.save()
}

// load reads the registry from disk
func (r *Registry) load() error {
	data, err := os.ReadFile(r.configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.Arcs)
}

// Unregister removes an arc entry
func (r *Registry) Unregister(id string) error {
	delete(r.Arcs, id)
	return r.save()
}

// GetByID returns arc entry by ID
func (r *Registry) GetByID(id string) (*ArcEntry, bool) {
	entry, exists := r.Arcs[id]
	return entry, exists
}

// GetByName returns arc entry by name
func (r *Registry) GetByName(name string) (*ArcEntry, bool) {
	for _, entry := range r.Arcs {
		if entry.Name == name {
			return entry, true
		}
	}
	return nil, false
}

// ListAll returns all arc entries
func (r *Registry) ListAll() []*ArcEntry {
	entries := make([]*ArcEntry, 0, len(r.Arcs))
	for _, entry := range r.Arcs {
		entries = append(entries, entry)
	}
	return entries
}

// FindArc tries to find an arc by ID or name
func (r *Registry) FindArc(idOrName string) (*ArcEntry, error) {
	if entry, exists := r.GetByID(idOrName); exists {
		return entry, nil
	}

	if entry, exists := r.GetByName(idOrName); exists {
		return entry, nil
	}

	for id, entry := range r.Arcs {
		if len(idOrName) >= 8 && id[:len(idOrName)] == idOrName {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("arc not found: %s", idOrName)
}

// save writes the registry to disk
func (r *Registry) save() error {
	dir := filepath.Dir(r.configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(r.Arcs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.configPath, data, 0600)
}
