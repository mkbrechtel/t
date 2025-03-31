package core

import (
	"github.com/gofrs/uuid/v5"
)

// Repository interface defines operations for managing all entities in the system
type Repository interface {
	// CreateTask stores a new task and returns its ID
	CreateTask(task Task) (uuid.UUID, error)

	// GetTask retrieves a task by its ID
	GetTask(id uuid.UUID) (Task, error)

	// UpdateTask modifies an existing task
	UpdateTask(task Task) error

	// DeleteTask removes a task by its ID
	DeleteTask(id uuid.UUID) error

	// ListTasks retrieves tasks based on optional filters
	ListTasks(filters []TaskFilter) ([]Task, error)

	// SearchTasks finds tasks matching the given query
	SearchTasks(query string) ([]Task, error)

	// RecordTaskUpdate stores a task update event
	//RecordTaskUpdate(update TaskUpdate) error

	// RecordTaskStartTime stores a task start time event
	RecordTaskStartTime(startTime TaskStartTime) error

	// RecordTaskEndTime stores a task end time event
	RecordTaskEndTime(endTime TaskEndTime) error

	// CreateProject stores a new project and returns its ID
	CreateProject(project Project) (uuid.UUID, error)

	// GetProject retrieves a project by its ID
	GetProject(id uuid.UUID) (Project, error)

	// UpdateProject modifies an existing project
	UpdateProject(project Project) error

	// DeleteProject removes a project by its ID
	DeleteProject(id uuid.UUID) error

	// ListProjects retrieves all projects
	ListProjects() ([]Project, error)

	// GetTasksForProject retrieves all tasks associated with a project
	GetTasksForProject(projectID uuid.UUID) ([]Task, error)
}
