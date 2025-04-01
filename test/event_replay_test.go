package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEventReplay tests the event replay system by modifying tasks and verifying state
func TestEventReplay(t *testing.T) {
	cleanup, configFile, eventStoreFile := setupTestEnv(t)
	defer cleanup()

	// Run update to initialize the event store
	_, stderr, err := runT5Command(t, "--config", configFile, "todo", "update")
	if err != nil {
		t.Fatalf("Update command failed: %v\nStderr: %s", err, stderr)
	}

	// Get the initial task count from list output
	stdout, stderr, err := runT5Command(t, "--config", configFile, "list")
	if err != nil {
		t.Fatalf("List command failed: %v\nStderr: %s", err, stderr)
	}
	initialOutput := stdout
	
	// Count the initial number of events
	initialEventCount := countEventsInFile(t, eventStoreFile)
	assert.Equal(t, 1, initialEventCount, "Should have one event after initial update")

	// Add a new task to todo.txt
	todoFilePath := filepath.Join(filepath.Dir(configFile), "todo.txt")
	content, err := os.ReadFile(todoFilePath)
	require.NoError(t, err, "Should be able to read the todo.txt file")
	
	newTaskText := "Test task for event replay"
	newContent := string(content) + "\n" + newTaskText
	
	err = os.WriteFile(todoFilePath, []byte(newContent), 0644)
	require.NoError(t, err, "Should be able to write to the todo.txt file")

	// Run update again to create a new event
	_, stderr, err = runT5Command(t, "--config", configFile, "todo", "update")
	if err != nil {
		t.Fatalf("Update command failed: %v\nStderr: %s", err, stderr)
	}

	// Count the updated number of events
	updatedEventCount := countEventsInFile(t, eventStoreFile)
	assert.Equal(t, 2, updatedEventCount, "Should have two events after second update")

	// List the tasks again to verify the new task was added
	stdout, stderr, err = runT5Command(t, "--config", configFile, "list")
	if err != nil {
		t.Fatalf("List command failed: %v\nStderr: %s", err, stderr)
	}
	updatedOutput := stdout
	
	// Verify that the new task appears in the list output
	assert.Contains(t, updatedOutput, newTaskText, 
		"The new test task should be in the list output")
	
	// Verify that there are more lines in the updated output
	assert.Greater(t, strings.Count(updatedOutput, "\n"), 
		strings.Count(initialOutput, "\n"), 
		"Updated output should have more lines")

	// Validate all events in the event store
	validateEventStore(t, eventStoreFile)

	t.Logf("Event replay test completed successfully")
}