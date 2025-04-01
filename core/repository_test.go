package core

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.state)
	assert.NotNil(t, repo.EventStore)
	
	events, err := repo.GetEvents()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(events))
	assert.Equal(t, 0, len(repo.state.Tasks))
	assert.Equal(t, 0, len(repo.state.Projects))
}

func TestInMemoryEventStore(t *testing.T) {
	store := NewInMemoryEventStore()
	
	// Check initial state
	events, err := store.GetEvents()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(events))
	
	// Save an event
	event1 := NewTodoTxtTaskUpdate("Task 1 +project1", "test")
	err = store.SaveEvent(event1)
	assert.NoError(t, err)
	
	// Check that the event was stored
	events, err = store.GetEvents()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, event1, events[0])
	
	// Save another event
	event2 := NewTodoTxtTaskUpdate("Task 2 +project2", "test")
	err = store.SaveEvent(event2)
	assert.NoError(t, err)
	
	// Check that both events were stored
	events, err = store.GetEvents()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(events))
	
	// Verify the events are returned in the order they were added
	assert.Equal(t, event1, events[0])
	assert.Equal(t, event2, events[1])
	
	// Verify that modifying the returned slice doesn't affect the store
	events = append(events, NewTodoTxtTaskUpdate("Task 3", "test"))
	assert.Equal(t, 3, len(events))
	
	eventsAfter, err := store.GetEvents()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(eventsAfter))
}

func TestRepository_SaveEvent(t *testing.T) {
	repo := NewRepository()
	
	event := NewTodoTxtTaskUpdate("Test task +project1 @context1", "test")
	err := repo.SaveEvent(event)
	assert.NoError(t, err)
	
	// Verify the event was stored
	events, err := repo.GetEvents()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(events))
	
	// Verify the state was updated
	assert.Equal(t, 1, len(repo.state.Tasks))
	assert.Equal(t, 1, len(repo.state.Projects))
	
	// Check if the project was created
	_, hasProject := repo.state.Projects["project1"]
	assert.True(t, hasProject)
}

func TestRepository_GetEvents(t *testing.T) {
	repo := NewRepository()
	
	// Add some events
	event1 := NewTodoTxtTaskUpdate("Task 1 +project1", "test")
	event2 := NewTodoTxtTaskUpdate("Task 2 +project2", "test")
	
	assert.NoError(t, repo.SaveEvent(event1))
	assert.NoError(t, repo.SaveEvent(event2))
	
	// Get events
	events, err := repo.GetEvents()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(events))
	
	// Verify that modifying the returned slice doesn't affect the repository
	events = append(events, NewTodoTxtTaskUpdate("Task 3", "test"))
	assert.Equal(t, 3, len(events))
	
	// Check events again from repository
	eventsAfter, err := repo.GetEvents()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(eventsAfter))
}

func TestRepository_RebuildState(t *testing.T) {
	// Create a repository with a custom InMemoryEventStore for direct manipulation
	eventStore := NewInMemoryEventStore()
	repo := &Repository{
		EventStore: eventStore,
		state:      NewAppState(),
	}
	
	// Add events with timestamps out of order
	mockTime3 := time.Date(2023, 7, 15, 12, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime3 }
	event3 := NewTodoTxtTaskUpdate("Task 3 +project3", "test")
	
	mockTime1 := time.Date(2023, 7, 15, 10, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime1 }
	event1 := NewTodoTxtTaskUpdate("Task 1 +project1", "test")
	
	mockTime2 := time.Date(2023, 7, 15, 11, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime2 }
	event2 := NewTodoTxtTaskUpdate("Task 2 +project2", "test")
	
	// Add events in wrong order (3, 1, 2) directly to the event store
	eventStore.events = append(eventStore.events, event3, event1, event2)
	
	// Rebuild state
	err := repo.RebuildState()
	assert.NoError(t, err)
	
	// Verify all projects exist
	_, hasProject1 := repo.state.Projects["project1"]
	assert.True(t, hasProject1)
	_, hasProject2 := repo.state.Projects["project2"]
	assert.True(t, hasProject2)
	_, hasProject3 := repo.state.Projects["project3"]
	assert.True(t, hasProject3)
	
	// Verify all tasks exist
	assert.Equal(t, 3, len(repo.state.Tasks))
}

func TestRepository_TaskManagement(t *testing.T) {
	repo := NewRepository()
	
	// Reset now function
	originalNow := now
	defer func() { now = originalNow }()
	mockTime := time.Date(2023, 7, 15, 10, 0, 0, 0, time.UTC)
	now = func() time.Time { return mockTime }
	
	// Create a task
	task := Task{
		Todo:      "Test task",
		Priority:  "A",
		Projects:  []string{"project1"},
		Contexts:  []string{"context1"},
		CreatedDate: mockTime,
	}
	
	// Create the task
	id, err := repo.CreateTask(task)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
	
	// Get the task
	retrievedTask, err := repo.GetTask(id)
	assert.NoError(t, err)
	assert.Equal(t, task.Todo, retrievedTask.Todo)
	assert.Equal(t, task.Priority, retrievedTask.Priority)
	assert.Equal(t, 1, len(retrievedTask.Projects))
	assert.Equal(t, "project1", retrievedTask.Projects[0])
	
	// Update the task
	retrievedTask.Priority = "B"
	retrievedTask.Todo = "Updated task"
	err = repo.UpdateTask(retrievedTask)
	assert.NoError(t, err)
	
	// Get the updated task
	updatedTask, err := repo.GetTask(id)
	assert.NoError(t, err)
	assert.Equal(t, "Updated task", updatedTask.Todo)
	assert.Equal(t, "B", updatedTask.Priority)
	
	// List tasks with no filters
	tasks, err := repo.ListTasks(nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	
	// Delete the task
	err = repo.DeleteTask(id)
	assert.NoError(t, err)
	
	// List tasks after deletion
	tasks, err = repo.ListTasks(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tasks))
}

func TestRepository_ProjectManagement(t *testing.T) {
	repo := NewRepository()
	
	// Create a project directly
	project := Project{
		Name: "TestProject",
	}
	
	// Create the project
	id, err := repo.CreateProject(project)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
	
	// Get the project
	retrievedProject, err := repo.GetProject(id)
	assert.NoError(t, err)
	assert.Equal(t, project.Name, retrievedProject.Name)
	// No description to check in basic Project struct
	
	// List projects
	projects, err := repo.ListProjects()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(projects))
	
	// Update the project - in this simple version we just change the name
	retrievedProject.Name = "UpdatedProject"
	err = repo.UpdateProject(retrievedProject)
	assert.NoError(t, err)
	
	// Get the updated project
	updatedProject, err := repo.GetProject(id)
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedProject", updatedProject.Name)
	
	// Create a task in this project
	task := Task{
		Todo:     "Test task for project",
		Projects: []string{"UpdatedProject"},
	}
	taskID, err := repo.CreateTask(task)
	assert.NoError(t, err)
	
	// Get tasks for the project
	tasks, err := repo.GetTasksForProject(id)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, taskID, tasks[0].ID)
	
	// Delete the project
	err = repo.DeleteProject(id)
	assert.NoError(t, err)
	
	// List projects after deletion
	projects, err = repo.ListProjects()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(projects))
}