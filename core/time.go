package core

import (
	"fmt"
	"time"

	uuid "github.com/gofrs/uuid/v5"
	"t5.mkbrechtel.dev/t5/utils"
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
	Note      string
}

// GetType returns the type of the event
func (e *TaskStartTime) GetType() string {
	return "TaskStartTime"
}

// apply implements the Event interface
func (e *TaskStartTime) apply(state *AppState) error {
	// If there's already an active task, implicitly end it
	if state.ActiveTask != nil {
		implicitEnd := &TaskEndTime{
			ID:      utils.NewUUID(),
			TaskID:  state.ActiveTask.TaskID,
			EndTime: e.StartTime,
		}
		implicitEnd.apply(state)
	}

	// Set the new active task
	state.ActiveTask = &ActiveTaskInfo{
		TaskID:      e.TaskID,
		StartTime:   e.StartTime,
		TotalPaused: 0,
		EventID:     e.ID,
	}

	// Initialize task time if this is the first time tracking for this task
	task, exists := state.Tasks[e.TaskID]
	if exists {
		// Update the task to show it's now active
		task.Active = true
		state.Tasks[e.TaskID] = task
	} else {
		return fmt.Errorf("task with ID %s not found", e.TaskID)
	}

	return nil
}

// TaskEndTime event records when a task was stopped
type TaskEndTime struct {
	BaseEvent
	ID      uuid.UUID
	TaskID  uuid.UUID
	EndTime time.Time
	Note    string
}

// GetType returns the type of the event
func (e *TaskEndTime) GetType() string {
	return "TaskEndTime"
}

// apply implements the Event interface
func (e *TaskEndTime) apply(state *AppState) error {
	// Check if there's an active task
	if state.ActiveTask == nil || state.ActiveTask.TaskID != e.TaskID {
		return fmt.Errorf("no active task with ID %s found", e.TaskID)
	}

	// Calculate effective work duration, accounting for pauses
	var effectiveDuration time.Duration
	if state.ActiveTask.LastPaused.IsZero() {
		// Not paused, calculate from start time to end time
		effectiveDuration = e.EndTime.Sub(state.ActiveTask.StartTime) - state.ActiveTask.TotalPaused
	} else {
		// Currently paused, calculate from start time to last pause time
		effectiveDuration = state.ActiveTask.LastPaused.Sub(state.ActiveTask.StartTime) - state.ActiveTask.TotalPaused
	}

	// Create a time record
	timeRecord := TimeRecord{
		StartTime: state.ActiveTask.StartTime,
		EndTime:   e.EndTime,
		Duration:  effectiveDuration,
		TaskID:    e.TaskID,
		Note:      e.Note,
	}

	// Add the time record to the global time records
	state.TimeRecords = append(state.TimeRecords, timeRecord)

	// Update the task's used time
	task, exists := state.Tasks[e.TaskID]
	if exists {
		task.UsedTime += effectiveDuration
		task.Active = false
		
		// Update daily time spent
		today := e.EndTime.Format("2006-01-02")
		if task.DailyTimeSpent == nil {
			task.DailyTimeSpent = make(map[string]time.Duration)
		}
		task.DailyTimeSpent[today] += effectiveDuration
		
		state.Tasks[e.TaskID] = task

		// Update project used time if the task has projects
		for _, projectName := range task.Projects {
			if project, ok := state.Projects[projectName]; ok {
				project.UsedTime += effectiveDuration
				state.Projects[projectName] = project
			}
		}
	}

	// Clear the active task
	state.ActiveTask = nil

	return nil
}

// TaskPauseTime event records when work on a task was paused
type TaskPauseTime struct {
	BaseEvent
	ID        uuid.UUID
	TaskID    uuid.UUID
	PauseTime time.Time
	Reason    string
}

// GetType returns the type of the event
func (e *TaskPauseTime) GetType() string {
	return "TaskPauseTime"
}

// apply implements the Event interface
func (e *TaskPauseTime) apply(state *AppState) error {
	// Check if there's an active task
	if state.ActiveTask == nil || state.ActiveTask.TaskID != e.TaskID {
		return fmt.Errorf("no active task with ID %s found", e.TaskID)
	}

	// Check if the task is already paused
	if !state.ActiveTask.LastPaused.IsZero() {
		return fmt.Errorf("task is already paused")
	}

	// Mark the task as paused
	state.ActiveTask.LastPaused = e.PauseTime

	return nil
}

