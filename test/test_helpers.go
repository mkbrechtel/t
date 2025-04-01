package test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"t5.mkbrechtel.dev/cmd"
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
	workFile := filepath.Join(cwd, "work.txt")
	eventStoreFile := filepath.Join(cwd, "events.jsonl")

	// Remove the event store file if it exists to start fresh
	os.Remove(eventStoreFile)

	// Make backup of the original todo.txt file
	backupTodoFile := todoFile + ".bak"
	todoContent, err := os.ReadFile(todoFile)
	if err != nil {
		t.Fatalf("Failed to read todo.txt file: %v", err)
	}
	err = os.WriteFile(backupTodoFile, todoContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create backup of todo.txt: %v", err)
	}

	// Create work.txt file if it doesn't exist
	if _, err := os.Stat(workFile); os.IsNotExist(err) {
		err = os.WriteFile(workFile, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create work.txt file: %v", err)
		}
	} else {
		// Make backup of work.txt if it exists
		backupWorkFile := workFile + ".bak"
		workContent, err := os.ReadFile(workFile)
		if err != nil {
			t.Fatalf("Failed to read work.txt file: %v", err)
		}
		err = os.WriteFile(backupWorkFile, workContent, 0644)
		if err != nil {
			t.Fatalf("Failed to create backup of work.txt: %v", err)
		}
	}

	// Return cleanup function, config file path, and event store path
	cleanup := func() {
		os.Remove(eventStoreFile)
		
		// Restore the original todo.txt file
		restoreContent, err := os.ReadFile(backupTodoFile)
		if err != nil {
			t.Logf("Failed to read backup of todo.txt: %v", err)
		} else {
			err = os.WriteFile(todoFile, restoreContent, 0644)
			if err != nil {
				t.Logf("Failed to restore todo.txt: %v", err)
			}
		}
		os.Remove(backupTodoFile)
		
		// Restore the original work.txt file if it had a backup
		backupWorkFile := workFile + ".bak"
		if _, err := os.Stat(backupWorkFile); !os.IsNotExist(err) {
			restoreWorkContent, err := os.ReadFile(backupWorkFile)
			if err != nil {
				t.Logf("Failed to read backup of work.txt: %v", err)
			} else {
				err = os.WriteFile(workFile, restoreWorkContent, 0644)
				if err != nil {
					t.Logf("Failed to restore work.txt: %v", err)
				}
			}
			os.Remove(backupWorkFile)
		}
	}

	return cleanup, configFile, eventStoreFile
}

// runT5Command runs the t5 command with the given arguments and returns stdout, stderr, and error
func runT5Command(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	// Capture stdout and stderr
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	
	// Prepare arguments with the program name as first argument (like os.Args)
	fullArgs := append([]string{"t5"}, args...)
	
	// Execute the command with our captured output
	err := cmd.ExecuteWithArgs(fullArgs, stdout, stderr)
	
	// Wait a moment to ensure file operations complete
	// This helps with event store file creation and visibility
	time.Sleep(100 * time.Millisecond)
	
	// For config test to work
	stdoutStr := stdout.String()
	if len(args) > 0 && args[len(args)-1] == "config" {
		t.Logf("Config output: %s", stdoutStr)
	}
	
	return stdoutStr, stderr.String(), err
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
		assert.Contains(t, event, "type", "Event should have a type")
		assert.Contains(t, event, "timestamp", "Event should have a timestamp")
		assert.Contains(t, event, "data", "Event should have data")
	}
}