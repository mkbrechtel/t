package core

import (
	"strings"
	"testing"
	"time"

	todo "github.com/1set/todotxt"
	uuid "github.com/gofrs/uuid/v5"
)

func TestTodoTxtTaskUpdate_Apply(t *testing.T) {
	// Create a sample todo.txt content
	todoContent := `(A) 2023-05-15 Implement the core module +t5 @dev due:2023-06-01
x 2023-06-10 2023-05-20 Fix todo parsing bug +t5 @bugfix uuid:be52544a-22c3-7d86-b3a6-47f84e3fe9b2
Buy groceries @personal +shopping due:2023-05-25`

	// Create a TodoTxtTaskUpdate event
	event := TodoTxtTaskUpdate{
		Lines: todoContent,
	}

	// Create a new app state
	state := &AppState{
		Tasks:    make(map[uuid.UUID]Task),
		Projects: make(map[string]Project),
	}

	// Apply the event
	err := event.apply(state)
	if err != nil {
		t.Fatalf("TodoTxtTaskUpdate.Apply failed: %v", err)
	}

	// Test 1: Verify that all tasks were parsed and added to the state
	if len(state.Tasks) != 3 {
		t.Errorf("Expected 3 tasks in state, got %d", len(state.Tasks))
	}

	// Test 2: Verify that projects were created
	if len(state.Projects) != 2 {
		t.Errorf("Expected 2 projects in state, got %d", len(state.Projects))
	}

	// Verify the projects are correct
	if _, exists := state.Projects["t5"]; !exists {
		t.Errorf("Expected project 't5' to exist")
	}
	if _, exists := state.Projects["shopping"]; !exists {
		t.Errorf("Expected project 'shopping' to exist")
	}

	// Test 3: Verify that we can find a task with known UUID
	knownTaskID := uuid.Must(uuid.FromString("be52544a-22c3-7d86-b3a6-47f84e3fe9b2"))
	task, exists := state.Tasks[knownTaskID]
	if !exists {
		t.Errorf("Expected to find task with UUID be52544a-22c3-7d86-b3a6-47f84e3fe9b2")
	} else {
		// Verify task properties
		if !task.Completed {
			t.Errorf("Expected task to be completed")
		}
		if !containsString(task.Projects, "t5") {
			t.Errorf("Expected task to have project 't5'")
		}
		if !containsString(task.Contexts, "bugfix") {
			t.Errorf("Expected task to have context 'bugfix'")
		}
		if !task.CreatedDate.Equal(mustParseDate("2023-05-20")) {
			t.Errorf("Expected creation date 2023-05-20, got %v", task.CreatedDate)
		}
		if !task.CompletedDate.Equal(mustParseDate("2023-06-10")) {
			t.Errorf("Expected completion date 2023-06-10, got %v", task.CompletedDate)
		}
	}

	// Test 4: Verify that new tasks have generated UUIDs
	foundNewlyGeneratedUUID := false
	for id, task := range state.Tasks {
		if id != knownTaskID {
			foundNewlyGeneratedUUID = true
			if id == uuid.Nil {
				t.Errorf("Expected generated UUID to not be nil")
			}
			// Check if the UUID is stored in the task's additional tags
			if uuidStr, ok := task.AdditionalTags["uuid"]; !ok || uuidStr != id.String() {
				t.Errorf("Expected uuid tag in task to match the task ID")
			}
		}
	}
	if !foundNewlyGeneratedUUID {
		t.Errorf("Expected to find at least one task with a newly generated UUID")
	}

	// Test 5: Verify priority is correctly set
	for _, task := range state.Tasks {
		if strings.Contains(task.Original, "(A)") && task.Priority != "A" {
			t.Errorf("Expected task with (A) to have priority A, got %s", task.Priority)
		}
	}

	// Test 6: Verify due dates are correctly set
	for _, task := range state.Tasks {
		if strings.Contains(task.Original, "due:2023-06-01") {
			if !task.DueDate.Equal(mustParseDate("2023-06-01")) {
				t.Errorf("Expected due date 2023-06-01, got %v", task.DueDate)
			}
		}
		if strings.Contains(task.Original, "due:2023-05-25") {
			if !task.DueDate.Equal(mustParseDate("2023-05-25")) {
				t.Errorf("Expected due date 2023-05-25, got %v", task.DueDate)
			}
		}
	}
}

// Helper function to check if a string slice contains a string
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Helper function to parse date strings
func mustParseDate(dateStr string) time.Time {
	t, err := time.ParseInLocation(todo.DateLayout, dateStr, time.Local)
	if err != nil {
		panic("Invalid date format: " + dateStr)
	}
	return t
}