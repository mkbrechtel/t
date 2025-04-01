package core

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gofrs/uuid/v5"
)

// Repository is an in-memory implementation of event storage and state management
type Repository struct {
	events []Event
	state  *AppState
	mu     sync.RWMutex
}

// NewRepository creates a new empty repository
func NewRepository() *Repository {
	return &Repository{
		events: []Event{},
		state:  NewAppState(),
	}
}

// SaveEvent persists an event and applies it to the current state
func (r *Repository) SaveEvent(event Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Apply the event to the current state
	if err := event.apply(r.state); err != nil {
		return err
	}

	// Store the event
	r.events = append(r.events, event)
	return nil
}

// GetEvents retrieves all events from storage
func (r *Repository) GetEvents() ([]Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy of the events slice to prevent modification
	eventsCopy := make([]Event, len(r.events))
	copy(eventsCopy, r.events)
	return eventsCopy, nil
}

// RecordTodoTxtTaskUpdate stores a todo.txt task update event
func (r *Repository) RecordTodoTxtTaskUpdate(update *TodoTxtTaskUpdate) error {
	return r.SaveEvent(update)
}

// RecordTaskStartTime stores a task start time event
func (r *Repository) RecordTaskStartTime(startTime *TaskStartTime) error {
	return r.SaveEvent(startTime)
}

// RecordTaskEndTime stores a task end time event
func (r *Repository) RecordTaskEndTime(endTime *TaskEndTime) error {
	return r.SaveEvent(endTime)
}

// GetAppState returns the current application state
func (r *Repository) GetAppState() *AppState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Return the current state
	return r.state
}

// RebuildState rebuilds the state from events
func (r *Repository) RebuildState() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Create a new empty state
	newState := NewAppState()
	
	// Sort events by timestamp
	sort.Slice(r.events, func(i, j int) bool {
		return r.events[i].GetTimestamp().Before(r.events[j].GetTimestamp())
	})
	
	// Apply all events in order
	for _, event := range r.events {
		if err := event.apply(newState); err != nil {
			return fmt.Errorf("failed to apply event %s: %w", event.GetType(), err)
		}
	}
	
	// Update the state
	r.state = newState
	return nil
}

// *** Task Management ***

// CreateTask creates a new task and returns its ID
// This is implemented via a TodoTxtTaskUpdate event
func (r *Repository) CreateTask(task Task) (uuid.UUID, error) {
	// Ensure the task has a UUID
	if task.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		task.ID = id
	}
	
	// Convert task to todo.txt format
	todoTxt := task.ToTodoTxt()
	
	// Create and record the event
	event := NewTodoTxtTaskUpdate(todoTxt)
	err := r.SaveEvent(event)
	if err != nil {
		return uuid.Nil, err
	}
	
	return task.ID, nil
}

// GetTask retrieves a task by its ID
func (r *Repository) GetTask(id uuid.UUID) (Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	task, ok := r.state.Tasks[id]
	if !ok {
		return Task{}, errors.New("task not found")
	}
	
	return task, nil
}

// UpdateTask modifies an existing task
func (r *Repository) UpdateTask(task Task) error {
	// Make sure the task exists
	r.mu.RLock()
	_, ok := r.state.Tasks[task.ID]
	r.mu.RUnlock()
	
	if !ok {
		return errors.New("task not found")
	}
	
	// Convert task to todo.txt format
	todoTxt := task.ToTodoTxt()
	
	// Create and record the event
	event := NewTodoTxtTaskUpdate(todoTxt)
	return r.SaveEvent(event)
}

// DeleteTask removes a task by its ID
// Note: In an event-sourced system, we don't physically delete tasks,
// but mark them as deleted through an event
func (r *Repository) DeleteTask(id uuid.UUID) error {
	// For now, we'll implement this by creating a task update
	// that removes all content but keeps the ID
	r.mu.RLock()
	_, ok := r.state.Tasks[id]
	r.mu.RUnlock()
	
	if !ok {
		return errors.New("task not found")
	}
	
	// Create an empty task with just the ID preserved
	emptyTask := Task{
		ID:            id,
		AdditionalTags: map[string]string{"deleted": "true"},
	}
	
	todoTxt := emptyTask.ToTodoTxt()
	event := NewTodoTxtTaskUpdate(todoTxt)
	return r.SaveEvent(event)
}

