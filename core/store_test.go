package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendLogEventStore(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "events.jsonl")

	// Create a new store
	store, err := NewAppendLogEventStore(filePath)
	require.NoError(t, err)
	require.NotNil(t, store)

	// Create a mock time provider for deterministic testing
	originalNow := now
	defer func() { now = originalNow }()

	mockTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return mockTime
	}

	// Create some test events
	taskID, err := uuid.NewV4()
	require.NoError(t, err)

	events := []Event{
		&TodoTxtTaskUpdate{
			BaseEvent: BaseEvent{Timestamp: now()},
			Lines:     "(A) 2023-07-15 Test task +project @context",
		},
		&TaskStartTime{
			BaseEvent: BaseEvent{Timestamp: now()},
			ID:        uuid.Must(uuid.NewV4()),
			TaskID:    taskID,
			StartTime: now(),
		},
		&TaskEndTime{
			BaseEvent: BaseEvent{Timestamp: now()},
			ID:        uuid.Must(uuid.NewV4()),
			TaskID:    taskID,
			EndTime:   now(),
		},
	}

	// Save the events
	for _, event := range events {
		err := store.SaveEvent(event)
		require.NoError(t, err)
	}

	// Retrieve the events
	retrievedEvents, err := store.GetEvents()
	require.NoError(t, err)
	require.Len(t, retrievedEvents, len(events))

	// Check that each event was correctly saved and retrieved
	for i, original := range events {
		retrieved := retrievedEvents[i]
		assert.Equal(t, original.GetType(), retrieved.GetType())
		assert.Equal(t, original.GetTimestamp(), retrieved.GetTimestamp())

		// Check specific fields based on event type
		switch orig := original.(type) {
		case *TodoTxtTaskUpdate:
			ret, ok := retrieved.(*TodoTxtTaskUpdate)
			require.True(t, ok)
			assert.Equal(t, orig.Lines, ret.Lines)
		case *TaskStartTime:
			ret, ok := retrieved.(*TaskStartTime)
			require.True(t, ok)
			assert.Equal(t, orig.TaskID, ret.TaskID)
			assert.Equal(t, orig.StartTime, ret.StartTime)
		case *TaskEndTime:
			ret, ok := retrieved.(*TaskEndTime)
			require.True(t, ok)
			assert.Equal(t, orig.TaskID, ret.TaskID)
			assert.Equal(t, orig.EndTime, ret.EndTime)
		}
	}

	// Test reopening the store
	newStore, err := NewAppendLogEventStore(filePath)
	require.NoError(t, err)

	// Retrieve events from the new store instance
	retrievedEvents, err = newStore.GetEvents()
	require.NoError(t, err)
	require.Len(t, retrievedEvents, len(events))

	// Check file content
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"type":"TodoTxtTaskUpdate"`)
	assert.Contains(t, string(content), `"type":"TaskStartTime"`)
	assert.Contains(t, string(content), `"type":"TaskEndTime"`)
}

func TestNewRepositoryWithAppendLog(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "events.jsonl")

	// Create a new append log store
	store, err := NewAppendLogEventStore(filePath)
	require.NoError(t, err)

	// Create a repository with the append log store
	repo := NewRepositoryWithEventStore(store)
	require.NotNil(t, repo)

	// Create a test task
	task := Task{
		Todo:     "Test task",
		Priority: "A",
		Projects: []string{"project"},
		Contexts: []string{"context"},
	}

	// Save the task via repository
	taskID, err := repo.CreateTask(task)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, taskID)

	// Retrieve the task
	retrievedTask, err := repo.GetTask(taskID)
	require.NoError(t, err)
	assert.Equal(t, task.Todo, retrievedTask.Todo)
	assert.Equal(t, task.Priority, retrievedTask.Priority)

	// Create a new repository with the same append log
	newStore, err := NewAppendLogEventStore(filePath)
	require.NoError(t, err)
	newRepo := NewRepositoryWithEventStore(newStore)
	require.NotNil(t, newRepo)

	// Rebuild the state
	err = newRepo.RebuildState()
	require.NoError(t, err)

	// Check if the task is in the new repository
	retrievedTask, err = newRepo.GetTask(taskID)
	require.NoError(t, err)
	assert.Equal(t, task.Todo, retrievedTask.Todo)
	assert.Equal(t, task.Priority, retrievedTask.Priority)
}