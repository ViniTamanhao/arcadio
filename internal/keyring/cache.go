package keyring

import (
	"path/filepath"
	"time"
)

type SessionCache struct {
	cachePath string
	passwords map[string]*CachedPassword
	ttl 			time.Duration
}

type CachedPassword struct {
	ArcID 		string
	Password 	string
	ExpiresAt time.Time
}

func NewSessionCache(cacheDir string, ttl time.Duration) *SessionCache {
	return &SessionCache{
		cachePath: 	filepath.Join(cacheDir, ".session_cache"),
		passwords: 	make(map[string]*CachedPassword),
		ttl: 				ttl,
	}
}

// Set adds a password to the session cache
func (sc *SessionCache) Set(arcID, password string) {
	sc.passwords[arcID] = &CachedPassword{
		ArcID: arcID,
		Password: password,
		ExpiresAt: time.Now().Add(sc.ttl),
	}
}

// Get retrieves a password from the session cache
func (sc *SessionCache) Get(arcID string) (string, bool) {
	cached, exists := sc.passwords[arcID]
	if !exists {
		return "", false
	}

	if time.Now().After(cached.ExpiresAt) {
		delete(sc.passwords, arcID)
		return "", false
	}

	return cached.Password, true
}

// Clear removes a password from the session cache
func (sc *SessionCache) Clear(arcID string) {
	delete(sc.passwords, arcID)
}

// ClearAll removes all passwords from the session cache
func (sc *SessionCache) ClearAll() {
	sc.passwords = make(map[string]*CachedPassword)
}


