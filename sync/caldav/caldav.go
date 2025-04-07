package caldav

import (
	"fmt"
	"time"

	"t5.mkbrechtel.dev/t5/core"
	"t5.mkbrechtel.dev/t5/sync"
	"t5.mkbrechtel.dev/t5/utils"
)

// CalDAVTask represents a task in a CalDAV server
type CalDAVTask struct {
	UID             string
	URL             string
	Summary         string
	Description     string
	Categories      []string
	Priority        int    // 0 (undefined), 1 (high), 5 (medium), 9 (low)
	Status          string // NEEDS-ACTION, IN-PROCESS, COMPLETED, CANCELLED
	Created         time.Time
	Modified        time.Time
	Due             time.Time
	Start           time.Time
	Completed       time.Time
	PercentComplete int // 0-100
}

// CalDAVProvider is a mock implementation of a CalDAV sync provider
// In a real implementation, this would connect to a CalDAV server
type CalDAVProvider struct {
	serverURL  string
	username   string
	password   string
	calendarID string
	taskPrefix string
}

// NewCalDAVProvider creates a new CalDAV sync provider
func NewCalDAVProvider(serverURL, username, password, calendarID, taskPrefix string) *CalDAVProvider {
	return &CalDAVProvider{
		serverURL:  serverURL,
		username:   username,
		password:   password,
		calendarID: calendarID,
		taskPrefix: taskPrefix,
	}
}

// Name returns the name of the sync provider
func (p *CalDAVProvider) Name() string {
	return "CalDAV"
}

// Description returns a description of the sync provider
func (p *CalDAVProvider) Description() string {
	return fmt.Sprintf("CalDAV sync provider for %s", p.serverURL)
}

// SupportedDirections returns the sync directions supported by this provider
func (p *CalDAVProvider) SupportedDirections() []sync.SyncDirection {
	return []sync.SyncDirection{
		sync.SyncDirectionImport, // Only import is implemented for now
		// Would need a full CalDAV client library to support export and bidirectional sync
	}
}

// Import imports tasks from a CalDAV server into the t5 repository
// This is a mock implementation - in a real app, this would connect to a CalDAV server
func (p *CalDAVProvider) Import(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// This is a placeholder for CalDAV implementation
	return &sync.SyncResult{}, fmt.Errorf("CalDAV import not implemented yet - needs a CalDAV client library")
}

// Export exports tasks from the t5 repository to a CalDAV server - NOT IMPLEMENTED
func (p *CalDAVProvider) Export(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// CalDAV export is not implemented yet
	return &sync.SyncResult{}, fmt.Errorf("exporting to CalDAV is not yet implemented")
}

// Sync synchronizes tasks between a CalDAV server and the t5 repository
func (p *CalDAVProvider) Sync(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*sync.SyncResult, error) {
	// Since export isn't implemented, sync is just an import
	return p.Import(repo, filter, modifier)
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

// NewCalDAVProviderFactory creates a factory function for CalDAV providers
func NewCalDAVProviderFactory() sync.SyncProviderFactory {
	return func(config sync.ProviderConfig) (sync.SyncProvider, error) {
		// Get parameters from config
		serverURL, ok := config.Params["server_url"].(string)
		if !ok || serverURL == "" {
			return nil, fmt.Errorf("missing or invalid server_url parameter for CalDAV provider")
		}

		username, ok := config.Params["username"].(string)
		if !ok || username == "" {
			return nil, fmt.Errorf("missing or invalid username parameter for CalDAV provider")
		}

		password, ok := config.Params["password"].(string)
		if !ok || password == "" {
			return nil, fmt.Errorf("missing or invalid password parameter for CalDAV provider")
		}

		calendarID, ok := config.Params["calendar_id"].(string)
		if !ok || calendarID == "" {
			return nil, fmt.Errorf("missing or invalid calendar_id parameter for CalDAV provider")
		}

		taskPrefix, ok := config.Params["task_prefix"].(string)
		if !ok {
			taskPrefix = "CalDAV: "
		}

		return NewCalDAVProvider(serverURL, username, password, calendarID, taskPrefix), nil
	}
}

// Register the CalDAV provider factory
func init() {
	sync.Register("caldav", NewCalDAVProviderFactory())
}
