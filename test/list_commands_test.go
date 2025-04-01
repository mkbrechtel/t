package test

import (
	"path/filepath"
	"strings"
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
	stdout, stderr, err := runT5Command(t, "--config", configFile, "list")
	if err != nil {
		t.Fatalf("List command failed: %v\nStderr: %s", err, stderr)
	}

	// Verify that we have some output from the list command
	assert.Contains(t, stdout, "Found", "List output should contain 'Found'")
	assert.Contains(t, stdout, "tasks", "List output should contain 'tasks'")

	// Verify that the number of tasks matches what's in todo.txt
	// (subtract header line and any empty lines)
	nonEmptyLines := 0
	for _, line := range strings.Split(stdout, "\n") {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines++
		}
	}
	// The output should have a header line "Found X tasks:" plus one line per task
	assert.Equal(t, len(tasks)+1, nonEmptyLines, 
		"Number of tasks in list output should match todo.txt")

	// Check that all task descriptions appear in the output
	for _, task := range tasks {
		assert.Contains(t, stdout, task.Todo, 
			"Each task description should appear in the list output")
	}

	t.Logf("List command completed successfully")
}