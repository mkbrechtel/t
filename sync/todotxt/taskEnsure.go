package todo

import (
	"time"
	todo "github.com/1set/todotxt"
	"t5.mkbrechtel.dev/t5/utils"
	uuidv7 "github.com/gofrs/uuid/v5"
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
	if config.EnforceCreationDate {
		ensureCreationDate(task)
	}
	if config.EnforceCompletionDate && task.IsCompleted() && !task.HasCompletedDate() {
		ensureCompletionDate(task)
	}
	ensureIdentifier(task, config.PreferShortIDs)
	ensureTags(task, config.DefaultTags)
}

// EnsureTaskListProperties applies property assurance to all tasks in a list
func EnsureTaskListProperties(taskList todo.TaskList, config TaskEnsureConfig) todo.TaskList {
	for i := range taskList {
		EnsureTaskProperties(&taskList[i], config)
	}
	return taskList
}

// ensureCreationDate ensures tasks have a creation date
func ensureCreationDate(task *todo.Task) {
	if !task.HasCreatedDate() {
		task.CreatedDate = time.Now()
	}
}

// ensureCompletionDate ensures completed tasks have a completion date
func ensureCompletionDate(task *todo.Task) {
	task.CompletedDate = time.Now()
}

// ensureIdentifier ensures tasks have proper identifiers and returns the UUID
// If the task has a non-UUID id, it preserves that id and adds a uuid.
// If the task has a UUID id, it formats it according to preference.
func ensureIdentifier(task *todo.Task, preferShortIDs bool) (uuidv7.UUID) {
	if task.AdditionalTags == nil {
		task.AdditionalTags = make(map[string]string)
	}

	var id = utils.NewUUID() // Default to new UUID
	var hasUUID = false
	var hasID = false

	// Check for existing UUID tag
	if uuidStr, exists := task.AdditionalTags["uuid"]; exists && uuidStr != "" {
		if parsedId, err := utils.DecodeUUID(uuidStr); err == nil {
			id = parsedId
			hasUUID = true
		}
	}

	// Check for existing ID tag
	if idStr, exists := task.AdditionalTags["id"]; exists && idStr != "" {
		if utils.IsUUID(idStr) {
			// ID is a valid UUID
			parsedId, _ := utils.DecodeUUID(idStr)
			id = parsedId
			hasID = true
			// This is a UUID formatted as ID, so handle according to preference
			if preferShortIDs {
				task.AdditionalTags["id"] = utils.ShortEncodeUUID(id)
				delete(task.AdditionalTags, "uuid")
			} else {
				task.AdditionalTags["uuid"] = utils.LongEncodeUUID(id)
				delete(task.AdditionalTags, "id")
			}
		} else {
			// ID exists but is not a UUID - keep it and ensure a UUID is present
			hasID = true
			task.AdditionalTags["uuid"] = utils.LongEncodeUUID(id)
		}
	}

	// If no UUID or ID exists, create based on preference
	if !hasUUID && !hasID {
		if preferShortIDs {
			task.AdditionalTags["id"] = utils.ShortEncodeUUID(id)
		} else {
			task.AdditionalTags["uuid"] = utils.LongEncodeUUID(id)
		}
	} else if hasUUID && !hasID && preferShortIDs {
		// Has UUID but no ID, and we prefer short IDs
		task.AdditionalTags["id"] = utils.ShortEncodeUUID(id)
		delete(task.AdditionalTags, "uuid")
	} else if hasID && !hasUUID && !preferShortIDs {
		// Has ID but no UUID, and we don't prefer short IDs
		// Only convert to UUID if the ID is actually a UUID
		if utils.IsUUID(task.AdditionalTags["id"]) {
			task.AdditionalTags["uuid"] = utils.LongEncodeUUID(id)
			delete(task.AdditionalTags, "id")
		} else {
			// ID is not a UUID, so add a UUID tag
			task.AdditionalTags["uuid"] = utils.LongEncodeUUID(id)
		}
	}

	return id
}

// ensureDefaultTags ensures all default tags are present
func ensureTags(task *todo.Task, tags map[string]string) {
	if task.AdditionalTags == nil {
		task.AdditionalTags = make(map[string]string)
	}
	
	for key, value := range tags {
		if _, exists := task.AdditionalTags[key]; !exists {
			task.AdditionalTags[key] = value
		}
	}
}
