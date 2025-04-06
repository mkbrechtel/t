# Time Tracking System Design

## Overview

The time tracking system in t5 records when users work on tasks, allowing them to track time spent on different activities and manage project budgets. This document outlines the event types, state changes, and commands related to time tracking.

## New Event Types

### TaskStartTime

Records when a user starts working on a task.

- **ID**: Unique identifier for this start event
- **TaskID**: ID of the task being worked on
- **StartTime**: Timestamp when work started
- **Note**: Optional note about the activity

```go
// TaskStartTime event records when a task was started
type TaskStartTime struct {
    BaseEvent
    ID        uuid.UUID
    TaskID    uuid.UUID
    StartTime time.Time
    Note      string
}
```

### TaskEndTime

Records when a user stops working on a task.

- **ID**: Unique identifier for this end event
- **TaskID**: ID of the task that was worked on
- **EndTime**: Timestamp when work ended
- **Note**: Optional note about the completed work

```go
// TaskEndTime event records when a task was stopped
type TaskEndTime struct {
    BaseEvent
    ID      uuid.UUID
    TaskID  uuid.UUID
    EndTime time.Time
    Note    string
}
```

### TaskPauseTime

Records when a user temporarily pauses work on a task.

- **ID**: Unique identifier for this pause event
- **TaskID**: ID of the task being paused
- **PauseTime**: Timestamp when work was paused
- **Reason**: Optional reason for pausing

```go
// TaskPauseTime event records when work on a task was paused
type TaskPauseTime struct {
    BaseEvent
    ID        uuid.UUID
    TaskID    uuid.UUID
    PauseTime time.Time
    Reason    string
}
```

### TaskResumeTime

Records when a user resumes work on a previously paused task.

- **ID**: Unique identifier for this resume event
- **TaskID**: ID of the task being resumed
- **ResumeTime**: Timestamp when work was resumed

```go
// TaskResumeTime event records when work on a task was resumed
type TaskResumeTime struct {
    BaseEvent
    ID         uuid.UUID
    TaskID     uuid.UUID
    ResumeTime time.Time
}
```

### SetTimeBudgetForProject

Sets a time budget for a project.

- **ID**: Unique identifier for this budget event
- **ProjectID**: ID of the project getting a budget
- **Budget**: Total allocated time for the project
- **UsedTime**: Time already used (for imported projects)
- **Description**: Optional description of the budget

```go
// SetTimeBudgetForProject event sets a time budget for a project
type SetTimeBudgetForProject struct {
    BaseEvent
    ID          uuid.UUID
    ProjectID   uuid.UUID
    Budget      time.Duration
    UsedTime    time.Duration
    Description string
}
```

## State Design

The time tracking information is stored in these parts of the application state:

1. **Active Task**: The task currently being worked on
   ```go
   type ActiveTaskInfo struct {
       TaskID      uuid.UUID
       StartTime   time.Time
       LastPaused  time.Time // Zero if not paused
       TotalPaused time.Duration
       EventID     uuid.UUID // ID of the TaskStartTime event
   }
   ```

2. **Task Time Records**: History of time spent on tasks
   ```go
   type TimeRecord struct {
       StartTime   time.Time
       EndTime     time.Time
       Duration    time.Duration
       TaskID      uuid.UUID
       Note        string
   }
   ```

3. **Project Time Budgets**:
   ```go
   type ProjectBudget struct {
       ProjectID   uuid.UUID
       Budget      time.Duration
       UsedTime    time.Duration
       Description string
   }
   ```

## Event Application Logic

### Applying TaskStartTime

When a `TaskStartTime` event is received:

1. If there's already an active task:
   - Generate an implicit `TaskEndTime` event for the current active task
   - Apply that event first to properly close the current session

2. Set the new task as active:
   ```go
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
       }

       return nil
   }
   ```

### Applying TaskEndTime

When a `TaskEndTime` event is received:

1. Calculate the duration worked on the task
2. Add a time record to the task's history
3. Update the task's total used time
4. Update project used time
5. Clear the active task

```go
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

    // Add the time record to the task's history
    task, exists := state.Tasks[e.TaskID]
    if exists {
        task.UsedTime += effectiveDuration
        task.TimeRecords = append(task.TimeRecords, timeRecord)
        task.Active = false
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
```

### Applying TaskPauseTime

When a `TaskPauseTime` event is received:

1. Mark the active task as paused
2. Record the pause time

```go
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
```

### Applying TaskResumeTime

When a `TaskResumeTime` event is received:

1. Calculate the pause duration
2. Add to the total paused time
3. Clear the pause time

```go
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
```

### Applying SetTimeBudgetForProject

When a `SetTimeBudgetForProject` event is received:

1. Update or create the project budget information

```go
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

    // Only use the provided UsedTime if it's larger than what we've calculated
    if e.UsedTime > project.UsedTime {
        project.UsedTime = e.UsedTime
    }

    // Update the project in the state
    state.Projects[projectName] = project

    return nil
}
```

## CLI Commands

The following CLI commands will be implemented for time tracking:

### Start Command

```
t5 start [task-id] [--note "Starting work on feature X"]
```

This command creates a `TaskStartTime` event and starts tracking time for the specified task.

### Stop Command

```
t5 stop [--note "Completed implementation"]
```

This command creates a `TaskEndTime` event for the currently active task.

### Pause Command

```
t5 pause [--reason "Lunch break"]
```

This command creates a `TaskPauseTime` event for the currently active task.

### Resume Command

```
t5 resume
```

This command creates a `TaskResumeTime` event for the currently paused task.

### Status Command

```
t5 status
```

Shows the current active task, how long it's been running, and if it's paused.

### Budget Command

```
t5 budget set [project] [budget] [--description "Q3 allocation"]
t5 budget show [project]
t5 budget list
```

The `set` subcommand creates a `SetTimeBudgetForProject` event.
The `show` and `list` subcommands display budget information.

### Report Command

```
t5 report [--daily|--weekly|--monthly] [--from DATE] [--to DATE] [--project PROJECT]
```

Generates a time report showing time spent on tasks within the specified period and project.

## Implementation Plan

1. Define event types and state changes in the core package
2. Implement event application logic for each new event type
3. Create CLI commands for time tracking operations
4. Add reporting functionality
5. Implement project budget tracking
6. Add i3 integration for status display

## Future Enhancements

- Pomodoro technique integration with automated breaks
- Calendar integration for time blocking
- Visualization of time tracking data in the web UI
- Automatic idle detection and handling
- Multi-device synchronization for time tracking
- Export time reports to various formats (CSV, PDF, etc.)