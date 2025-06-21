package todo

import (
	"fmt"
	"os"

	todotxtlib "github.com/1set/todotxt"

	"t5.mkbrechtel.dev/t5/core"
	"t5.mkbrechtel.dev/t5/sync"
	"t5.mkbrechtel.dev/t5/utils"
)

// TodoTxtProvider is a sync provider for todo.txt files
type TodoTxtProvider struct {
	filePath string
}

// NewTodoTxtProvider creates a new todo.txt sync provider
func NewTodoTxtProvider(filePath string) *TodoTxtProvider {
	return &TodoTxtProvider{
		filePath: filePath,
	}
}

// Name returns the name of the sync provider
func (p *TodoTxtProvider) Name() string {
	return "Todo.txt"
}

// Description returns a description of the sync provider
func (p *TodoTxtProvider) Description() string {
	return fmt.Sprintf("Todo.txt file sync provider for %s", p.filePath)
}

// SupportedDirections returns the sync directions supported by this provider
func (p *TodoTxtProvider) SupportedDirections() []sync.SyncDirection {
	return []sync.SyncDirection{
		sync.SyncDirectionImport,
		sync.SyncDirectionExport,
		sync.SyncDirectionBoth,
	}
}

// Import imports tasks from the todo.txt file into the t5 repository
func (p *TodoTxtProvider) Import(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	result := &sync.SyncResult{}

	// Read the todo.txt file
	todoList, err := ReadTodoFile(p.filePath)
	if err != nil {
		// If file doesn't exist, nothing to import
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, err
	}

	// Get current application state
	appState := repo.GetAppState()

	// First pass - collect tasks without UUIDs and with UUIDs
	todoTasksByUUID := make(map[string]*todotxtlib.Task)
	tasksWithoutUUID := make([]*todotxtlib.Task, 0)

	for i := range todoList {
		task := &todoList[i]

		// Apply filter if provided
		if filter != nil {
			repoTask := todoTaskToRepoTask(*task)
			if !filter.FilterTask(repoTask) {
				result.Skipped++
				continue
			}
		}

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
						// Log the conflict
						fmt.Printf("Conflict detected for task with UUID %s\n", uuidStr)
					} else {
						// Convert to repo task
						repoTask := todoTaskToRepoTask(*todoTask)

						// Apply modifier if provided
						if modifier != nil {
							repoTask = modifier.ModifyTask(repoTask)
						}

						// Update repository from todo.txt changes
						todoTxtContent := todoTaskToString(repoTask)
						event := core.NewTodoTxtTaskUpdate(todoTxtContent, p.filePath)
						err := repo.SaveEvent(event)
						if err != nil {
							return nil, fmt.Errorf("failed to save todo.txt task update: %w", err)
						}
						result.FromSource++
					}
				} else {
					result.Skipped++
				}
			} else {
				// Convert to repo task
				repoTask := todoTaskToRepoTask(*todoTask)

				// Apply modifier if provided
				if modifier != nil {
					repoTask = modifier.ModifyTask(repoTask)
				}

				// Task exists in todo.txt but not in repo or has invalid UUID - add it
				todoTxtContent := todoTaskToString(repoTask)
				event := core.NewTodoTxtTaskUpdate(todoTxtContent, p.filePath)
				err := repo.SaveEvent(event)
				if err != nil {
					return nil, fmt.Errorf("failed to save todo.txt task update: %w", err)
				}
				result.Added++
			}
		}
	}

	// Handle tasks without UUIDs - create them in the repository
	for _, task := range tasksWithoutUUID {
		// Convert to repo task
		repoTask := todoTaskToRepoTask(*task)

		// Apply modifier if provided
		if modifier != nil {
			repoTask = modifier.ModifyTask(repoTask)
		}

		// Create new task in repository
		todoTxtContent := todoTaskToString(repoTask)
		event := core.NewTodoTxtTaskUpdate(todoTxtContent, p.filePath)
		err := repo.SaveEvent(event)
		if err != nil {
			return nil, fmt.Errorf("failed to save todo.txt task without UUID: %w", err)
		}
		result.Added++
	}

	return result, nil
}

