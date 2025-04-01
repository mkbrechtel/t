package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTodoUpdate tests the 'todo update' command ensures proper task properties
func TestTodoUpdate(t *testing.T) {
	cleanup, configFile, eventStoreFile := setupTestEnv(t)
	defer cleanup()

	// Get the path to the todo.txt file
	todoFilePath := filepath.Join(filepath.Dir(configFile), "todo.txt")
	
	// Count tasks before update
	originalTasks := parseTodoFile(t, todoFilePath)
	
	// Run the todo update command
	_, stderr, err := runT5Command(t, "--config", configFile, "todo", "update")
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr)
	}

	// Verify the event store file was created
	_, err = os.Stat(eventStoreFile)
	assert.NoError(t, err, "Event store file should exist")

	// Read the updated todo.txt file
	updatedTasks := parseTodoFile(t, todoFilePath)
	
	// Number of tasks should remain the same
	assert.Equal(t, len(originalTasks), len(updatedTasks), 
		"Number of tasks should not change")
	
	// Check that all completed tasks have completion dates
	for _, task := range updatedTasks {
		if task.Completed {
			assert.NotEmpty(t, task.CompletionDate, 
				"Completed task should have a completion date")
		}
	}

	// Check that all tasks have IDs
	for _, task := range updatedTasks {
		assert.NotEmpty(t, task.ID, "All tasks should have an ID")
	}

	// Verify the event store contains one event
	assert.Equal(t, 1, countEventsInFile(t, eventStoreFile),
		"Event store should contain one event after update")
	
	// Verify the event is valid
	validateEventStore(t, eventStoreFile)

	t.Logf("Todo update command completed successfully")
}

// TestTaskLifecycle tests the complete lifecycle of a task
func TestTaskLifecycle(t *testing.T) {
	cleanup, configFile, eventStoreFile := setupTestEnv(t)
	defer cleanup()

	// Initialize the event store
	_, stderr, err := runT5Command(t, "--config", configFile, "todo", "update")
	if err != nil {
		t.Fatalf("Update command failed: %v\nStderr: %s", err, stderr)
	}

	// Add a new task with specific attributes
	todoFilePath := filepath.Join(filepath.Dir(configFile), "todo.txt")
	content, err := os.ReadFile(todoFilePath)
	require.NoError(t, err, "Should be able to read the todo.txt file")
	
	today := time.Now().Format("2006-01-02")
	taskText := "(A) " + today + " Important test task +project @context"
	newContent := string(content) + "\n" + taskText
	
	err = os.WriteFile(todoFilePath, []byte(newContent), 0644)
	require.NoError(t, err, "Should be able to write to the todo.txt file")

	// Run update to process the new task
	_, stderr, err = runT5Command(t, "--config", configFile, "todo", "update")
	if err != nil {
		t.Fatalf("Update command failed: %v\nStderr: %s", err, stderr)
	}

	// Verify the task was added with correct attributes
	tasks := parseTodoFile(t, todoFilePath)
	var taskID string
	var foundTask bool
	
	for _, task := range tasks {
		if strings.Contains(task.Raw, "Important test task") {
			foundTask = true
			assert.Equal(t, "A", task.Priority, "Task should have priority A")
			assert.Equal(t, today, task.CreatedAt, "Task should have today's date")
			assert.NotEmpty(t, task.ID, "Task should have an ID")
			taskID = task.ID
			break
		}
	}
	assert.True(t, foundTask, "Should find the added task")

	// Now mark the task as completed
	content, err = os.ReadFile(todoFilePath)
	require.NoError(t, err, "Should be able to read the todo.txt file")
	
	lines := strings.Split(string(content), "\n")
	var newLines []string
	
	for _, line := range lines {
		if strings.Contains(line, "Important test task") && strings.Contains(line, taskID) {
			// Mark as completed but don't add completion date (should be added by update)
			newLines = append(newLines, "x "+line)
		} else {
			newLines = append(newLines, line)
		}
	}
	
	err = os.WriteFile(todoFilePath, []byte(strings.Join(newLines, "\n")), 0644)
	require.NoError(t, err, "Should be able to write to the todo.txt file")

	// Run update to process the completion with explicit flag for enforce-completion-date
	_, stderr, err = runT5Command(t, "--config", configFile, "--enforce-completion-date=true", "todo", "update")
	if err != nil {
		t.Fatalf("Update command failed: %v\nStderr: %s", err, stderr)
	}

	// Verify the task is now completed with a completion date
	tasks = parseTodoFile(t, todoFilePath)
	foundTask = false
	
	for _, task := range tasks {
		if strings.Contains(task.Raw, "Important test task") && strings.Contains(task.Raw, taskID) {
			foundTask = true
			assert.True(t, task.Completed, "Task should be marked as completed")
			assert.NotEmpty(t, task.CompletionDate, "Completed task should have a completion date")
			break
		}
	}
	assert.True(t, foundTask, "Should find the completed task")

	// Verify we now have 3 events (initial, add task, complete task)
	assert.Equal(t, 3, countEventsInFile(t, eventStoreFile),
		"Should have three events after task lifecycle")

	t.Logf("Task lifecycle test completed successfully")
}

