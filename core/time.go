package core

import (
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

// Time provider function variable
// Can be overridden in tests to provide deterministic time
var now = time.Now

// TaskStartTime event records when a task was started
type TaskStartTime struct {
	BaseEvent
	ID        uuid.UUID
	TaskID    uuid.UUID
	StartTime time.Time
}

// GetType returns the type of the event
func (e *TaskStartTime) GetType() string {
	return "TaskStartTime"
}

// Apply applies this event to the app state
func (e *TaskStartTime) apply(state *AppState) error {
	// This is a placeholder for future implementation
	return nil
}

// TaskEndTime event records when a task was completed
type TaskEndTime struct {
	BaseEvent
	ID      uuid.UUID
	TaskID  uuid.UUID
	EndTime time.Time
}

// GetType returns the type of the event
func (e *TaskEndTime) GetType() string {
	return "TaskEndTime"
}

// Apply applies this event to the app state
func (e *TaskEndTime) apply(state *AppState) error {
	// This is a placeholder for future implementation
	return nil
}

// SetTimeBudgetForProject event sets a time budget for a project
type SetTimeBudgetForProject struct {
	BaseEvent
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Budget      time.Duration
	UsedTime    time.Duration
	Description string
}

// GetType returns the type of the event
func (e *SetTimeBudgetForProject) GetType() string {
	return "SetTimeBudgetForProject"
}

// Apply applies this event to the app state
func (e *SetTimeBudgetForProject) apply(state *AppState) error {
	// This is a placeholder for future implementation
	return nil
}
