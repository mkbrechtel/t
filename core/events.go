package core

import (
	"fmt"
	"sort"
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
func NewTodoTxtTaskUpdate(lines string, source string) *TodoTxtTaskUpdate {
	event := &TodoTxtTaskUpdate{
		Lines: lines,
		Source: source,
	}
	event.Timestamp = now()
	return event
}
