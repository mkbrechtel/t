package core

import uuid "github.com/gofrs/uuid/v5"

// AppState represents the current state of the application
// It is built by replaying all events from the beginning
type AppState struct {
	Tasks    map[uuid.UUID]Task
	Projects map[string]Project
}

// NewAppState creates a new empty application state
func NewAppState() *AppState {
	return &AppState{
		Tasks:    make(map[uuid.UUID]Task),
		Projects: make(map[string]Project),
	}
}
