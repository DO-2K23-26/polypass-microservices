package main

import "time"

// CredentialCreated represents an event when a credential is created
type CredentialCreated struct {
	UserID       string    `json:"user_id"`
	GroupID      string    `json:"group_id"`
	CredentialID string    `json:"credential_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// CredentialUpdated represents an event when a credential is updated
type CredentialUpdated struct {
	UserID       string    `json:"user_id"`
	GroupID      string    `json:"group_id"`
	CredentialID string    `json:"credential_id"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CredentialDeleted represents an event when a credential is deleted
type CredentialDeleted struct {
	UserID       string    `json:"user_id"`
	GroupID      string    `json:"group_id"`
	CredentialID string    `json:"credential_id"`
	DeletedAt    time.Time `json:"deleted_at"`
}

// CredentialShared represents an event when a credential is shared
type CredentialShared struct {
	UserID       string    `json:"user_id"`
	GroupID      string    `json:"group_id"`
	CredentialID string    `json:"credential_id"`
	SharedAt     time.Time `json:"shared_at"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	IsOneTime    bool      `json:"is_one_time"`
}