// ListTasks retrieves tasks based on optional filters
func (r *Repository) ListTasks(filters []TaskFilter) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []Task
	
	// If no filters, return all tasks
	if len(filters) == 0 {
		result = make([]Task, 0, len(r.state.Tasks))
		for _, task := range r.state.Tasks {
			// Skip deleted tasks
			if _, isDeleted := task.AdditionalTags["deleted"]; !isDeleted {
				result = append(result, task)
			}
		}
		return result, nil
	}
	
	// Apply filters
	for _, task := range r.state.Tasks {
		// Skip deleted tasks
		if _, isDeleted := task.AdditionalTags["deleted"]; isDeleted {
			continue
		}
		
		// Check if task matches all filters
		match := true
		for _, filter := range filters {
			if !filter.FilterTask(task) {
				match = false
				break
			}
		}
		
		if match {
			result = append(result, task)
		}
	}
	
	return result, nil
}

// SearchTasks finds tasks matching the given query string
func (r *Repository) SearchTasks(query string) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	query = strings.ToLower(query)
	var result []Task
	
	for _, task := range r.state.Tasks {
		// Skip deleted tasks
		if _, isDeleted := task.AdditionalTags["deleted"]; isDeleted {
			continue
		}
		
		// Simple text search in the todo content
		if strings.Contains(strings.ToLower(task.Todo), query) {
			result = append(result, task)
		}
	}
	
	return result, nil
}

// *** Project Management ***

// CreateProject stores a new project and returns its ID
func (r *Repository) CreateProject(project Project) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Ensure the project has an ID
	if project.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return uuid.Nil, err
		}
		project.ID = id
	}
	
	// Add to state
	r.state.Projects[project.Name] = project
	
	return project.ID, nil
}

// GetProject retrieves a project by its ID
func (r *Repository) GetProject(id uuid.UUID) (Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Search for project by ID
	for _, project := range r.state.Projects {
		if project.ID == id {
			return project, nil
		}
	}
	
	return Project{}, errors.New("project not found")
}

// UpdateProject modifies an existing project
func (r *Repository) UpdateProject(project Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Find the old project name
	var oldName string
	found := false
	for name, p := range r.state.Projects {
		if p.ID == project.ID {
			oldName = name
			found = true
			break
		}
	}
	
	if !found {
		return errors.New("project not found")
	}
	
	// If the name changed, remove the old entry
	if oldName != project.Name {
		delete(r.state.Projects, oldName)
	}
	
	// Update project
	r.state.Projects[project.Name] = project
	
	return nil
}

// DeleteProject removes a project by its ID
func (r *Repository) DeleteProject(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Find the project
	var projectName string
	found := false
	for name, project := range r.state.Projects {
		if project.ID == id {
			projectName = name
			found = true
			break
		}
	}
	
	if !found {
		return errors.New("project not found")
	}
	
	// Remove the project
	delete(r.state.Projects, projectName)
	
	return nil
}

// ListProjects retrieves all projects
func (r *Repository) ListProjects() ([]Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	projects := make([]Project, 0, len(r.state.Projects))
	for _, project := range r.state.Projects {
		projects = append(projects, project)
	}
	
	return projects, nil
}

// GetTasksForProject retrieves all tasks associated with a project
func (r *Repository) GetTasksForProject(projectID uuid.UUID) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Find the project name
	var projectName string
	found := false
	for name, project := range r.state.Projects {
		if project.ID == projectID {
			projectName = name
			found = true
			break
		}
	}
	
	if !found {
		return nil, errors.New("project not found")
	}
	
	// Find tasks that have this project
	var tasks []Task
	for _, task := range r.state.Tasks {
		// Skip deleted tasks
		if _, isDeleted := task.AdditionalTags["deleted"]; isDeleted {
			continue
		}
		
		for _, p := range task.Projects {
			if p == projectName {
				tasks = append(tasks, task)
				break
			}
		}
	}
	
	return tasks, nil
}