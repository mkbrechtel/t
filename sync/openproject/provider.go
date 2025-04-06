package openproject

import (
	"fmt"
	"time"
	
	"t5.mkbrechtel.dev/t5/core"
	"t5.mkbrechtel.dev/t5/sync"
	"t5.mkbrechtel.dev/t5/utils"
)

// OpenProjectProvider handles synchronization with OpenProject work packages
type OpenProjectProvider struct {
	baseURL    string
	apiKey     string
	queryIDs   []string
	taskPrefix string
}

// NewOpenProjectProvider creates a new OpenProject sync provider
func NewOpenProjectProvider(baseURL, apiKey string, queryIDs []string, taskPrefix string) *OpenProjectProvider {
	return &OpenProjectProvider{
		baseURL:    baseURL,
		apiKey:     apiKey,
		queryIDs:   queryIDs,
		taskPrefix: taskPrefix,
	}
}

// Name returns the name of the sync provider
func (p *OpenProjectProvider) Name() string {
	return "OpenProject"
}

// Description returns a description of the sync provider
func (p *OpenProjectProvider) Description() string {
	return "OpenProject work packages sync provider"
}

// SupportedDirections returns the sync directions supported by this provider
func (p *OpenProjectProvider) SupportedDirections() []sync.SyncDirection {
	return []sync.SyncDirection{
		sync.SyncDirectionImport, // Only import from OpenProject is supported for now
	}
}

// Import imports tasks from OpenProject into the t5 repository
func (p *OpenProjectProvider) Import(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	result := &sync.SyncResult{}
	
	// Get current application state
	appState := repo.GetAppState()
	
	// Track tasks by URL for quick lookup
	tasksByURL := make(map[string]*core.Task)
	for id, task := range appState.Tasks {
		if url, exists := task.AdditionalTags["url"]; exists {
			taskCopy := task
			taskCopy.ID = id
			tasksByURL[url] = &taskCopy
		}
	}
	
	// Process each query
	for _, queryID := range p.queryIDs {
		// Fetch work packages for this query
		workPackages, err := GetWorkPackages(p.baseURL, p.apiKey, queryID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch OpenProject work packages for query %s: %w", queryID, err)
		}
		
		// Create task list from work packages
		taskList := CreateTaskList(workPackages, p.taskPrefix, p.baseURL)
		
		// Process each task from the task list
		for _, openProjectTask := range taskList {
			url, hasURL := openProjectTask.AdditionalTags["url"]
			if !hasURL {
				result.Skipped++
				continue
			}
			
			// Convert from todo.txt task to repo task
			repoTask := core.Task{
				Todo:      openProjectTask.Todo,
				Priority:  openProjectTask.Priority,
				Projects:  openProjectTask.Projects,
				Contexts:  openProjectTask.Contexts,
				Completed: openProjectTask.Completed,
			}
			
			// Set dates
			if !openProjectTask.CreatedDate.IsZero() {
				repoTask.CreatedDate = openProjectTask.CreatedDate
			} else {
				repoTask.CreatedDate = time.Now()
			}
			
			if !openProjectTask.DueDate.IsZero() {
				repoTask.DueDate = openProjectTask.DueDate
			}
			
			if !openProjectTask.CompletedDate.IsZero() {
				repoTask.CompletedDate = openProjectTask.CompletedDate
			}
			
			// Copy additional tags
			repoTask.AdditionalTags = make(map[string]string)
			for k, v := range openProjectTask.AdditionalTags {
				repoTask.AdditionalTags[k] = v
			}
			
			// Set source
			repoTask.Source = "openproject"
			
			// Apply filter if provided
			if filter != nil && !filter.FilterTask(repoTask) {
				result.Skipped++
				continue
			}
			
			// Apply modifier if provided
			if modifier != nil {
				repoTask = modifier.ModifyTask(repoTask)
			}
			
			// Check if task already exists
			if existingTask, exists := tasksByURL[url]; exists {
				// Update existing task if needed
				repoTask.ID = existingTask.ID // Preserve ID
				
				// Check if task has changed
				if taskNeedsUpdate(existingTask, &repoTask) {
					// Save task update
					taskString := taskToString(repoTask)
					event := core.NewTodoTxtTaskUpdate(taskString, "openproject")
					if err := repo.SaveEvent(event); err != nil {
						return nil, fmt.Errorf("failed to update task from OpenProject: %w", err)
					}
					result.Updated++
					result.FromSource++
				} else {
					result.Skipped++
				}
			} else {
				// Create new task
				taskString := taskToString(repoTask)
				event := core.NewTodoTxtTaskUpdate(taskString, "openproject")
				if err := repo.SaveEvent(event); err != nil {
					return nil, fmt.Errorf("failed to create task from OpenProject: %w", err)
				}
				result.Added++
				result.FromSource++
			}
		}
	}
	
	return result, nil
}

