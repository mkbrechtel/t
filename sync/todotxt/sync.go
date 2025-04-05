package todo

import (
	"fmt"
	"os"
	"time"
	
	todo "github.com/1set/todotxt"
	
	"t5.mkbrechtel.dev/core"
	"t5.mkbrechtel.dev/utils"
)

// SyncResult contains statistics about the sync operation
type SyncResult struct {
	Added      int
	Updated    int
	Skipped    int
	Conflicts  int
	FromTodoTxt int
	ToTodoTxt   int
}

// SyncWithRepository synchronizes a todo.txt file with the repository's event store
func SyncWithRepository(repo *core.Repository, todoFilePath string) (*SyncResult, error) {
	result := &SyncResult{}
	
	// Get current application state
	appState := repo.GetAppState()
	
	// Read the todo.txt file
	todoList, err := ReadTodoFile(todoFilePath)
	if err != nil {
		// If file doesn't exist, create an empty one
		if os.IsNotExist(err) {
			// Create new tasks from the repository's state
			return createTodoFile(repo, todoFilePath)
		}
		return nil, err
	}
	
	// Map of UUID to todo.txt tasks for easier lookup
	todoTasksByUUID := make(map[string]*todo.Task)
	tasksWithoutUUID := make([]*todo.Task, 0)
	
	// First pass - build lookup maps
	for i := range todoList {
		task := &todoList[i]
		
		// Check if task has UUID
		if uuid, exists := task.AdditionalTags["uuid"]; exists && uuid != "" {
			todoTasksByUUID[uuid] = task
		} else {
			tasksWithoutUUID = append(tasksWithoutUUID, task)
		}
	}
	
	// Second pass - check for updates from todo.txt to repo
	for uuidStr, todoTask := range todoTasksByUUID {
		id, err := utils.DecodeUUID(uuidStr)
		if err == nil {
			if repoTask, exists := appState.Tasks[id]; exists {
				// Compare tasks, detect changes
				if isTaskModified(todoTask, repoTask) {
					// If both have different changes, detect conflict
					if hasConflict(todoTask, repoTask) {
						result.Conflicts++
						// Log the conflict - in a real application, you might want a more sophisticated
						// conflict resolution strategy, but for now we just report it
						fmt.Printf("Conflict detected for task with UUID %s\n", uuidStr)
					} else {
						// Update repository from todo.txt changes
						todoTxtContent := todoTask.String()
						event := core.NewTodoTxtTaskUpdate(todoTxtContent, todoFilePath)
						err := repo.SaveEvent(event)
						if err != nil {
							return nil, fmt.Errorf("failed to save todo.txt task update: %w", err)
						}
						result.FromTodoTxt++
					}
				}
			} 
		} else {
			// Task exists in todo.txt but not in repo or has invalid UUID - add it
			todoTxtContent := todoTask.String()
			event := core.NewTodoTxtTaskUpdate(todoTxtContent, todoFilePath)
			err := repo.SaveEvent(event)
			if err != nil {
				return nil, fmt.Errorf("failed to save todo.txt task update: %w", err)
			}
			result.Added++
		}
	}
	
	// Handle tasks without UUIDs - create them in the repository
	for _, task := range tasksWithoutUUID {
		todoTxtContent := task.String()
		event := core.NewTodoTxtTaskUpdate(todoTxtContent, todoFilePath)
		err := repo.SaveEvent(event)
		if err != nil {
			return nil, fmt.Errorf("failed to save todo.txt task without UUID: %w", err)
		}
		result.Added++
	}
	
	// Third pass - check for tasks in repo not in todo.txt
	var tasksToWrite []todo.Task
	for id, repoTask := range appState.Tasks {
		// Skip tasks that aren't from this todo.txt file
		if repoTask.Source != todoFilePath && repoTask.Source != "todo.txt" {
			continue
		}
		
		// Check if task exists in todo.txt
		uuidStr := id.String()
		if _, exists := todoTasksByUUID[uuidStr]; !exists {
			// Task exists in repo but not in todo.txt - add it to todo.txt
			todoTask := repoTaskToTodoTask(repoTask)
			tasksToWrite = append(tasksToWrite, todoTask)
			result.ToTodoTxt++
		}
	}
	
	// Write any changes back to todo.txt if needed
	if result.ToTodoTxt > 0 {
		// Add the existing tasks
		for _, task := range todoList {
			tasksToWrite = append(tasksToWrite, task)
		}
		
		err = WriteTodoFile(tasksToWrite, todoFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to write tasks to todo.txt: %w", err)
		}
	}
	
	return result, nil
}

