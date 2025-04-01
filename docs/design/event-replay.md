# Event Replay System

## Overview

The t5, task manager implements an event sourcing pattern for its core functionality. This document explains how the event replay system works and how to use it in your own code.

## Core Concepts

### Event Sourcing

Event sourcing is a pattern where all changes to application state are captured as a sequence of events. Instead of storing the current state, the system stores a complete history of events that led to the current state. The application state can be reconstructed at any point by replaying these events.

Key advantages:
- Complete audit trail of all changes
- Ability to reconstruct the state at any point in time
- Natural support for event-driven architectures
- Separation of write and read models

### Events

Events in t5 are immutable records of something that happened. Each event:
- Has a timestamp when it occurred
- Has a type identifier
- Contains all data necessary to apply the change
- Has an `apply` method that modifies the application state

The base `Event` interface requires:
```go
type Event interface {
    apply(*AppState) error
    GetTimestamp() time.Time
    GetType() string
}
```

All events should embed the `BaseEvent` struct to get common functionality.

### AppState

The `AppState` is the current snapshot of the application. It contains:
- Tasks (indexed by UUID)
- Projects (indexed by name)

The state is built by applying events in chronological order.

## Implemented Events

### TodoTxtTaskUpdate

This event represents tasks imported from or updated via todo.txt format. When applied:
1. Parses the todo.txt lines into tasks
2. Assigns UUIDs if needed
3. Updates the application state with tasks
4. Creates projects referenced by the tasks

Usage:
```go
// Create a new TodoTxtTaskUpdate event
event := NewTodoTxtTaskUpdate("(A) 2023-07-15 Implement feature +project @context")

// The event can be saved to an EventStore
store.SaveEvent(event)

// Or applied directly to the state
state := NewAppState()
ApplyEvent(state, event)
```

## Repository and Event Storage

The `Repository` is a concrete in-memory implementation that manages both event storage and the application state:

```go
type Repository struct {
    events []Event
    state  *AppState
    mu     sync.RWMutex
}
```

The repository handles:
- Storing events
- Maintaining the current application state
- Thread-safety with a mutex
- Rebuilding the state from events when needed
- Providing access to tasks and projects

To use the repository:

```go
// Create a new repository
repo := NewRepository()

// Add events that update the state
event := NewTodoTxtTaskUpdate("Task description +project @context")
repo.SaveEvent(event)

// Access the state
state := repo.GetAppState()
for id, task := range state.Tasks {
    // Process tasks
}

// Rebuild the state from events if needed
repo.RebuildState()
```

The repository provides these high-level operations:

1. **Task Management**:
   - `CreateTask(task Task) (uuid.UUID, error)`
   - `GetTask(id uuid.UUID) (Task, error)`
   - `UpdateTask(task Task) error`
   - `DeleteTask(id uuid.UUID) error`
   - `ListTasks(filters []TaskFilter) ([]Task, error)`
   - `SearchTasks(query string) ([]Task, error)`

2. **Project Management**:
   - `CreateProject(project Project) (uuid.UUID, error)`
   - `GetProject(id uuid.UUID) (Project, error)`
   - `UpdateProject(project Project) error`
   - `DeleteProject(id uuid.UUID) error`
   - `ListProjects() ([]Project, error)`
   - `GetTasksForProject(projectID uuid.UUID) ([]Task, error)`

3. **Event Management**:
   - `SaveEvent(event Event) error`
   - `GetEvents() ([]Event, error)`
   - `RecordTodoTxtTaskUpdate(update *TodoTxtTaskUpdate) error`
   - `RecordTaskStartTime(startTime *TaskStartTime) error`
   - `RecordTaskEndTime(endTime *TaskEndTime) error`

4. **State Management**:
   - `GetAppState() *AppState`
   - `RebuildState() error`

All state-changing operations are implemented by creating and storing events, ensuring the event sourcing pattern is followed.

Note that while a traditional event-sourced system would typically use an interface to abstract the storage, this implementation uses a concrete in-memory store for simplicity.

## Adding New Event Types

To add a new event type:

1. Define a new struct that embeds `BaseEvent`
2. Implement the `GetType()` method to return the event type name
3. Implement the `apply(*AppState) error` method to modify the state
4. Add a constructor function that sets the timestamp

Example:

```go
// Define the event
type TaskCompleted struct {
    BaseEvent
    TaskID uuid.UUID
}

// Implement the GetType method
func (e *TaskCompleted) GetType() string {
    return "TaskCompleted"
}

// Implement the apply method
func (e *TaskCompleted) apply(state *AppState) error {
    task, ok := state.Tasks[e.TaskID]
    if !ok {
        return fmt.Errorf("task not found: %s", e.TaskID)
    }
    task.Completed = true
    task.CompletedDate = e.Timestamp
    state.Tasks[e.TaskID] = task
    return nil
}

// Add a constructor
func NewTaskCompleted(taskID uuid.UUID) *TaskCompleted {
    event := &TaskCompleted{
        TaskID: taskID,
    }
    event.Timestamp = now()
    return event
}
```

## Testing Events

When testing events, you can:
1. Override the `now` function to provide deterministic timestamps
2. Use a `MockEventStore` to store and retrieve events for testing
3. Create test events and apply them to a state
4. Verify the state reflects the expected changes

## Future Enhancements

Planned enhancements to the event system:
- Event versioning for schema evolution
- Event filtering for efficient queries
- Snapshots for faster state reconstruction
- Event subscription for real-time updates