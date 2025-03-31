package repo

import (
	"github.com/gofrs/uuid/v5"

	"t5.mkbrechtel.dev/model"
)

// Repository interface defines operations for managing all entities in the system
type Repository interface {
	// CreateTask stores a new task and returns its ID
	CreateTask(task model.Task) (uuid.UUID, error)

	// GetTask retrieves a task by its ID
	GetTask(id uuid.UUID) (model.Task, error)

	// UpdateTask modifies an existing task
	UpdateTask(task model.Task) error

	// DeleteTask removes a task by its ID
	DeleteTask(id uuid.UUID) error

	// ListTasks retrieves tasks based on optional filters
	ListTasks(filters []TaskFilter) ([]model.Task, error)

	// SearchTasks finds tasks matching the given query
	SearchTasks(query string) ([]model.Task, error)

	// RecordTaskUpdate stores a task update event
	RecordTaskUpdate(update model.TaskUpdate) error

	// RecordTaskStartTime stores a task start time event
	RecordTaskStartTime(startTime model.TaskStartTime) error

	// RecordTaskEndTime stores a task end time event
	RecordTaskEndTime(endTime model.TaskEndTime) error

	// CreateProject stores a new project and returns its ID
	CreateProject(project model.Project) (uuid.UUID, error)

	// GetProject retrieves a project by its ID
	GetProject(id uuid.UUID) (model.Project, error)

	// UpdateProject modifies an existing project
	UpdateProject(project model.Project) error

	// DeleteProject removes a project by its ID
	DeleteProject(id uuid.UUID) error

	// ListProjects retrieves all projects
	ListProjects() ([]model.Project, error)

	// GetTasksForProject retrieves all tasks associated with a project
	GetTasksForProject(projectID uuid.UUID) ([]model.Task, error)
}
