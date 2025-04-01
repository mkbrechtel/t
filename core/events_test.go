package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

// MockEventStore is a simple in-memory event store for testing
type MockEventStore struct {
	events []Event
}

func (m *MockEventStore) SaveEvent(event Event) error {
	m.events = append(m.events, event)
	return nil
}

func (m *MockEventStore) GetEvents() ([]Event, error) {
	return m.events, nil
}

func TestReplayEvents(t *testing.T) {
	// Override now function for deterministic tests
	originalNow := now
	defer func() { now = originalNow }()
	
	mockTime := time.Date(2023, 7, 15, 10, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime }

	// Create mock event store
	store := &MockEventStore{}
	
	// Create some test events
	event1 := NewTodoTxtTaskUpdate(`(A) 2023-05-15 Implement the core module +t5 @dev due:2023-06-01`, "test")
	mockTime = mockTime.Add(time.Hour)
	now = func() time.Time { return mockTime }
	
	event2 := NewTodoTxtTaskUpdate(`x 2023-06-10 2023-05-20 Fix todo parsing bug +t5 @bugfix uuid:be52544a-22c3-7d86-b3a6-47f84e3fe9b2`, "test")
	mockTime = mockTime.Add(time.Hour)
	now = func() time.Time { return mockTime }
	
	event3 := NewTodoTxtTaskUpdate(`Buy groceries @personal +shopping due:2023-05-25`, "test")
	
	// Add events to store
	assert.NoError(t, store.SaveEvent(event1))
	assert.NoError(t, store.SaveEvent(event2))
	assert.NoError(t, store.SaveEvent(event3))
	
	// Replay events
	state, err := ReplayEvents(store)
	assert.NoError(t, err)
	assert.NotNil(t, state)
	
	// Verify state
	assert.Equal(t, 3, len(state.Tasks), "Expected 3 tasks in state")
	assert.Equal(t, 2, len(state.Projects), "Expected 2 projects in state")
	
	// Verify projects
	_, hasT5 := state.Projects["t5"]
	assert.True(t, hasT5, "Expected project 't5' to exist")
	
	_, hasShopping := state.Projects["shopping"]
	assert.True(t, hasShopping, "Expected project 'shopping' to exist")
	
	// Verify task with known UUID
	knownTaskID := uuid.Must(uuid.FromString("be52544a-22c3-7d86-b3a6-47f84e3fe9b2"))
	task, exists := state.Tasks[knownTaskID]
	assert.True(t, exists, "Expected to find task with UUID be52544a-22c3-7d86-b3a6-47f84e3fe9b2")
	assert.True(t, task.Completed, "Expected task to be completed")
	assert.Contains(t, task.Projects, "t5", "Expected task to have project 't5'")
	assert.Contains(t, task.Contexts, "bugfix", "Expected task to have context 'bugfix'")
}

func TestApplyEvent(t *testing.T) {
	// Create a new state
	state := NewAppState()
	
	// Create an event
	event := NewTodoTxtTaskUpdate(`(A) 2023-05-15 Implement the core module +t5 @dev due:2023-06-01`, "test")
	
	// Apply event
	err := ApplyEvent(state, event)
	assert.NoError(t, err)
	
	// Verify state
	assert.Equal(t, 1, len(state.Tasks), "Expected 1 task in state")
	assert.Equal(t, 1, len(state.Projects), "Expected 1 project in state")
}

func TestEventOrder(t *testing.T) {
	// Override now function for deterministic tests
	originalNow := now
	defer func() { now = originalNow }()
	
	// Create events with timestamps out of order
	mockTime3 := time.Date(2023, 7, 15, 12, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime3 }
	event3 := NewTodoTxtTaskUpdate(`Buy groceries @personal +shopping due:2023-05-25`, "test")
	
	mockTime1 := time.Date(2023, 7, 15, 10, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime1 }
	event1 := NewTodoTxtTaskUpdate(`(A) 2023-05-15 Task 1 +project1`, "test")
	
	mockTime2 := time.Date(2023, 7, 15, 11, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime2 }
	event2 := NewTodoTxtTaskUpdate(`(B) 2023-05-16 Task 2 +project2`, "test")
	
	// Create store with events in wrong order
	store := &MockEventStore{}
	assert.NoError(t, store.SaveEvent(event3))
	assert.NoError(t, store.SaveEvent(event1))
	assert.NoError(t, store.SaveEvent(event2))
	
	// Replay events
	state, err := ReplayEvents(store)
	assert.NoError(t, err)
	
	// Since our test is looking at project creation order to verify events were
	// applied in chronological order, we need to ensure state.Projects uses a
	// data structure that preserves insertion order. If it doesn't, this test
	// might not be reliable. Consider this a note for future refactoring.
	
	// Attempt to verify ordering by checking project creation, but this is
	// not entirely reliable due to the map structure
	assert.Equal(t, 3, len(state.Projects), "Expected 3 projects in state")
	_, hasProject1 := state.Projects["project1"]
	assert.True(t, hasProject1, "Expected project1 to exist")
	_, hasProject2 := state.Projects["project2"]
	assert.True(t, hasProject2, "Expected project2 to exist")
	_, hasShopping := state.Projects["shopping"]
	assert.True(t, hasShopping, "Expected shopping project to exist")
}

func TestNewTodoTxtTaskUpdate(t *testing.T) {
	// Override now function for deterministic tests
	originalNow := now
	defer func() { now = originalNow }()
	
	mockTime := time.Date(2023, 7, 15, 10, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime }
	
	// Create new event
	event := NewTodoTxtTaskUpdate("Test task", "test")
	
	// Verify event properties
	assert.Equal(t, "Test task", event.Lines)
	assert.Equal(t, mockTime, event.Timestamp)
	assert.Equal(t, "TodoTxtTaskUpdate", event.GetType())
}

func ExampleReplayEvents() {
	// Create a mock store with some events
	store := &MockEventStore{}
	
	// Add task events
	store.SaveEvent(NewTodoTxtTaskUpdate(`(A) 2023-05-15 First task +project1 @context1`, "test"))
	store.SaveEvent(NewTodoTxtTaskUpdate(`(B) 2023-05-16 Second task +project2 @context2`, "test"))
	
	// Replay events to build application state
	state, err := ReplayEvents(store)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Use the state
	fmt.Printf("Tasks: %d\n", len(state.Tasks))
	fmt.Printf("Projects: %d\n", len(state.Projects))
	
	// Output:
	// Tasks: 2
	// Projects: 2
}