// Export exports tasks from the t5 repository to the todo.txt file
func (p *TodoTxtProvider) Export(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	result := &sync.SyncResult{}

	// Get current application state
	appState := repo.GetAppState()

	// Read existing todo.txt file if it exists
	var todoList todotxtlib.TaskList
	var err error

	todoList, err = ReadTodoFile(p.filePath)
	if err != nil {
		// If file doesn't exist, create an empty one
		if os.IsNotExist(err) {
			todoList = todotxtlib.TaskList{}
		} else {
			return nil, err
		}
	}

	// Map of UUID to todo.txt tasks for easier lookup
	todoTasksByUUID := make(map[string]*todotxtlib.Task)

	// Build lookup map
	for i := range todoList {
		task := &todoList[i]
		if uuid, exists := task.AdditionalTags["uuid"]; exists && uuid != "" {
			todoTasksByUUID[uuid] = task
		}
	}

	// Export all relevant repository tasks to todo.txt
	var tasksToWrite []todotxtlib.Task

	// Keep track of which tasks were processed
	processedTasks := make(map[string]bool)

	for id, repoTask := range appState.Tasks {
		// Skip tasks that aren't from this todo.txt file and don't match filter
		if filter != nil && !filter.FilterTask(repoTask) {
			continue
		}

		// Skip tasks marked as deleted
		if _, isDeleted := repoTask.AdditionalTags["deleted"]; isDeleted {
			continue
		}

		// Apply modifier if provided
		if modifier != nil {
			repoTask = modifier.ModifyTask(repoTask)
		}

		// Convert to todo.txt task
		todoTask := repoTaskToTodoTask(repoTask)

		// Check if task exists in todo.txt
		uuidStr := id.String()
		if existingTask, exists := todoTasksByUUID[uuidStr]; exists {
			// Update existing task
			*existingTask = todoTask
			result.ToSource++
			// Mark as processed
			processedTasks[uuidStr] = true
		} else {
			// Add new task
			tasksToWrite = append(tasksToWrite, todoTask)
			result.ToSource++
		}
	}

	// Add existing tasks that weren't updated (avoiding duplicates)
	for i := range todoList {
		task := todoList[i]
		// Only add tasks that weren't already processed
		if uuid, hasUUID := task.AdditionalTags["uuid"]; hasUUID {
			if !processedTasks[uuid] {
				tasksToWrite = append(tasksToWrite, task)
			}
		} else {
			tasksToWrite = append(tasksToWrite, task)
		}
	}

	// Write changes to todo.txt
	err = WriteTodoFile(tasksToWrite, p.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to write tasks to todo.txt: %w", err)
	}

	return result, nil
}

// Sync synchronizes tasks between the todo.txt file and the t5 repository
func (p *TodoTxtProvider) Sync(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// First import from todo.txt to repository
	importResult, err := p.Import(repo, filter, modifier)
	if err != nil {
		return nil, fmt.Errorf("import failed: %w", err)
	}

	// Then export from repository to todo.txt
	exportResult, err := p.Export(repo, filter, modifier)
	if err != nil {
		return nil, fmt.Errorf("export failed: %w", err)
	}

	// Combine results
	result := &sync.SyncResult{
		Added:      importResult.Added,
		Updated:    importResult.Updated + exportResult.Updated,
		Skipped:    importResult.Skipped + exportResult.Skipped,
		Conflicts:  importResult.Conflicts + exportResult.Conflicts,
		FromSource: importResult.FromSource,
		ToSource:   exportResult.ToSource,
	}

	return result, nil
}

// ProvidesImport returns true if this provider supports importing tasks
func (p *TodoTxtProvider) ProvidesImport() bool {
	return true
}

// ProvidesExport returns true if this provider supports exporting tasks
func (p *TodoTxtProvider) ProvidesExport() bool {
	return true
}

// todoTaskToRepoTask converts a todo.txt task to a repository task
func todoTaskToRepoTask(task todotxtlib.Task) core.Task {
	repoTask := core.Task{
		Todo:          task.Todo,
		Priority:      task.Priority,
		Projects:      task.Projects,
		Contexts:      task.Contexts,
		Completed:     task.Completed,
		CreatedDate:   task.CreatedDate,
		DueDate:       task.DueDate,
		CompletedDate: task.CompletedDate,
	}

	// Copy additional tags
	repoTask.AdditionalTags = make(map[string]string)
	for k, v := range task.AdditionalTags {
		repoTask.AdditionalTags[k] = v
	}

	// Set UUID if it exists
	if uuid, exists := task.AdditionalTags["uuid"]; exists && uuid != "" {
		id, err := utils.DecodeUUID(uuid)
		if err == nil {
			repoTask.ID = id
		}
	}

	return repoTask
}

// todoTaskToString converts a repository task to a todo.txt string
func todoTaskToString(task core.Task) string {
	// Convert to todo.txt task
	todoTask := repoTaskToTodoTask(task)

	// Convert to string
	return todoTask.String()
}

// NewTodoTxtProviderFactory creates a factory function for todo.txt providers
func NewTodoTxtProviderFactory() sync.SyncProviderFactory {
	return func(config sync.ProviderConfig) (sync.SyncProvider, error) {
		// Get file path from config
		filePath, ok := config.Params["file_path"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid file_path parameter for todo.txt provider")
		}

		return NewTodoTxtProvider(filePath), nil
	}
}
