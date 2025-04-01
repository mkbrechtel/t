package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

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

// EventData represents the serialized form of an event
type EventData struct {
	Type      string          `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// AppendLogEventStore implements EventStore using a JSONL file
type AppendLogEventStore struct {
	filePath string
	mu       sync.RWMutex
}

// NewAppendLogEventStore creates a new append-only log store using the specified file
func NewAppendLogEventStore(filePath string) (*AppendLogEventStore, error) {
	// Create the file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create event log file: %w", err)
		}
		file.Close()
	}

	return &AppendLogEventStore{
		filePath: filePath,
	}, nil
}

// SaveEvent persists an event to the append-only log file
func (s *AppendLogEventStore) SaveEvent(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Open the file in append mode
	file, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open event log file: %w", err)
	}
	defer file.Close()

	// Create the EventData wrapper
	wrapper := EventData{
		Type:      event.GetType(),
		Timestamp: event.GetTimestamp(),
		Data:      nil, // Will be filled below
	}

	// Marshal the event based on its type directly into the wrapper's Data field
	switch e := event.(type) {
	case *TodoTxtTaskUpdate:
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal TodoTxtTaskUpdate event: %w", err)
		}
		wrapper.Data = json.RawMessage(data)
	case *TaskStartTime:
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal TaskStartTime event: %w", err)
		}
		wrapper.Data = json.RawMessage(data)
	case *TaskEndTime:
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal TaskEndTime event: %w", err)
		}
		wrapper.Data = json.RawMessage(data)
	case *SetTimeBudgetForProject:
		data, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal SetTimeBudgetForProject event: %w", err)
		}
		wrapper.Data = json.RawMessage(data)
	default:
		return fmt.Errorf("unsupported event type: %s", event.GetType())
	}

	// Marshal the wrapper
	data, err := json.Marshal(wrapper)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Write to the file
	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write event to log: %w", err)
	}

	return nil
}

// GetEvents retrieves all events from the append-only log file
func (s *AppendLogEventStore) GetEvents() ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Open the file for reading
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open event log file: %w", err)
	}
	defer file.Close()

	var events []Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse the EventData wrapper
		var wrapper EventData
		if err := json.Unmarshal([]byte(line), &wrapper); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		// Create the appropriate event type based on the wrapper type
		var event Event
		switch wrapper.Type {
		case "TodoTxtTaskUpdate":
			e := &TodoTxtTaskUpdate{}
			if err := json.Unmarshal(wrapper.Data, e); err != nil {
				return nil, fmt.Errorf("failed to unmarshal TodoTxtTaskUpdate event: %w", err)
			}
			e.Timestamp = wrapper.Timestamp
			event = e
		case "TaskStartTime":
			e := &TaskStartTime{}
			if err := json.Unmarshal(wrapper.Data, e); err != nil {
				return nil, fmt.Errorf("failed to unmarshal TaskStartTime event: %w", err)
			}
			e.Timestamp = wrapper.Timestamp
			event = e
		case "TaskEndTime":
			e := &TaskEndTime{}
			if err := json.Unmarshal(wrapper.Data, e); err != nil {
				return nil, fmt.Errorf("failed to unmarshal TaskEndTime event: %w", err)
			}
			e.Timestamp = wrapper.Timestamp
			event = e
		case "SetTimeBudgetForProject":
			e := &SetTimeBudgetForProject{}
			if err := json.Unmarshal(wrapper.Data, e); err != nil {
				return nil, fmt.Errorf("failed to unmarshal SetTimeBudgetForProject event: %w", err)
			}
			e.Timestamp = wrapper.Timestamp
			event = e
		default:
			return nil, fmt.Errorf("unknown event type: %s", wrapper.Type)
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading event log file: %w", err)
	}

	return events, nil
}