// Export exports tasks from the t5 repository to OpenProject - NOT IMPLEMENTED
func (p *OpenProjectProvider) Export(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// OpenProject export is not implemented yet
	return &sync.SyncResult{}, fmt.Errorf("exporting to OpenProject is not yet implemented")
}

// Sync synchronizes tasks between OpenProject and the t5 repository
func (p *OpenProjectProvider) Sync(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// Since export isn't implemented, sync is just an import
	return p.Import(repo, filter, modifier)
}

// taskNeedsUpdate checks if the task needs updating
func taskNeedsUpdate(existing, source *core.Task) bool {
	// Check if todo text has changed
	if existing.Todo != source.Todo {
		return true
	}
	
	// Check if due date has changed
	if !existing.DueDate.Equal(source.DueDate) {
		return true
	}
	
	// Check if completion status has changed
	if existing.Completed != source.Completed {
		return true
	}
	
	// Check if threshold date has changed
	existingT, existingHasT := existing.AdditionalTags["t"]
	sourceT, sourceHasT := source.AdditionalTags["t"]
	if existingHasT != sourceHasT || existingT != sourceT {
		return true
	}
	
	return false
}

// taskToString converts a Task to a todo.txt string
func taskToString(task core.Task) string {
	todoTxt := ""
	
	// Add completion mark if completed
	if task.Completed {
		todoTxt += "x "
		// Add completion date if available
		if !task.CompletedDate.IsZero() {
			todoTxt += task.CompletedDate.Format("2006-01-02") + " "
		}
	}
	
	// Add priority if available
	if task.Priority != "" {
		todoTxt += "(" + task.Priority + ") "
	}
	
	// Add creation date if available
	if !task.CreatedDate.IsZero() {
		todoTxt += task.CreatedDate.Format("2006-01-02") + " "
	}
	
	// Add main text
	todoTxt += task.Todo
	
	// Add projects
	for _, project := range task.Projects {
		todoTxt += " +" + project
	}
	
	// Add contexts
	for _, context := range task.Contexts {
		todoTxt += " @" + context
	}
	
	// Add due date if available
	if !task.DueDate.IsZero() {
		todoTxt += " due:" + task.DueDate.Format("2006-01-02")
	}
	
	// Add other tags
	for key, value := range task.AdditionalTags {
		// Skip internal tags
		if key == "uuid" || key == "id" {
			continue
		}
		todoTxt += " " + key + ":" + value
	}
	
	// Add UUID
	todoTxt += " uuid:" + utils.LongEncodeUUID(task.ID)
	
	return todoTxt
}

// NewOpenProjectProviderFactory creates a factory function for OpenProject providers
func NewOpenProjectProviderFactory() sync.SyncProviderFactory {
	return func(config sync.ProviderConfig) (sync.SyncProvider, error) {
		// Get parameters from config
		baseURL, ok := config.Params["base_url"].(string)
		if !ok || baseURL == "" {
			return nil, fmt.Errorf("missing or invalid base_url parameter for OpenProject provider")
		}
		
		apiKey, ok := config.Params["api_key"].(string)
		if !ok || apiKey == "" {
			return nil, fmt.Errorf("missing or invalid api_key parameter for OpenProject provider")
		}
		
		// Get query IDs as an array
		var queryIDs []string
		queryIDsInterface, ok := config.Params["query_ids"]
		if !ok {
			return nil, fmt.Errorf("missing query_ids parameter for OpenProject provider")
		}
		
		// Handle different ways query_ids could be specified
		switch v := queryIDsInterface.(type) {
		case []interface{}:
			queryIDs = make([]string, len(v))
			for i, id := range v {
				queryIDs[i] = fmt.Sprintf("%v", id)
			}
		case string:
			queryIDs = []string{v}
		default:
			return nil, fmt.Errorf("invalid query_ids parameter for OpenProject provider, must be a string or array of strings")
		}
		
		taskPrefix, ok := config.Params["task_prefix"].(string)
		if !ok {
			taskPrefix = "OpenProject: "
		}
		
		return NewOpenProjectProvider(baseURL, apiKey, queryIDs, taskPrefix), nil
	}
}

// Register the OpenProject provider factory
func init() {
	sync.Register("openproject", NewOpenProjectProviderFactory())
}