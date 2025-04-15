package models

import "time"

// UserMetrics represents the metrics for a user's credentials
type UserMetrics struct {
	UserID           string    `json:"user_id"`
	TotalCredentials int       `json:"total_credentials"`
	LastUpdated      time.Time `json:"last_updated"`
}

// GroupMetrics represents the metrics for a group's credentials
type GroupMetrics struct {
	GroupID          string    `json:"group_id"`
	TotalCredentials int       `json:"total_credentials"`
	ActiveUsers      int       `json:"active_users"`
	LastUpdated      time.Time `json:"last_updated"`
}

// CredentialTrend represents the trend of credential creation or usage
type CredentialTrend struct {
	StartDate  time.Time        `json:"start_date"`
	EndDate    time.Time        `json:"end_date"`
	DataPoints []TrendDataPoint `json:"data_points"`
}

// TrendDataPoint represents a single point in a trend
type TrendDataPoint struct {
	Date      time.Time `json:"date"`
	Creations int       `json:"creations"`
	Accesses  int       `json:"accesses"`
}

// SharedCredentialStats represents statistics about shared credentials
type SharedCredentialStats struct {
	CredentialID  string    `json:"credential_id"`
	TotalViews    int       `json:"total_views"`
	UniqueViewers int       `json:"unique_viewers"`
	FirstShared   time.Time `json:"first_shared"`
	LastAccessed  time.Time `json:"last_accessed"`
	OneTimeViews  int       `json:"one_time_views"`
}

// PasswordStrength represents the analysis of a password's strength
type PasswordStrength struct {
	CredentialID string `json:"credential_id"`
	Strength     string `json:"strength"`
	Score        int    `json:"score"`
}

// ReusedPassword represents a password that is used across multiple credentials
type ReusedPassword struct {
	Password      string   `json:"password"`
	CredentialIDs []string `json:"credential_ids"`
}

// BreachedCredential represents a credential that has been found in a data breach
type BreachedCredential struct {
	CredentialID string `json:"credential_id"`
	Password     string `json:"password"`
}

// OldPassword represents a password that hasn't been changed in a long time
type OldPassword struct {
	CredentialID string `json:"credential_id"`
	Password     string `json:"password"`
	Age          int    `json:"age"`
}

// CredentialAccess represents a record of a credential being accessed
type CredentialAccess struct {
	CredentialID string    `json:"credential_id"`
	UserID       string    `json:"user_id"`
	GroupID      string    `json:"group_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	Timestamp    time.Time `json:"timestamp"`
	IsOneTime    bool      `json:"is_one_time"`
}