// createTodoFile creates a new todo.txt file from the repository's state
func createTodoFile(repo *core.Repository, todoFilePath string) (*SyncResult, error) {
	result := &SyncResult{}
	
	// Get current application state
	appState := repo.GetAppState()
	
	// Create tasks list from repository
	var tasksToWrite []todo.Task
	for _, repoTask := range appState.Tasks {
		// Skip tasks that are marked as deleted
		if _, isDeleted := repoTask.AdditionalTags["deleted"]; isDeleted {
			continue
		}
		
		todoTask := repoTaskToTodoTask(repoTask)
		tasksToWrite = append(tasksToWrite, todoTask)
		result.ToTodoTxt++
	}
	
	// Write tasks to todo.txt
	err := WriteTodoFile(tasksToWrite, todoFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo.txt file: %w", err)
	}
	
	return result, nil
}

// isTaskModified checks if a todo.txt task has been modified compared to the repository task
func isTaskModified(todoTask *todo.Task, repoTask core.Task) bool {
	// Compare basic properties
	if todoTask.Todo != repoTask.Todo {
		return true
	}
	if todoTask.Priority != repoTask.Priority {
		return true
	}
	if todoTask.Completed != repoTask.Completed {
		return true
	}
	
	// Compare projects (order independent)
	if !stringSlicesEqual(todoTask.Projects, repoTask.Projects) {
		return true
	}
	
	// Compare contexts (order independent)
	if !stringSlicesEqual(todoTask.Contexts, repoTask.Contexts) {
		return true
	}
	
	// Compare dates
	if !datesEqual(todoTask.CreatedDate, repoTask.CreatedDate) {
		return true
	}
	if !datesEqual(todoTask.DueDate, repoTask.DueDate) {
		return true
	}
	if !datesEqual(todoTask.CompletedDate, repoTask.CompletedDate) {
		return true
	}
	
	return false
}

// hasConflict checks if there are conflicting changes between a todo.txt task and repo task
func hasConflict(todoTask *todo.Task, repoTask core.Task) bool {
	// This is a simplified conflict detection
	// In a real application, you might want to track when each field was last modified
	// and use that to determine conflicts
	
	// For now, we'll use the modified timestamp if available
	todoModified := todoTask.AdditionalTags["modified"]
	repoModified := repoTask.AdditionalTags["modified"]
	
	if todoModified != "" && repoModified != "" && todoModified != repoModified {
		// Parse timestamps
		todoTime, todoErr := time.Parse(time.RFC3339, todoModified)
		repoTime, repoErr := time.Parse(time.RFC3339, repoModified)
		
		// If we can parse both timestamps, compare them
		if todoErr == nil && repoErr == nil {
			// If timestamps are very close (within 1 second), consider it a conflict
			diff := todoTime.Sub(repoTime)
			if diff < 1*time.Second && diff > -1*time.Second {
				return true
			}
			// Otherwise, the more recent change wins (already handled by our sync logic)
			return false
		}
	}
	
	// If we can't determine based on timestamps, we'll be cautious
	// and flag a conflict if both have substantive changes
	// This could be improved with more sophisticated conflict detection
	return false
}

// repoTaskToTodoTask converts a repository Task to a todo.txt Task
func repoTaskToTodoTask(task core.Task) todo.Task {
	// Create a todo.txt task
	todoTask := todo.Task{
		Todo:           task.Todo,
		Priority:       task.Priority,
		Projects:       task.Projects,
		Contexts:       task.Contexts,
		Completed:      task.Completed,
		CreatedDate:    task.CreatedDate,
		DueDate:        task.DueDate,
		CompletedDate:  task.CompletedDate,
	}
	
	// Copy additional tags
	todoTask.AdditionalTags = make(map[string]string)
	for k, v := range task.AdditionalTags {
		todoTask.AdditionalTags[k] = v
	}
	
	// Ensure UUID is set with long format
	todoTask.AdditionalTags["uuid"] = utils.LongEncodeUUID(task.ID)
	
	// Set modified timestamp
	todoTask.AdditionalTags["modified"] = time.Now().Format(time.RFC3339)
	
	return todoTask
}

// stringSlicesEqual compares two string slices for equality, ignoring order
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	
	// Create maps for efficient lookup
	aMap := make(map[string]struct{}, len(a))
	bMap := make(map[string]struct{}, len(b))
	
	for _, s := range a {
		aMap[s] = struct{}{}
	}
	
	for _, s := range b {
		bMap[s] = struct{}{}
	}
	
	// Compare maps
	if len(aMap) != len(bMap) {
		return false
	}
	
	for s := range aMap {
		if _, ok := bMap[s]; !ok {
			return false
		}
	}
	
	return true
}

// datesEqual compares two dates for equality, ignoring time parts if desired
func datesEqual(a, b time.Time) bool {
	if a.IsZero() && b.IsZero() {
		return true
	}
	if a.IsZero() != b.IsZero() {
		return false
	}
	
	// Compare only the date parts
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()
	
	return aYear == bYear && aMonth == bMonth && aDay == bDay
}