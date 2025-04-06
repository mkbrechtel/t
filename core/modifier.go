package core

import (
	"fmt"
	"strings"
	"time"
)

// TaskModifier interface for modifying tasks
type TaskModifier interface {
	ModifyTask(Task) Task
	// IsActive returns true if the modifier has been initialized with values
	IsActive() bool
}

// CompletionModifier modifies the completion status of a task
type CompletionModifier struct {
	Complete   bool // Flag for --complete
	Uncomplete bool // Flag for --uncomplete
}

func (m CompletionModifier) ModifyTask(task Task) Task {
	// If both flags are set, the last one wins
	if m.Complete {
		task.Completed = true
		task.CompletedDate = time.Now()
	} else if m.Uncomplete {
		task.Completed = false
		task.CompletedDate = time.Time{}
	}
	return task
}

func (m CompletionModifier) IsActive() bool {
	return m.Complete || m.Uncomplete
}

// PriorityModifier modifies the priority of a task
type PriorityModifier struct {
	Priority       string
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

func (m PriorityModifier) IsActive() bool {
	return m.Priority != "" || m.RemovePriority
}

// Implement flag.Value interface
func (m *PriorityModifier) String() string {
	return m.Priority
}

func (m *PriorityModifier) Set(value string) error {
	if value == "" {
		return nil
	}
	m.Priority = strings.ToUpper(strings.TrimSpace(value))
	return nil
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

func (m ProjectModifier) IsActive() bool {
	return len(m.AddProjects) > 0 || len(m.RemoveProjects) > 0
}

// AddProjectFlag implements flag.Value for adding projects
type AddProjectFlag struct {
	ProjectModifier *ProjectModifier
}

func (f *AddProjectFlag) String() string {
	return strings.Join(f.ProjectModifier.AddProjects, ",")
}

func (f *AddProjectFlag) Set(value string) error {
	if value == "" {
		return nil
	}
	
	projects := strings.Split(value, ",")
	for i, project := range projects {
		project = strings.TrimSpace(project)
		// Remove leading + if present
		if strings.HasPrefix(project, "+") {
			project = project[1:]
		}
		projects[i] = project
	}
	
	f.ProjectModifier.AddProjects = projects
	return nil
}

// RemoveProjectFlag implements flag.Value for removing projects
type RemoveProjectFlag struct {
	ProjectModifier *ProjectModifier
}

func (f *RemoveProjectFlag) String() string {
	return strings.Join(f.ProjectModifier.RemoveProjects, ",")
}

func (f *RemoveProjectFlag) Set(value string) error {
	if value == "" {
		return nil
	}
	
	projects := strings.Split(value, ",")
	for i, project := range projects {
		project = strings.TrimSpace(project)
		// Remove leading + if present
		if strings.HasPrefix(project, "+") {
			project = project[1:]
		}
		projects[i] = project
	}
	
	f.ProjectModifier.RemoveProjects = projects
	return nil
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

func (m ContextModifier) IsActive() bool {
	return len(m.AddContexts) > 0 || len(m.RemoveContexts) > 0
}

// AddContextFlag implements flag.Value for adding contexts
type AddContextFlag struct {
	ContextModifier *ContextModifier
}

func (f *AddContextFlag) String() string {
	return strings.Join(f.ContextModifier.AddContexts, ",")
}

func (f *AddContextFlag) Set(value string) error {
	if value == "" {
		return nil
	}
	
	contexts := strings.Split(value, ",")
	for i, context := range contexts {
		context = strings.TrimSpace(context)
		// Remove leading @ if present
		if strings.HasPrefix(context, "@") {
			context = context[1:]
		}
		contexts[i] = context
	}
	
	f.ContextModifier.AddContexts = contexts
	return nil
}

// RemoveContextFlag implements flag.Value for removing contexts
type RemoveContextFlag struct {
	ContextModifier *ContextModifier
}

func (f *RemoveContextFlag) String() string {
	return strings.Join(f.ContextModifier.RemoveContexts, ",")
}

func (f *RemoveContextFlag) Set(value string) error {
	if value == "" {
		return nil
	}
	
	contexts := strings.Split(value, ",")
	for i, context := range contexts {
		context = strings.TrimSpace(context)
		// Remove leading @ if present
		if strings.HasPrefix(context, "@") {
			context = context[1:]
		}
		contexts[i] = context
	}
	
	f.ContextModifier.RemoveContexts = contexts
	return nil
}

// DueDateModifier modifies the due date of a task
type DueDateModifier struct {
	DueDate       *time.Time
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

func (m DueDateModifier) IsActive() bool {
	return m.DueDate != nil || m.RemoveDueDate
}

// Implement flag.Value interface
func (m *DueDateModifier) String() string {
	if m.DueDate != nil {
		return m.DueDate.Format("2006-01-02")
	}
	return ""
}

func (m *DueDateModifier) Set(value string) error {
	if value == "" {
		return nil
	}
	
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	m.DueDate = &date
	return nil
}

// TextModifier appends or prepends text to a task
type TextModifier struct {
	AppendText  string
	PrependText string
}

func (m TextModifier) ModifyTask(task Task) Task {
	if m.AppendText != "" {
		task.Todo = task.Todo + m.AppendText
	}
	if m.PrependText != "" {
		task.Todo = m.PrependText + task.Todo
	}
	return task
}

func (m TextModifier) IsActive() bool {
	return m.AppendText != "" || m.PrependText != ""
}

// AppendTextFlag implements flag.Value for appending text
type AppendTextFlag struct {
	TextModifier *TextModifier
}

func (f *AppendTextFlag) String() string {
	return f.TextModifier.AppendText
}

func (f *AppendTextFlag) Set(value string) error {
	f.TextModifier.AppendText = value
	return nil
}

// PrependTextFlag implements flag.Value for prepending text
type PrependTextFlag struct {
	TextModifier *TextModifier
}

func (f *PrependTextFlag) String() string {
	return f.TextModifier.PrependText
}

func (f *PrependTextFlag) Set(value string) error {
	f.TextModifier.PrependText = value
	return nil
}

// ChainModifier applies multiple modifiers in sequence
type ChainModifier struct {
	Modifiers []TaskModifier
}

func (m ChainModifier) ModifyTask(task Task) Task {
	for _, modifier := range m.Modifiers {
		task = modifier.ModifyTask(task)
	}
	return task
}

func (m ChainModifier) IsActive() bool {
	return len(m.Modifiers) > 0
}

// AddModifier adds a modifier to the ChainModifier if it's active
func (m *ChainModifier) AddModifier(modifier TaskModifier) {
	if modifier != nil && modifier.IsActive() {
		m.Modifiers = append(m.Modifiers, modifier)
	}
}

// ModifierComposer manages a collection of modifiers and builds a ChainModifier
// It integrates with flag parsing to only include modifiers that are actually selected
type ModifierComposer struct {
	modifiers []TaskModifier
}

// NewModifierComposer creates a new ModifierComposer
func NewModifierComposer() *ModifierComposer {
	return &ModifierComposer{modifiers: make([]TaskModifier, 0)}
}

// AddModifier adds a modifier to the composer
// Only active modifiers will be included in the final composition
func (mc *ModifierComposer) AddModifier(modifier TaskModifier) {
	if modifier != nil {
		mc.modifiers = append(mc.modifiers, modifier)
	}
}

// ComposeModifier creates a ChainModifier containing all active modifiers
// If no modifiers are active, it returns nil
// If only one modifier is active, it returns that modifier directly
func (mc *ModifierComposer) ComposeModifier() TaskModifier {
	activeModifiers := make([]TaskModifier, 0)
	
	for _, modifier := range mc.modifiers {
		if modifier.IsActive() {
			activeModifiers = append(activeModifiers, modifier)
		}
	}
	
	if len(activeModifiers) == 0 {
		return nil
	}
	
	if len(activeModifiers) == 1 {
		return activeModifiers[0]
	}
	
	return &ChainModifier{Modifiers: activeModifiers}
}

// ApplyModifier applies a modifier to a task
func ApplyModifier(task Task, modifier TaskModifier) Task {
	if modifier == nil {
		return task
	}
	return modifier.ModifyTask(task)
}