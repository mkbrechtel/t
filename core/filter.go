package core

import (
	"log"
	"regexp"
	"strings"
	"time"
)

type TaskFilter interface {
	FilterTask(Task) bool
}

// ProjectFilter filters tasks that have specific projects
type ProjectFilter struct {
	Projects []string
}

func (f ProjectFilter) FilterTask(task Task) bool {
	if len(f.Projects) == 0 {
		return true
	}
	for _, project := range f.Projects {
		for _, taskProject := range task.Projects {
			if project == taskProject {
				return true
			}
		}
	}
	return false
}

// ContextFilter filters tasks that have specific contexts
type ContextFilter struct {
	Contexts []string
}

func (f ContextFilter) FilterTask(task Task) bool {
	if len(f.Contexts) == 0 {
		return true
	}
	for _, context := range f.Contexts {
		for _, taskContext := range task.Contexts {
			if context == taskContext {
				return true
			}
		}
	}
	return false
}

// CompletionFilter filters tasks based on completion status
type CompletionFilter struct {
	Completed *bool
}

func (f CompletionFilter) FilterTask(task Task) bool {
	if f.Completed == nil {
		return true
	}
	return task.Completed == *f.Completed
}

// PriorityFilter filters tasks by priority
type PriorityFilter struct {
	Priority string
}

func (f PriorityFilter) FilterTask(task Task) bool {
	if f.Priority == "" {
		return true
	}
	return task.Priority == f.Priority
}

// DueDateFilter filters tasks by due date range
type DueDateFilter struct {
	Before *time.Time
	After  *time.Time
}

func (f DueDateFilter) FilterTask(task Task) bool {
	if f.Before == nil && f.After == nil {
		return true
	}
	if task.DueDate.IsZero() {
		return false
	}
	if f.Before != nil && task.DueDate.After(*f.Before) {
		return false
	}
	if f.After != nil && task.DueDate.Before(*f.After) {
		return false
	}
	return true
}

// CreatedDateFilter filters tasks by creation date range
type CreatedDateFilter struct {
	Before *time.Time
	After  *time.Time
}

func (f CreatedDateFilter) FilterTask(task Task) bool {
	if f.Before == nil && f.After == nil {
		return true
	}
	if task.CreatedDate.IsZero() {
		return false
	}
	if f.Before != nil && task.CreatedDate.After(*f.Before) {
		return false
	}
	if f.After != nil && task.CreatedDate.Before(*f.After) {
		return false
	}
	return true
}

// TagFilter filters tasks by the presence of a specific tag
type TagFilter struct {
	Tag string
}

func (f TagFilter) FilterTask(task Task) bool {
	if f.Tag == "" {
		return true
	}
	
	// First check if it's a key-only tag
	if _, exists := task.AdditionalTags[f.Tag]; exists {
		return true
	}
	
	// Then check for any tag key that starts with the given tag
	for key := range task.AdditionalTags {
		if strings.HasPrefix(key, f.Tag) {
			return true
		}
	}
	
	return false
}

// RegexFilter filters tasks matching a regular expression pattern
type RegexFilter struct {
	Pattern string
}

func (f RegexFilter) FilterTask(task Task) bool {
	if f.Pattern == "" {
		return true
	}
	
	// Compile the regex pattern
	pattern, err := regexp.Compile(f.Pattern)
	if err != nil {
		// If pattern is invalid, log an error and return false
		log.Printf("Invalid regex pattern: %v", err)
		return false
	}
	
	// Check task text
	if pattern.MatchString(task.Todo) {
		return true
	}
	
	// Check projects
	for _, project := range task.Projects {
		if pattern.MatchString(project) {
			return true
		}
	}
	
	// Check contexts
	for _, context := range task.Contexts {
		if pattern.MatchString(context) {
			return true
		}
	}
	
	// Check tag keys and values
	for key, value := range task.AdditionalTags {
		if pattern.MatchString(key) || pattern.MatchString(value) {
			return true
		}
	}
	
	return false
}

// AndFilter combines multiple filters with AND logic
type AndFilter struct {
	Filters []TaskFilter
}

func (f AndFilter) FilterTask(task Task) bool {
	if len(f.Filters) == 0 {
		return true
	}
	
	for _, filter := range f.Filters {
		if !filter.FilterTask(task) {
			return false
		}
	}
	
	return true
}

// OrFilter combines multiple filters with OR logic
type OrFilter struct {
	Filters []TaskFilter
}

func (f OrFilter) FilterTask(task Task) bool {
	if len(f.Filters) == 0 {
		return true
	}
	
	for _, filter := range f.Filters {
		if filter.FilterTask(task) {
			return true
		}
	}
	
	return false
}

// NotFilter negates the result of a filter
type NotFilter struct {
	Filter TaskFilter
}

func (f NotFilter) FilterTask(task Task) bool {
	return !f.Filter.FilterTask(task)
}
