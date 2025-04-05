# t5

**t5**; the **t**ask manager, **t**odo list, **t**ime tracker and [**t**imeboxing](https://en.wikipedia.org/wiki/Timeboxing) **t**ool

It uses the [todo.txt](http://todotxt.org/) format for your todo list.

_WIP_

## Usage

Show help:
```
t5 help
```

Available commands:

* `t5 list` - List all tasks with their priorities and metadata
* `t5 update [file]` - Update and ensure properties of tasks (IDs, dates)
* `t5 todo [update [file]]` - Alias for update command
* `t5 sync [file]` - Synchronize tasks between todo.txt and event store
* `t5 config` - Display current configuration settings

Common flags:
* `--todo, -t` - Specify todo.txt file path
* `--config, -c` - Specify config file path
* `--eventstore` - Specify event store file path
* `--prefer-short-ids` - Use short form IDs
* `--enforce-creation-date` - Ensure tasks have creation dates
* `--enforce-completion-date` - Ensure completed tasks have completion dates

## Build

```
go mod download
go build -o t5
```

run `./t5`

## Tests

```
go test -v ./...
```

End-to-end tests are located in the `test` package and are run using the `go test` command too.

## Architecture

t5 is built on an event-sourcing architecture, where all changes to the application state are recorded as events. The current state of the application can be reconstructed by replaying these events in order.

### Core Concepts

- **Events**: Immutable records of something that happened (e.g., a task was created, updated, or completed)
- **AppState**: The current snapshot of the application, built by applying events
- **Event Store**: Persistent storage for events
- **Replay**: The process of rebuilding the application state by applying events in chronological order

For more details on the event system, see [Event Replay Documentation](docs/design/event-replay.md).

### Event Replay System

The event replay system works as follows:

1. **Event Creation**: When a change occurs (e.g., adding a task), an event is created with all necessary data
2. **Event Storage**: Events are stored in chronological order with timestamps
3. **State Application**: Events modify the application state through their `apply()` method
4. **State Reconstruction**: The current state can be rebuilt at any time by replaying all events

Key benefits of this approach:
- Complete history of all changes
- Ability to reconstruct state at any point in time
- Natural support for undo/redo operations
- Separation of write and read models

Currently implemented event types:
- `TodoTxtTaskUpdate`: Updates tasks from todo.txt format

Example usage:
```go
// Create a repository that manages events and state
repo := NewRepository()

// Add a new todo task via an event
event := NewTodoTxtTaskUpdate("(A) 2023-07-15 Implement feature +project @context")
repo.SaveEvent(event)

// Access the current state
state := repo.GetAppState()
for id, task := range state.Tasks {
    fmt.Printf("Task: %s (Priority: %s)\n", task.Todo, task.Priority)
}
```

See the complete [Event Replay Documentation](docs/design/event-replay.md) for more information.

### Data Model

The data model for t5 includes:

- **Tasks**: Individual tasks with properties like priority, due date, and completion status
- **Projects**: Collections of related tasks, with optional time budgets
- **Events**: Various event types that modify the application state

See the complete [Data Model Documentation](docs/design/data-model.md) for more information.

## Feature Implementation Plan

### Phase 1: Core Task Management Enhancements

- [ ] **Task Command Improvements**
   - [ ] Implement `t5 add todo` to add tasks directly from command line
   - [ ] Support stdin input for task creation
   - [ ] Expand task filtering capabilities based on todo.txt properties
   - [ ] Implement regex-based filters
   - [ ] Add boolean combinators for filters (AND, OR, NOT)

- [ ] **Task Priority Management**
   - [ ] Implement priority assignment/reassignment
   - [ ] Implement priority-based sorting for display
   - [ ] Add priority filtering in task listings

- [ ] **Task Modification Commands**
   - [ ] Add commands to modify existing tasks
   - [ ] Implement project/context addition/removal
   - [ ] Add bulk task modification options

### Phase 2: Time Tracking and Activity Management

- [ ] **Activity Selection**
   - [ ] Implement `t5 start [task-id]` to begin time tracking
   - [ ] Design `TaskStartTime` and `TaskEndTime` event applications
   - [ ] Add `t5 pause` and `t5 resume` commands

- [ ] **Time Budget Management**
   - [ ] Implement `t5 budget set [project] [budget]` command
   - [ ] Implement budget usage tracking
   - [ ] Add time budget reporting

- [ ] **Activity Reporting**
   - [ ] Implement `t5 status` to show current activity
   - [ ] Create daily/weekly/monthly time reports
   - [ ] Add formatting options for reports

### Phase 3: Synchronization Capabilities

- [ ] **Enhanced Sync Options**
   - [ ] Complete todo.txt synchronization
   - [ ] Add support for multiple todo.txt files
   - [ ] Implement intelligent merging of todo files

- [ ] **External Services Integration**
   - [ ] Complete GitHub sync
   - [ ] Complete GitLab sync
   - [ ] Complete OpenProject sync
   - [ ] Implement CalDAV sync

- [ ] **Advanced Sync Configuration**
   - [ ] Add filter-based sync options
   - [ ] Implement task transformation during sync
   - [ ] Support for bi-directional sync

### Phase 4: User Experience Improvements

- [ ] **UI Enhancements**
   - [ ] Improve CLI output formatting
   - [ ] Add color support for terminal output
   - [ ] Implement interactive mode for task selection

- [ ] **i3 Integration**
   - [ ] Create i3 bar component for current activity
   - [ ] Add i3blocks status output
   - [ ] Implement i3 mode templates

- [ ] **Timeboxing Features**
   - [ ] Add Pomodoro technique support
   - [ ] Implement countdown timers
   - [ ] Add notifications for time boundaries

### Phase 5: Web Application

- [ ] **Backend API**
   - [ ] Design and implement REST API endpoints
   - [ ] Create WebSocket server for live updates
   - [ ] Implement authentication and authorization

- [ ] **Web Frontend**
   - [ ] Create responsive web interface
   - [ ] Implement task management UI
   - [ ] Add time tracking features
   - [ ] Develop reporting dashboards

- [ ] **Mobile Support**
   - [ ] Ensure responsive design works on mobile devices
   - [ ] Add PWA support for offline usage

## Implementation Milestones

### Milestone 1: Complete Core Task Management (1-2 weeks)
- [ ] Full task CRUD operations
- [ ] Comprehensive filtering system
- [ ] Improved command line experience

### Milestone 2: Time Tracking System (2-3 weeks)
- [ ] Working activity tracking
- [ ] Budget management
- [ ] Basic reporting

### Milestone 3: Synchronization Framework (3-4 weeks)
- [ ] Complete todo.txt sync
- [ ] At least two external service integrations
- [ ] Multi-source merging capabilities

### Milestone 4: User Experience & i3 Integration (2 weeks)
- [ ] Enhanced CLI experience
- [ ] i3 integration components
- [ ] Timeboxing features

### Milestone 5: Web Application (4-6 weeks)
- [ ] Working REST API
- [ ] Basic web interface
- [ ] Mobile-responsive design

## Immediate Next Steps

- [ ] **Implement Task Command Enhancements**
   - [ ] Add support for adding tasks from command line
   - [ ] Implement stdin task creation
   - [ ] Create basic task modification commands

- [ ] **Complete Time Tracking Logic**
   - [ ] Implement the `apply()` method for `TaskStartTime` and `TaskEndTime` events
   - [ ] Create task duration calculations
   - [ ] Add start/pause/resume commands

- [ ] **Enhance Task Filtering**
   - [ ] Implement todo.txt property filters
   - [ ] Add regex filter support
   - [ ] Create boolean combinators for complex filters

- [ ] **Improve Task Listing Display**
   - [ ] Add formatted output options
   - [ ] Implement sorting capabilities
   - [ ] Add grouping by project/context

- [ ] **Setup CI/CD Pipeline**
   - [ ] Configure GitHub Actions for testing and deployment
   - [ ] Add release automation
   - [ ] Implement code quality checks