// TaskResumeTime event records when work on a task was resumed
type TaskResumeTime struct {
	BaseEvent
	ID         uuid.UUID
	TaskID     uuid.UUID
	ResumeTime time.Time
}

// GetType returns the type of the event
func (e *TaskResumeTime) GetType() string {
	return "TaskResumeTime"
}

// apply implements the Event interface
func (e *TaskResumeTime) apply(state *AppState) error {
	// Check if there's an active task
	if state.ActiveTask == nil || state.ActiveTask.TaskID != e.TaskID {
		return fmt.Errorf("no active task with ID %s found", e.TaskID)
	}

	// Check if the task is actually paused
	if state.ActiveTask.LastPaused.IsZero() {
		return fmt.Errorf("task is not paused")
	}

	// Calculate pause duration and add to total
	pauseDuration := e.ResumeTime.Sub(state.ActiveTask.LastPaused)
	state.ActiveTask.TotalPaused += pauseDuration

	// Clear the pause time
	state.ActiveTask.LastPaused = time.Time{}

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

// apply implements the Event interface
func (e *SetTimeBudgetForProject) apply(state *AppState) error {
	// Find the project by ID
	var projectName string
	found := false

	for name, project := range state.Projects {
		if project.ID == e.ProjectID {
			projectName = name
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("project with ID %s not found", e.ProjectID)
	}

	// Update the project with budget information
	project := state.Projects[projectName]
	project.Budget = e.Budget
	project.Description = e.Description

	// Only use the provided UsedTime if it's larger than what we've calculated
	if e.UsedTime > project.UsedTime {
		project.UsedTime = e.UsedTime
	}

	// Update the project in the state
	state.Projects[projectName] = project

	// Also update the project budget tracking
	state.ProjectBudgets[e.ProjectID] = ProjectBudget{
		ProjectID:   e.ProjectID,
		Budget:      e.Budget,
		UsedTime:    project.UsedTime,
		Description: e.Description,
	}

	return nil
}

// NewTaskStartTime creates a new TaskStartTime event with current timestamp
func NewTaskStartTime(taskID uuid.UUID, note string) *TaskStartTime {
	event := &TaskStartTime{
		ID:        utils.NewUUID(),
		TaskID:    taskID,
		StartTime: now(),
		Note:      note,
	}
	event.Timestamp = event.StartTime
	return event
}

// NewTaskEndTime creates a new TaskEndTime event with current timestamp
func NewTaskEndTime(taskID uuid.UUID, note string) *TaskEndTime {
	event := &TaskEndTime{
		ID:      utils.NewUUID(),
		TaskID:  taskID,
		EndTime: now(),
		Note:    note,
	}
	event.Timestamp = event.EndTime
	return event
}

// NewTaskPauseTime creates a new TaskPauseTime event with current timestamp
func NewTaskPauseTime(taskID uuid.UUID, reason string) *TaskPauseTime {
	event := &TaskPauseTime{
		ID:        utils.NewUUID(),
		TaskID:    taskID,
		PauseTime: now(),
		Reason:    reason,
	}
	event.Timestamp = event.PauseTime
	return event
}

// NewTaskResumeTime creates a new TaskResumeTime event with current timestamp
func NewTaskResumeTime(taskID uuid.UUID) *TaskResumeTime {
	event := &TaskResumeTime{
		ID:         utils.NewUUID(),
		TaskID:     taskID,
		ResumeTime: now(),
	}
	event.Timestamp = event.ResumeTime
	return event
}

// NewSetTimeBudgetForProject creates a new SetTimeBudgetForProject event with current timestamp
func NewSetTimeBudgetForProject(projectID uuid.UUID, budget time.Duration, usedTime time.Duration, description string) *SetTimeBudgetForProject {
	event := &SetTimeBudgetForProject{
		ID:          utils.NewUUID(),
		ProjectID:   projectID,
		Budget:      budget,
		UsedTime:    usedTime,
		Description: description,
	}
	event.Timestamp = now()
	return event
}