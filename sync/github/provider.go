package github

import (
	"fmt"
	"time"
	
	"t5.mkbrechtel.dev/t5/core"
	"t5.mkbrechtel.dev/t5/sync"
	"t5.mkbrechtel.dev/t5/utils"
)

// GitHubProvider handles synchronization with GitHub issues and PRs
type GitHubProvider struct {
	token       string
	baseURL     string
	issuesURL   string
	prsURL      string
	issuePrefix string
	prPrefix    string
}

// NewGitHubProvider creates a new GitHub sync provider
func NewGitHubProvider(token, baseURL, issuesURL, prsURL, issuePrefix, prPrefix string) *GitHubProvider {
	return &GitHubProvider{
		token:       token,
		baseURL:     baseURL,
		issuesURL:   issuesURL,
		prsURL:      prsURL,
		issuePrefix: issuePrefix,
		prPrefix:    prPrefix,
	}
}

// Name returns the name of the sync provider
func (p *GitHubProvider) Name() string {
	return "GitHub"
}

// Description returns a description of the sync provider
func (p *GitHubProvider) Description() string {
	return "GitHub issues and pull requests sync provider"
}

// SupportedDirections returns the sync directions supported by this provider
func (p *GitHubProvider) SupportedDirections() []sync.SyncDirection {
	return []sync.SyncDirection{
		sync.SyncDirectionImport, // Only import from GitHub is supported for now
	}
}

// Import imports tasks from GitHub into the t5 repository
func (p *GitHubProvider) Import(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	result := &sync.SyncResult{}
	
	// Fetch issues
	issues, err := GetUserIssues(p.token, p.baseURL, p.issuesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub issues: %w", err)
	}
	
	// Create task list from issues
	taskList := CreateTaskList(issues, p.issuePrefix, p.prPrefix)
	
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
	
	// Process each task from GitHub
	for _, githubTask := range taskList {
		url, hasURL := githubTask.AdditionalTags["url"]
		if !hasURL {
			result.Skipped++
			continue
		}
		
		// Convert from todo.txt task to repo task
		repoTask := core.Task{
			Todo:      githubTask.Todo,
			Priority:  githubTask.Priority,
			Projects:  githubTask.Projects,
			Contexts:  githubTask.Contexts,
			Completed: githubTask.Completed,
		}
		
		// Set dates
		if !githubTask.CreatedDate.IsZero() {
			repoTask.CreatedDate = githubTask.CreatedDate
		} else {
			repoTask.CreatedDate = time.Now()
		}
		
		if !githubTask.DueDate.IsZero() {
			repoTask.DueDate = githubTask.DueDate
		}
		
		if !githubTask.CompletedDate.IsZero() {
			repoTask.CompletedDate = githubTask.CompletedDate
		}
		
		// Copy additional tags
		repoTask.AdditionalTags = make(map[string]string)
		for k, v := range githubTask.AdditionalTags {
			repoTask.AdditionalTags[k] = v
		}
		
		// Set source
		repoTask.Source = "github"
		
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
				event := core.NewTodoTxtTaskUpdate(taskString, "github")
				if err := repo.SaveEvent(event); err != nil {
					return nil, fmt.Errorf("failed to update task from GitHub: %w", err)
				}
				result.Updated++
				result.FromSource++
			} else {
				result.Skipped++
			}
		} else {
			// Create new task
			taskString := taskToString(repoTask)
			event := core.NewTodoTxtTaskUpdate(taskString, "github")
			if err := repo.SaveEvent(event); err != nil {
				return nil, fmt.Errorf("failed to create task from GitHub: %w", err)
			}
			result.Added++
			result.FromSource++
		}
	}
	
	return result, nil
}

// Export exports tasks from the t5 repository to GitHub - NOT IMPLEMENTED
func (p *GitHubProvider) Export(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// GitHub export is not implemented yet
	return &sync.SyncResult{}, fmt.Errorf("exporting to GitHub is not yet implemented")
}

// Sync synchronizes tasks between GitHub and the t5 repository
func (p *GitHubProvider) Sync(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
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

// NewGitHubProviderFactory creates a factory function for GitHub providers
func NewGitHubProviderFactory() sync.SyncProviderFactory {
	return func(config sync.ProviderConfig) (sync.SyncProvider, error) {
		// Get parameters from config
		token, ok := config.Params["token"].(string)
		if !ok || token == "" {
			return nil, fmt.Errorf("missing or invalid token parameter for GitHub provider")
		}
		
		baseURL, ok := config.Params["base_url"].(string)
		if !ok || baseURL == "" {
			baseURL = "https://api.github.com"
		}
		
		issuesURL, ok := config.Params["issues_url"].(string)
		if !ok || issuesURL == "" {
			issuesURL = "/issues"
		}
		
		prsURL, ok := config.Params["prs_url"].(string)
		if !ok || prsURL == "" {
			prsURL = "/pulls"
		}
		
		issuePrefix, ok := config.Params["issue_prefix"].(string)
		if !ok {
			issuePrefix = "GitHub: "
		}
		
		prPrefix, ok := config.Params["pr_prefix"].(string)
		if !ok {
			prPrefix = "GitHub PR: "
		}
		
		return NewGitHubProvider(token, baseURL, issuesURL, prsURL, issuePrefix, prPrefix), nil
	}
}

// Register the GitHub provider factory
func init() {
	sync.Register("github", NewGitHubProviderFactory())
}