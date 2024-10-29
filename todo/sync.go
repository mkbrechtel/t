package todo

import (
    "time"
    todo "github.com/1set/todotxt"
)

// SyncResult contains statistics about the sync operation
type SyncResult struct {
    Added   int
    Updated int
    Skipped int
}

// SyncTaskLists merges tasks from source into target list, using URL as unique identifier
// Note: This function now modifies the target list directly and returns it along with the result
func SyncTaskLists(target, source todo.TaskList) (todo.TaskList, *SyncResult, error) {
    result := &SyncResult{}
    
    // Create a map of existing tasks by URL for efficient lookup
    existingTasks := make(map[string]*todo.Task)
    for i := range target {
        if url, exists := target[i].AdditionalTags["url"]; exists {
            existingTasks[url] = &target[i]
        }
    }
    
    // Process each task from the source
    for _, sourceTask := range source {
        sourceURL, hasURL := sourceTask.AdditionalTags["url"]
        if !hasURL {
            result.Skipped++
            continue
        }
        
        if existingTask, exists := existingTasks[sourceURL]; exists {
            // Update existing task if needed
            if shouldUpdateTask(existingTask, &sourceTask) {
                updateTaskContent(existingTask, &sourceTask)
                result.Updated++
            } else {
                result.Skipped++
            }
        } else {
            // Add new task
            target = append(target, sourceTask)
            result.Added++
        }
    }
    
    return target, result, nil
}

// shouldUpdateTask determines if the existing task needs updating
func shouldUpdateTask(existing, source *todo.Task) bool {
    // Check if todo text has changed
    if existing.Todo != source.Todo {
        return true
    }
    
    // Check if due date has changed
    if !existing.DueDate.Equal(source.DueDate) {
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

// updateTaskContent updates the content of an existing task from a source task
func updateTaskContent(existing, source *todo.Task) {
    // Preserve completion status and completed date
    wasCompleted := existing.Completed
    completedDate := existing.CompletedDate
    
    // Update main task fields
    existing.Todo = source.Todo
    existing.DueDate = source.DueDate
    existing.Projects = source.Projects
    existing.Contexts = source.Contexts
    
    // Update additional tags while preserving uuid
    uuid := existing.AdditionalTags["uuid"]
    existing.AdditionalTags = make(map[string]string)
    for k, v := range source.AdditionalTags {
        existing.AdditionalTags[k] = v
    }
    existing.AdditionalTags["uuid"] = uuid
    
    // Restore completion status
    existing.Completed = wasCompleted
    existing.CompletedDate = completedDate
    
    // Update modified timestamp
    existing.AdditionalTags["modified"] = time.Now().Format(time.RFC3339)
}
