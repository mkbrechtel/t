package core

import (
	"time"
)

// Event interface defines the common methods all events must implement
type Event interface {
	// Apply changes the application state based on the event
	apply(*AppState) error
	// GetTimestamp returns when the event occurred
	GetTimestamp() time.Time
	// GetType returns the type of the event for serialization
	GetType() string
}

// BaseEvent provides common fields for all events
type BaseEvent struct {
	Timestamp time.Time
}

// GetTimestamp returns when the event occurred
func (e BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}
