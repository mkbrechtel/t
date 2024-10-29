package todo

import (
	"time"
	todo "github.com/1set/todotxt"
	"t/utils"
)

// TaskEnsureConfig holds configuration options for ensuring task properties
type TaskEnsureConfig struct {
	// PreferShortIDs determines whether to use short-form IDs (true) or long-form UUIDs (false)
	PreferShortIDs bool
	// EnforceCompletionDate determines whether completed tasks must have a completion date
	EnforceCompletionDate bool
	// EnforceCreationDate determines whether tasks must have a creation date
	EnforceCreationDate bool
	// DefaultTags are additional tags that should be present on all tasks
	DefaultTags map[string]string
}

// DefaultEnsureConfig provides sensible defaults for task properties
var DefaultEnsureConfig = TaskEnsureConfig{
	PreferShortIDs:        true,
	EnforceCompletionDate: true,
	EnforceCreationDate:   true,
	DefaultTags:           make(map[string]string),
}

// EnsureTaskProperties ensures a single task has all required properties according to the config
func EnsureTaskProperties(task *todo.Task, config TaskEnsureConfig) {
	ensureCreationDate(task, config)
	ensureCompletionDate(task, config)
	ensureIdentifier(task, config)
	ensureDefaultTags(task, config)
}

// EnsureTaskListProperties applies property assurance to all tasks in a list
func EnsureTaskListProperties(taskList todo.TaskList, config TaskEnsureConfig) todo.TaskList {
	for i := range taskList {
		EnsureTaskProperties(&taskList[i], config)
	}
	return taskList
}

// ensureCreationDate ensures tasks have a creation date
func ensureCreationDate(task *todo.Task, config TaskEnsureConfig) {
	if config.EnforceCreationDate && !task.HasCreatedDate() {
		task.CreatedDate = time.Now()
	}
}

// ensureCompletionDate ensures completed tasks have a completion date
func ensureCompletionDate(task *todo.Task, config TaskEnsureConfig) {
	if config.EnforceCompletionDate && task.IsCompleted() && !task.HasCompletedDate() {
		task.CompletedDate = time.Now()
	}
}

// ensureIdentifier ensures tasks have either a short-form ID or UUID
func ensureIdentifier(task *todo.Task, config TaskEnsureConfig) {
	if task.AdditionalTags == nil {
		task.AdditionalTags = make(map[string]string)
	}

	// Try to get existing UUID first
	var id = utils.NewUUID() // Default to new UUID
	if uuidStr, hasUuid := task.AdditionalTags["uuid"]; hasUuid {
		if parsedId, err := utils.DecodeUUID(uuidStr); err == nil {
			id = parsedId
			delete(task.AdditionalTags, "uuid")
		}
	}

	// If no valid UUID was found in uuid tag, try id tag
	if idStr, hasId := task.AdditionalTags["id"]; hasId {
		if parsedId, err := utils.DecodeUUID(idStr); err == nil {
			id = parsedId
		}
	}

	// Set the ID in preferred format
	if config.PreferShortIDs {
		task.AdditionalTags["id"] = utils.ShortEncodeUUID(id)
		delete(task.AdditionalTags, "uuid")
	} else {
		task.AdditionalTags["uuid"] = utils.LongEncodeUUID(id)
		delete(task.AdditionalTags, "id")
	}
}

// ensureDefaultTags ensures all default tags are present
func ensureDefaultTags(task *todo.Task, config TaskEnsureConfig) {
	if task.AdditionalTags == nil {
		task.AdditionalTags = make(map[string]string)
	}
	
	for key, value := range config.DefaultTags {
		if _, exists := task.AdditionalTags[key]; !exists {
			task.AdditionalTags[key] = value
		}
	}
}
