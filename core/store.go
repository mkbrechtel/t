package core

import "sync"

// InMemoryEventStore is a simple in-memory implementation of EventStore
type InMemoryEventStore struct {
	events []Event
	mu     sync.RWMutex
}

// NewInMemoryEventStore creates a new in-memory event store
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: []Event{},
	}
}

// SaveEvent persists an event to the in-memory store
func (s *InMemoryEventStore) SaveEvent(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the event
	s.events = append(s.events, event)
	return nil
}

// GetEvents retrieves all events from the in-memory store
func (s *InMemoryEventStore) GetEvents() ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy of the events slice to prevent modification
	eventsCopy := make([]Event, len(s.events))
	copy(eventsCopy, s.events)
	return eventsCopy, nil
}
