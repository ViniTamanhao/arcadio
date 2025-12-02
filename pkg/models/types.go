// Package models defines the types for arcadio
package models

import "time"


type Arc struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	CreatedAt         time.Time              `json:"created_at"`
	ModifiedAt        time.Time              `json:"modified_at"`
	Documents         map[string]*Document   `json:"documents"`
	Tags              map[string][]string    `json:"tags"` // doc_id -> tags
	EncryptionVersion string                 `json:"encryption_version"`
}

type Document struct {
	ID          string    `json:"id"`
	Filename    string    `json:"filename"`
	AddedAt     time.Time `json:"added_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	Size        int64     `json:"size"`
	ContentHash string    `json:"content_hash"` // SHA-256
	Compressed  bool      `json:"compressed"`
}

type SecurityConfig struct {
	Salt             []byte `json:"salt"`
	PasswordHash     []byte `json:"password_hash"`
	SecurityQuestion string `json:"security_question"`
	AnswerHash       []byte `json:"answer_hash"`
	KeyDerivation    string `json:"key_derivation"` // "argon2id"
}


