# t5

**t5**; the **t**ask manager, **t**odo list, **t**ime tracker and [**t**imeboxing](https://en.wikipedia.org/wiki/Timeboxing) **t**ool

It uses the [todo.txt](http://todotxt.org/) format for your todo list.

_WIP_

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
