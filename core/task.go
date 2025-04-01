package core

import (
	"strings"
	"time"

	uuid "github.com/gofrs/uuid/v5"
	todo "github.com/1set/todotxt"
)

type Task struct {
	ID             uuid.UUID
	Original       string
	Todo           string
	Priority       string
	Projects       []string
	Contexts       []string
	AdditionalTags map[string]string
	CreatedDate    time.Time
	DueDate        time.Time
	CompletedDate  time.Time
	Completed      bool
	UsedTime       time.Duration
	Source         string
}

// ToTodoTxt converts a Task to todo.txt format
func (t Task) ToTodoTxt() string {
	// Ensure the UUID is stored in AdditionalTags
	if t.AdditionalTags == nil {
		t.AdditionalTags = make(map[string]string)
	}
	t.AdditionalTags["uuid"] = t.ID.String()
	
	// Create a todo.txt task
	todoTask := todo.Task{
		Todo:           t.Todo,
		Priority:       t.Priority,
		Projects:       t.Projects,
		Contexts:       t.Contexts,
		AdditionalTags: t.AdditionalTags,
		Completed:      t.Completed,
	}
	
	// Set dates if they are not zero
	if !t.CreatedDate.IsZero() {
		todoTask.CreatedDate = t.CreatedDate
	}
	if !t.DueDate.IsZero() {
		todoTask.DueDate = t.DueDate
	}
	if !t.CompletedDate.IsZero() {
		todoTask.CompletedDate = t.CompletedDate
	}
	
	// Generate string representation
	return todoTask.String()
}

// TodoTxtTaskUpdate represents an event for updating tasks from todo.txt format
type TodoTxtTaskUpdate struct {
	BaseEvent
	Lines string
	Source string
}

// GetType returns the type of the event
func (e *TodoTxtTaskUpdate) GetType() string {
	return "TodoTxtTaskUpdate"
}

func (e *TodoTxtTaskUpdate) apply(state *AppState) error {
	// Parse the todo.txt lines into tasks
	var todoTasks todo.TaskList
	for _, line := range strings.Split(e.Lines, "\n") {
		if line == "" {
			continue
		}
		task, err := todo.ParseTask(line)
		if err != nil {
			return err
		}
		todoTasks = append(todoTasks, *task)
	}

	// Process each todo.txt task and update the app state
	for _, todoTask := range todoTasks {
		// Generate or extract a UUID for the task
		taskID := uuid.FromStringOrNil(todoTask.AdditionalTags["uuid"])
		if taskID == uuid.Nil {
			// Create a new UUID if none exists
			newID, err := uuid.NewV4()
			if err != nil {
				return err
			}
			taskID = newID

			// Add the UUID back to the task
			if todoTask.AdditionalTags == nil {
				todoTask.AdditionalTags = make(map[string]string)
			}
			todoTask.AdditionalTags["uuid"] = taskID.String()
		}

		// Convert the todo.txt task to our Task model
		task := Task{
			ID:             taskID,
			Original:       todoTask.Original,
			Todo:           todoTask.Todo,
			Priority:       todoTask.Priority,
			Projects:       todoTask.Projects,
			Contexts:       todoTask.Contexts,
			AdditionalTags: todoTask.AdditionalTags,
			CreatedDate:    todoTask.CreatedDate,
			DueDate:        todoTask.DueDate,
			CompletedDate:  todoTask.CompletedDate,
			Completed:      todoTask.Completed,
			Source:         e.Source,
		}

		// Update the state
		state.Tasks[taskID] = task

		// Update project references if needed
		for _, projectName := range task.Projects {
			// Check if project exists, if not create it
			if _, exists := state.Projects[projectName]; !exists {
				projectID, err := uuid.NewV4()
				if err != nil {
					return err
				}

				state.Projects[projectName] = Project{
					ID:   projectID,
					Name: projectName,
				}
			}
		}
	}

	return nil
}