// TestPriorityHandling tests task priority handling
func TestPriorityHandling(t *testing.T) {
	cleanup, configFile, _ := setupTestEnv(t)
	defer cleanup()

	// Initialize the event store
	_, stderr, err := runT5Command(t, "--config", configFile, "todo", "update")
	if err != nil {
		t.Fatalf("Update command failed: %v\nStderr: %s", err, stderr)
	}

	// Add tasks with different priorities
	todoFilePath := filepath.Join(filepath.Dir(configFile), "todo.txt")
	content, err := os.ReadFile(todoFilePath)
	require.NoError(t, err, "Should be able to read the todo.txt file")
	
	today := time.Now().Format("2006-01-02")
	
	// Add high, medium, and no priority tasks
	newTasks := []string{
		"(A) " + today + " High priority task",
		"(B) " + today + " Medium priority task",
		today + " No priority task",
	}
	
	newContent := string(content) + "\n" + strings.Join(newTasks, "\n")
	err = os.WriteFile(todoFilePath, []byte(newContent), 0644)
	require.NoError(t, err, "Should be able to write to the todo.txt file")

	// Run update to process the new tasks
	_, stderr, err = runT5Command(t, "--config", configFile, "todo", "update")
	if err != nil {
		t.Fatalf("Update command failed: %v\nStderr: %s", err, stderr)
	}

	// Verify tasks with priorities in the todo.txt file
	tasks := parseTodoFile(t, todoFilePath)
	
	var highPriorityFound, mediumPriorityFound, noPriorityFound bool
	
	for _, task := range tasks {
		if strings.Contains(task.Raw, "High priority task") {
			highPriorityFound = true
			assert.Equal(t, "A", task.Priority, "High priority task should have priority A")
		} else if strings.Contains(task.Raw, "Medium priority task") {
			mediumPriorityFound = true
			assert.Equal(t, "B", task.Priority, "Medium priority task should have priority B")
		} else if strings.Contains(task.Raw, "No priority task") {
			noPriorityFound = true
			assert.Empty(t, task.Priority, "No priority task should have empty priority")
		}
	}
	
	assert.True(t, highPriorityFound, "Should find the high priority task")
	assert.True(t, mediumPriorityFound, "Should find the medium priority task")
	assert.True(t, noPriorityFound, "Should find the no priority task")

	// Read the todo.txt file again to verify tasks still exist
	todoContent, err := os.ReadFile(todoFilePath)
	require.NoError(t, err, "Should be able to read todo.txt file")
	todoStr := string(todoContent)
	
	// Verify tasks with priorities in the todo.txt file
	assert.Contains(t, todoStr, "High priority task", "High priority task should be in the todo.txt file")
	assert.Contains(t, todoStr, "Medium priority task", "Medium priority task should be in the todo.txt file")
	assert.Contains(t, todoStr, "No priority task", "No priority task should be in the todo.txt file")

	t.Logf("Priority handling test completed successfully")
}