package core

import (
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

// ActiveTaskInfo contains information about the currently active task
type ActiveTaskInfo struct {
	TaskID      uuid.UUID
	StartTime   time.Time
	LastPaused  time.Time // Zero if not paused
	TotalPaused time.Duration
	EventID     uuid.UUID // ID of the TaskStartTime event
}

// TimeRecord represents a completed period of work on a task
type TimeRecord struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	TaskID    uuid.UUID
	Note      string
}

// ProjectBudget represents time budget information for a project
type ProjectBudget struct {
	ProjectID   uuid.UUID
	Budget      time.Duration
	UsedTime    time.Duration
	Description string
}

// AppState represents the current state of the application
// It is built by replaying all events from the beginning
type AppState struct {
	Tasks      map[uuid.UUID]Task
	Projects   map[string]Project
	ActiveTask *ActiveTaskInfo
	TimeRecords []TimeRecord
	ProjectBudgets map[uuid.UUID]ProjectBudget
}

// NewAppState creates a new empty application state
func NewAppState() *AppState {
	return &AppState{
		Tasks:         make(map[uuid.UUID]Task),
		Projects:      make(map[string]Project),
		ActiveTask:    nil,
		TimeRecords:   []TimeRecord{},
		ProjectBudgets: make(map[uuid.UUID]ProjectBudget),
	}
}
