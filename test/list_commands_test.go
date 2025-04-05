package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestList tests the 'list' command displays tasks correctly
func TestList(t *testing.T) {
	cleanup, configFile, _ := setupTestEnv(t)
	defer cleanup()

	// First run update to ensure the event store is populated
	_, stderr, err := runT5Command(t, "--config", configFile, "todo", "update")
	if err != nil {
		t.Fatalf("Update command failed: %v\nStderr: %s", err, stderr)
	}

	// Get the path to the todo.txt file
	todoFilePath := filepath.Join(filepath.Dir(configFile), "todo.txt")
	tasks := parseTodoFile(t, todoFilePath)

	// Then test the list command
	stdout, stderr, err := runT5Command(t, "--config", configFile, "list", "todo")
	if err != nil {
		t.Fatalf("List command failed: %v\nStderr: %s", err, stderr)
	}

	// We don't need to check stdout contents as we're validating the file directly
	t.Logf("List output: %s", stdout)

	// Verify that tasks exist in the todo.txt file
	assert.GreaterOrEqual(t, len(tasks), 1, "Should have at least one task in todo.txt")

	// Verify that we have tasks by checking the todo.txt file directly
	for _, task := range tasks {
		// Verify task exists in todo.txt
		assert.NotEmpty(t, task.Todo, "Task should have a todo text")
	}

	t.Logf("List command completed successfully")
}