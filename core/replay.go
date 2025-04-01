package core

import (
	"fmt"
	"sort"
)

// EventStore interface defines methods for storing and retrieving events
type EventStore interface {
	// SaveEvent persists an event to storage
	SaveEvent(event Event) error
	
	// GetEvents retrieves all events from storage, optionally filtered
	GetEvents() ([]Event, error)
}

// ReplayEvents applies all events from the event store to rebuild application state
func ReplayEvents(store EventStore) (*AppState, error) {
	// Create a new empty state
	state := NewAppState()
	
	// Get all events from the store
	events, err := store.GetEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve events: %w", err)
	}
	
	// Sort events by timestamp to ensure proper order
	sort.Slice(events, func(i, j int) bool {
		return events[i].GetTimestamp().Before(events[j].GetTimestamp())
	})
	
	// Apply events in order to build up the state
	for _, event := range events {
		if err := event.apply(state); err != nil {
			return nil, fmt.Errorf("failed to apply event %s: %w", event.GetType(), err)
		}
	}
	
	return state, nil
}

// ApplyEvent applies a single event to the application state
func ApplyEvent(state *AppState, event Event) error {
	return event.apply(state)
}

// NewTodoTxtTaskUpdate creates a new TodoTxtTaskUpdate event with current timestamp
func NewTodoTxtTaskUpdate(lines string) *TodoTxtTaskUpdate {
	event := &TodoTxtTaskUpdate{
		Lines: lines,
	}
	event.Timestamp = now()
	return event
}