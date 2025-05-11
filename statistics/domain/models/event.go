package models

import (
	"time"
)

// Event represents a base event in the system
type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
	Data      interface{} `json:"data"`
}

// EventType defines the possible types of events
type EventType string

const (
	// Define common event types here
	UserLoginEvent        EventType = "USER_LOGIN"
	UserLogoutEvent       EventType = "USER_LOGOUT"
	PasswordCreatedEvent  EventType = "PASSWORD_CREATED"
	PasswordUpdatedEvent  EventType = "PASSWORD_UPDATED"
	PasswordDeletedEvent  EventType = "PASSWORD_DELETED"
	PasswordAccessedEvent EventType = "PASSWORD_ACCESSED"
	// Add more event types as needed
)

// EventSource defines the possible sources of events
type EventSource string

const (
	// Define common event sources here
	WebApplication EventSource = "WEB_APP"
	MobileApp      EventSource = "MOBILE_APP"
	APIClient      EventSource = "API_CLIENT"
	// Add more sources as needed
)
