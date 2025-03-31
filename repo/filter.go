package repo

import (
	"time"

	"t5.mkbrechtel.dev/model"
)

type TaskFilter interface {
	FilterTask(model.Task) bool
}

// ProjectFilter filters tasks that have specific projects
type ProjectFilter struct {
	Projects []string
}

func (f ProjectFilter) FilterTask(task model.Task) bool {
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

func (f ContextFilter) FilterTask(task model.Task) bool {
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

func (f CompletionFilter) FilterTask(task model.Task) bool {
	if f.Completed == nil {
		return true
	}
	return task.Completed == *f.Completed
}

// PriorityFilter filters tasks by priority
type PriorityFilter struct {
	Priority string
}

func (f PriorityFilter) FilterTask(task model.Task) bool {
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

// func (f DueDateFilter) FilterTask(task model.Task) bool {
// 	if f.Before == nil && f.After == nil {
// 		return true
// 	}
// 	if task.DueDate == nil {
// 		return false
// 	}
// 	if f.Before != nil && task.DueDate.After(*f.Before) {
// 		return false
// 	}
// 	if f.After != nil && task.DueDate.Before(*f.After) {
// 		return false
// 	}
// 	return true
// }

// CreatedDateFilter filters tasks by creation date range
type CreatedDateFilter struct {
	Before *time.Time
	After  *time.Time
}

// func (f CreatedDateFilter) FilterTask(task model.Task) bool {
// 	if f.Before == nil && f.After == nil {
// 		return true
// 	}
// 	if task.CreatedDate == nil {
// 		return false
// 	}
// 	if f.Before != nil && task.CreatedDate.After(*f.Before) {
// 		return false
// 	}
// 	if f.After != nil && task.CreatedDate.Before(*f.After) {
// 		return false
// 	}
// 	return true
// }
