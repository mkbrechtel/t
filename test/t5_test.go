package test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Task represents a todo.txt task for testing purposes
type Task struct {
	Raw           string
	Todo          string
	Priority      string
	Completed     bool
	CompletionDate string
	CreatedAt     string
	ID            string
}

// setupTestEnv prepares the test environment by setting up the necessary directories
// and files for testing t5 commands
func setupTestEnv(t *testing.T) (func(), string, string) {
	t.Helper()

	// Get current directory as base for test files
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Use the test config file
	configFile := filepath.Join(cwd, "t5config.yaml")
	todoFile := filepath.Join(cwd, "todo.txt")
	eventStoreFile := filepath.Join(cwd, "events.jsonl")

	// Remove the event store file if it exists to start fresh
	os.Remove(eventStoreFile)

	// Make backup of the original todo file
	backupTodoFile := todoFile + ".bak"
	todoContent, err := os.ReadFile(todoFile)
	if err != nil {
		t.Fatalf("Failed to read todo.txt file: %v", err)
	}
	err = os.WriteFile(backupTodoFile, todoContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create backup of todo.txt: %v", err)
	}

	// Return cleanup function, config file path, and event store path
	cleanup := func() {
		os.Remove(eventStoreFile)
		// Restore the original todo.txt file
		restoreContent, err := os.ReadFile(backupTodoFile)
		if err != nil {
			t.Logf("Failed to read backup of todo.txt: %v", err)
			return
		}
		err = os.WriteFile(todoFile, restoreContent, 0644)
		if err != nil {
			t.Logf("Failed to restore todo.txt: %v", err)
		}
		os.Remove(backupTodoFile)
	}

	return cleanup, configFile, eventStoreFile
}

// runT5Command runs the t5 command with the given arguments and returns stdout, stderr, and error
func runT5Command(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	// Find the t5 binary relative to the test directory
	// Assuming the binary is in the project root
	binPath, err := filepath.Abs(filepath.Join("..", "t5"))
	if err != nil {
		t.Fatalf("Failed to get absolute path for t5 binary: %v", err)
	}
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Fatalf("t5 binary not found at %s. Did you build it?", binPath)
	}

	cmd := exec.Command(binPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	return stdout.String(), stderr.String(), err
}

// parseTodoFile parses a todo.txt file and returns a slice of Task structs
func parseTodoFile(t *testing.T, filePath string) []Task {
	t.Helper()
	
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read the todo.txt file")
	
	lines := strings.Split(string(content), "\n")
	tasks := make([]Task, 0, len(lines))
	
	for _, line := range lines {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		
		task := Task{Raw: line}
		
		// Parse completion status
		if strings.HasPrefix(line, "x ") {
			task.Completed = true
			line = line[2:] // Remove "x " prefix
			
			// Parse completion date if present
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 && isDateFormat(parts[0]) {
				task.CompletionDate = parts[0]
				line = parts[1]
			}
		}
		
		// Parse priority if present
		if len(line) >= 3 && line[0] == '(' && line[2] == ')' && line[1] >= 'A' && line[1] <= 'Z' {
			task.Priority = string(line[1])
			line = strings.TrimSpace(line[3:])
		}
		
		// Parse creation date if present
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 && isDateFormat(parts[0]) {
			task.CreatedAt = parts[0]
			line = parts[1]
		}
		
		// Parse ID if present
		idMatch := regexp.MustCompile(`id:([^\s]+)`).FindStringSubmatch(line)
		if len(idMatch) > 1 {
			task.ID = idMatch[1]
			// Remove the ID from the todo text to avoid duplication
			line = strings.Replace(line, idMatch[0], "", 1)
			line = strings.TrimSpace(line)
		}
		
		// The rest is the todo text
		task.Todo = line
		
		tasks = append(tasks, task)
	}
	
	return tasks
}

// isDateFormat checks if a string matches the YYYY-MM-DD format
func isDateFormat(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

// countEventsInFile counts the number of events in the event store file
func countEventsInFile(t *testing.T, filePath string) int {
	t.Helper()
	
	content, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return 0
	}
	require.NoError(t, err, "Should be able to read the event store file")
	
	lines := strings.Split(string(content), "\n")
	eventCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			eventCount++
		}
	}
	
	return eventCount
}

// validateEventStore checks that events in the event store file are valid JSON
// and have the required fields
func validateEventStore(t *testing.T, filePath string) {
	t.Helper()
	
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read the event store file")
	
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		
		var event map[string]interface{}
		err := json.Unmarshal([]byte(line), &event)
		assert.NoError(t, err, "Each line should be valid JSON")
		
		// Verify it has the required fields for an event
		assert.Contains(t, event, "id", "Event should have an ID")
		assert.Contains(t, event, "type", "Event should have a type")
		assert.Contains(t, event, "timestamp", "Event should have a timestamp")
	}
}

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

// TestConfig tests the 'config' command displays configuration correctly
func TestConfig(t *testing.T) {
	cleanup, configFile, _ := setupTestEnv(t)
	defer cleanup()

	// Execute the config command
	stdout, stderr, err := runT5Command(t, "--config", configFile, "config")
	if err != nil {
		t.Fatalf("Config command failed: %v\nStderr: %s", err, stderr)
	}

	// Verify the config file exists
	_, err = os.Stat(configFile)
	require.NoError(t, err, "Config file should exist")

	// Verify that the config output contains expected sections from the config file
	assert.Contains(t, stdout, "todo:", "Config output should contain 'todo:' section")
	assert.Contains(t, stdout, "eventstore:", "Config output should contain 'eventstore:' section")
	
	// Verify specific configuration values
	assert.Contains(t, stdout, "file: ./todo.txt", 
		"Config output should contain the correct todo file path")
	assert.Contains(t, stdout, "prefershortids: true", 
		"Config output should contain the correct ID preference")
	assert.Contains(t, stdout, "enforcecompletiondate: true", 
		"Config output should contain the correct completion date enforcement")

	t.Logf("Config command completed successfully")
}

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

	// Run update to process the completion
	_, stderr, err = runT5Command(t, "--config", configFile, "todo", "update")
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

	// Verify tasks with priorities in list output
	stdout, stderr, err := runT5Command(t, "--config", configFile, "list")
	if err != nil {
		t.Fatalf("List command failed: %v\nStderr: %s", err, stderr)
	}
	
	// In list output, priorities should be displayed correctly
	assert.Regexp(t, regexp.MustCompile(`\[\s\] \(A\).*High priority task`), stdout, 
		"High priority task should be displayed with (A)")
	assert.Regexp(t, regexp.MustCompile(`\[\s\] \(B\).*Medium priority task`), stdout, 
		"Medium priority task should be displayed with (B)")
	assert.Regexp(t, regexp.MustCompile(`\[\s\] \(-\).*No priority task`), stdout, 
		"No priority task should be displayed with (-)")

	t.Logf("Priority handling test completed successfully")
}