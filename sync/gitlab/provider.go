package gitlab

import (
	"fmt"
	"time"

	"t5.mkbrechtel.dev/t5/core"
	"t5.mkbrechtel.dev/t5/sync"
	"t5.mkbrechtel.dev/t5/utils"
)

// GitLabProvider handles synchronization with GitLab issues and MRs
type GitLabProvider struct {
	token            string
	baseURL          string
	issuesURL        string
	mergeRequestsURL string
	issuePrefix      string
	mrPrefix         string
}

// NewGitLabProvider creates a new GitLab sync provider
func NewGitLabProvider(token, baseURL, issuesURL, mergeRequestsURL, issuePrefix, mrPrefix string) *GitLabProvider {
	return &GitLabProvider{
		token:            token,
		baseURL:          baseURL,
		issuesURL:        issuesURL,
		mergeRequestsURL: mergeRequestsURL,
		issuePrefix:      issuePrefix,
		mrPrefix:         mrPrefix,
	}
}

// Name returns the name of the sync provider
func (p *GitLabProvider) Name() string {
	return "GitLab"
}

// Description returns a description of the sync provider
func (p *GitLabProvider) Description() string {
	return "GitLab issues and merge requests sync provider"
}

// SupportedDirections returns the sync directions supported by this provider
func (p *GitLabProvider) SupportedDirections() []sync.SyncDirection {
	return []sync.SyncDirection{
		sync.SyncDirectionImport, // Only import from GitLab is supported for now
	}
}

// Import imports tasks from GitLab into the t5 repository
func (p *GitLabProvider) Import(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	result := &sync.SyncResult{}

	// Fetch issues
	issues, err := GetUserIssues(p.token, p.baseURL, p.issuesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitLab issues: %w", err)
	}

	// Create task list from issues
	issueTaskList := CreateIssueTaskList(issues, p.issuePrefix)

	// Fetch merge requests
	mergeRequests, err := GetUserMergeRequests(p.token, p.baseURL, p.mergeRequestsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitLab merge requests: %w", err)
	}

	// Create task list from merge requests
	mrTaskList := CreateMergeRequestTaskList(mergeRequests, p.mrPrefix)

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

	// Process the combined tasks (issues and merge requests)
	for _, gitlabTask := range append(issueTaskList, mrTaskList...) {
		url, hasURL := gitlabTask.AdditionalTags["url"]
		if !hasURL {
			result.Skipped++
			continue
		}

		// Convert from todo.txt task to repo task
		repoTask := core.Task{
			Todo:      gitlabTask.Todo,
			Priority:  gitlabTask.Priority,
			Projects:  gitlabTask.Projects,
			Contexts:  gitlabTask.Contexts,
			Completed: gitlabTask.Completed,
		}

		// Set dates
		if !gitlabTask.CreatedDate.IsZero() {
			repoTask.CreatedDate = gitlabTask.CreatedDate
		} else {
			repoTask.CreatedDate = time.Now()
		}

		if !gitlabTask.DueDate.IsZero() {
			repoTask.DueDate = gitlabTask.DueDate
		}

		if !gitlabTask.CompletedDate.IsZero() {
			repoTask.CompletedDate = gitlabTask.CompletedDate
		}

		// Copy additional tags
		repoTask.AdditionalTags = make(map[string]string)
		for k, v := range gitlabTask.AdditionalTags {
			repoTask.AdditionalTags[k] = v
		}

		// Set source
		repoTask.Source = "gitlab"

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
				event := core.NewTodoTxtTaskUpdate(taskString, "gitlab")
				if err := repo.SaveEvent(event); err != nil {
					return nil, fmt.Errorf("failed to update task from GitLab: %w", err)
				}
				result.Updated++
				result.FromSource++
			} else {
				result.Skipped++
			}
		} else {
			// Create new task
			taskString := taskToString(repoTask)
			event := core.NewTodoTxtTaskUpdate(taskString, "gitlab")
			if err := repo.SaveEvent(event); err != nil {
				return nil, fmt.Errorf("failed to create task from GitLab: %w", err)
			}
			result.Added++
			result.FromSource++
		}
	}

	return result, nil
}

// Export exports tasks from the t5 repository to GitLab - NOT IMPLEMENTED
func (p *GitLabProvider) Export(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// GitLab export is not implemented yet
	return &sync.SyncResult{}, fmt.Errorf("exporting to GitLab is not yet implemented")
}

// Sync synchronizes tasks between GitLab and the t5 repository
func (p *GitLabProvider) Sync(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// Since export isn't implemented, sync is just an import
	return p.Import(repo, filter, modifier)
}

// ProvidesImport returns true if this provider supports importing tasks
func (p *GitLabProvider) ProvidesImport() bool {
	return true
}

// ProvidesExport returns true if this provider supports exporting tasks
func (p *GitLabProvider) ProvidesExport() bool {
	return false
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

// NewGitLabProviderFactory creates a factory function for GitLab providers
func NewGitLabProviderFactory() sync.SyncProviderFactory {
	return func(config sync.ProviderConfig) (sync.SyncProvider, error) {
		// Get parameters from config
		token, ok := config.Params["token"].(string)
		if !ok || token == "" {
			return nil, fmt.Errorf("missing or invalid token parameter for GitLab provider")
		}

		baseURL, ok := config.Params["base_url"].(string)
		if !ok || baseURL == "" {
			baseURL = "https://gitlab.com/api/v4"
		}

		issuesURL, ok := config.Params["issues_url"].(string)
		if !ok || issuesURL == "" {
			issuesURL = "/issues?state=opened&scope=assigned_to_me"
		}

		mergeRequestsURL, ok := config.Params["merge_requests_url"].(string)
		if !ok || mergeRequestsURL == "" {
			mergeRequestsURL = "/merge_requests?state=opened&scope=assigned_to_me"
		}

		issuePrefix, ok := config.Params["issue_prefix"].(string)
		if !ok {
			issuePrefix = "GitLab: "
		}

		mrPrefix, ok := config.Params["mr_prefix"].(string)
		if !ok {
			mrPrefix = "GitLab MR: "
		}

		return NewGitLabProvider(token, baseURL, issuesURL, mergeRequestsURL, issuePrefix, mrPrefix), nil
	}
}
