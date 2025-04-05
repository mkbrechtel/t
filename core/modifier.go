package core

import (
	"strings"
	"time"
)

// TaskModifier interface for modifying tasks
type TaskModifier interface {
	ModifyTask(Task) Task
}

// CompletionModifier modifies the completion status of a task
type CompletionModifier struct {
	Completed bool
}

func (m CompletionModifier) ModifyTask(task Task) Task {
	task.Completed = m.Completed
	if m.Completed {
		task.CompletedDate = time.Now()
	} else {
		task.CompletedDate = time.Time{}
	}
	return task
}

// PriorityModifier modifies the priority of a task
type PriorityModifier struct {
	Priority     string
	RemovePriority bool
}

func (m PriorityModifier) ModifyTask(task Task) Task {
	if m.RemovePriority {
		task.Priority = ""
	} else if m.Priority != "" {
		task.Priority = strings.ToUpper(m.Priority)
	}
	return task
}

// ProjectModifier adds or removes projects from a task
type ProjectModifier struct {
	AddProjects    []string
	RemoveProjects []string
}

func (m ProjectModifier) ModifyTask(task Task) Task {
	// Add projects
	for _, project := range m.AddProjects {
		project = strings.TrimSpace(project)
		if project == "" {
			continue
		}
		
		// Remove leading + if present
		if strings.HasPrefix(project, "+") {
			project = project[1:]
		}
		
		// Check if project already exists
		exists := false
		for _, p := range task.Projects {
			if p == project {
				exists = true
				break
			}
		}
		
		if !exists {
			task.Projects = append(task.Projects, project)
		}
	}
	
	// Remove projects
	for _, project := range m.RemoveProjects {
		project = strings.TrimSpace(project)
		if project == "" {
			continue
		}
		
		// Remove leading + if present
		if strings.HasPrefix(project, "+") {
			project = project[1:]
		}
		
		// Remove project from list
		for i, p := range task.Projects {
			if p == project {
				task.Projects = append(task.Projects[:i], task.Projects[i+1:]...)
				break
			}
		}
	}
	
	return task
}

// ContextModifier adds or removes contexts from a task
type ContextModifier struct {
	AddContexts    []string
	RemoveContexts []string
}

func (m ContextModifier) ModifyTask(task Task) Task {
	// Add contexts
	for _, context := range m.AddContexts {
		context = strings.TrimSpace(context)
		if context == "" {
			continue
		}
		
		// Remove leading @ if present
		if strings.HasPrefix(context, "@") {
			context = context[1:]
		}
		
		// Check if context already exists
		exists := false
		for _, c := range task.Contexts {
			if c == context {
				exists = true
				break
			}
		}
		
		if !exists {
			task.Contexts = append(task.Contexts, context)
		}
	}
	
	// Remove contexts
	for _, context := range m.RemoveContexts {
		context = strings.TrimSpace(context)
		if context == "" {
			continue
		}
		
		// Remove leading @ if present
		if strings.HasPrefix(context, "@") {
			context = context[1:]
		}
		
		// Remove context from list
		for i, c := range task.Contexts {
			if c == context {
				task.Contexts = append(task.Contexts[:i], task.Contexts[i+1:]...)
				break
			}
		}
	}
	
	return task
}

// DueDateModifier modifies the due date of a task
type DueDateModifier struct {
	DueDate     *time.Time
	RemoveDueDate bool
}

func (m DueDateModifier) ModifyTask(task Task) Task {
	if m.RemoveDueDate {
		task.DueDate = time.Time{}
	} else if m.DueDate != nil {
		task.DueDate = *m.DueDate
	}
	return task
}

// TextModifier appends or prepends text to a task
type TextModifier struct {
	AppendText  string
	PrependText string
}

func (m TextModifier) ModifyTask(task Task) Task {
	if m.AppendText != "" {
		task.Todo = task.Todo + " " + m.AppendText
	}
	if m.PrependText != "" {
		task.Todo = m.PrependText + " " + task.Todo
	}
	return task
}

// CompositeModifier applies multiple modifiers in sequence
type CompositeModifier struct {
	Modifiers []TaskModifier
}

func (m CompositeModifier) ModifyTask(task Task) Task {
	for _, modifier := range m.Modifiers {
		task = modifier.ModifyTask(task)
	}
	return task
}

// ApplyModifiers applies a list of modifiers to a task
func ApplyModifiers(task Task, modifiers []TaskModifier) Task {
	for _, modifier := range modifiers {
		task = modifier.ModifyTask(task)
	}
	return task
}