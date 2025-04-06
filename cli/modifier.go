package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"t5.mkbrechtel.dev/t5/core"
)

// ModifierComposer combines multiple modifiers and provides flag.Value
// implementations for easy flag integration
type ModifierComposer struct {
	Modifiers     core.ChainModifier
	CompleteFlag  CompleteFlag
	UncompleteFlag UncompleteFlag
	PriorityFlag  PriorityModifierFlag
	RemovePriorityFlag RemovePriorityFlag
	AddProjectFlag AddProjectFlag
	RemoveProjectFlag RemoveProjectFlag
	AddContextFlag AddContextFlag
	RemoveContextFlag RemoveContextFlag
	DueFlag       DueFlag
	RemoveDueFlag RemoveDueFlag
	AppendFlag    AppendFlag
	PrependFlag   PrependFlag
}

// NewModifierComposer creates a new ModifierComposer
func NewModifierComposer() *ModifierComposer {
	mc := &ModifierComposer{
		Modifiers: core.ChainModifier{Modifiers: make([]core.TaskModifier, 0)},
	}
	
	// Initialize flags with parent reference
	mc.CompleteFlag.parent = mc
	mc.UncompleteFlag.parent = mc
	mc.PriorityFlag.parent = mc
	mc.RemovePriorityFlag.parent = mc
	mc.AddProjectFlag.parent = mc
	mc.RemoveProjectFlag.parent = mc
	mc.AddContextFlag.parent = mc
	mc.RemoveContextFlag.parent = mc
	mc.DueFlag.parent = mc
	mc.RemoveDueFlag.parent = mc
	mc.AppendFlag.parent = mc
	mc.PrependFlag.parent = mc
	
	return mc
}

// AddModifierFlags adds all modifier flags to the provided flagset
func (mc *ModifierComposer) AddModifierFlags(flagset *flag.FlagSet) {
	flagset.Var(&mc.CompleteFlag, "complete", "Mark task as complete")
	flagset.Var(&mc.UncompleteFlag, "uncomplete", "Mark task as incomplete")
	flagset.Var(&mc.PriorityFlag, "priority", "Set task priority")
	flagset.Var(&mc.RemovePriorityFlag, "remove-priority", "Remove task priority")
	flagset.Var(&mc.AddProjectFlag, "add-project", "Add project to task")
	flagset.Var(&mc.RemoveProjectFlag, "remove-project", "Remove project from task")
	flagset.Var(&mc.AddContextFlag, "add-context", "Add context to task")
	flagset.Var(&mc.RemoveContextFlag, "remove-context", "Remove context from task")
	flagset.Var(&mc.DueFlag, "due", "Set due date (YYYY-MM-DD)")
	flagset.Var(&mc.RemoveDueFlag, "remove-due", "Remove due date")
	flagset.Var(&mc.AppendFlag, "append", "Append text to task")
	flagset.Var(&mc.PrependFlag, "prepend", "Prepend text to task")
}

// ProcessOriginalArgs has been removed as it's no longer needed.
// Flag parsing is now handled directly by each command with its own flag set.


// AddModifier adds a modifier to the composer
func (mc *ModifierComposer) AddModifier(modifier core.TaskModifier) {
	if modifier != nil {
		// Always add the modifier, since we may add values to it after creating it
		mc.Modifiers.Modifiers = append(mc.Modifiers.Modifiers, modifier)
	}
}

// ComposeModifier creates a modifier for use with task modification
func (mc *ModifierComposer) ComposeModifier() core.TaskModifier {
	// Filter out inactive modifiers
	activeModifiers := make([]core.TaskModifier, 0, len(mc.Modifiers.Modifiers))
	for _, mod := range mc.Modifiers.Modifiers {
		if mod.IsActive() {
			activeModifiers = append(activeModifiers, mod)
		}
	}
	
	if len(activeModifiers) == 0 {
		return nil
	}
	
	if len(activeModifiers) == 1 {
		return activeModifiers[0]
	}
	
	// Create a new chain with only active modifiers
	return &core.ChainModifier{Modifiers: activeModifiers}
}

// CompleteFlag implements flag.Value for marking a task as complete
type CompleteFlag struct {
	parent *ModifierComposer
	value  bool
}

func (f *CompleteFlag) String() string {
	return fmt.Sprintf("%v", f.value)
}

func (f *CompleteFlag) Set(value string) error {
	// For boolean flags, presence of the flag is enough
	f.value = true
	
	// Create complete modifier
	modifier := &core.CompletionModifier{Complete: true}
	f.parent.AddModifier(modifier)
	
	return nil
}

// UncompleteFlag implements flag.Value for marking a task as incomplete
type UncompleteFlag struct {
	parent *ModifierComposer
	value  bool
}

func (f *UncompleteFlag) String() string {
	return fmt.Sprintf("%v", f.value)
}

func (f *UncompleteFlag) Set(value string) error {
	// For boolean flags, presence of the flag is enough
	f.value = true
	
	// Create uncomplete modifier
	modifier := &core.CompletionModifier{Uncomplete: true}
	f.parent.AddModifier(modifier)
	
	return nil
}

// PriorityModifierFlag implements flag.Value for setting priority
type PriorityModifierFlag struct {
	parent *ModifierComposer
	value  string
}

func (f *PriorityModifierFlag) String() string {
	return f.value
}

func (f *PriorityModifierFlag) Set(value string) error {
	f.value = value
	
	if value == "" {
		return nil
	}
	
	// Create priority modifier
	modifier := &core.PriorityModifier{Priority: strings.ToUpper(strings.TrimSpace(value))}
	f.parent.AddModifier(modifier)
	
	return nil
}

// RemovePriorityFlag implements flag.Value for removing priority
type RemovePriorityFlag struct {
	parent *ModifierComposer
	value  bool
}

func (f *RemovePriorityFlag) String() string {
	return fmt.Sprintf("%v", f.value)
}

func (f *RemovePriorityFlag) Set(value string) error {
	// For boolean flags, presence of the flag is enough
	f.value = true
	
	// Create remove priority modifier
	modifier := &core.PriorityModifier{RemovePriority: true}
	f.parent.AddModifier(modifier)
	
	return nil
}

// AddProjectFlag implements flag.Value for adding projects
type AddProjectFlag struct {
	parent *ModifierComposer
	value  string
}

func (f *AddProjectFlag) String() string {
	return f.value
}

func (f *AddProjectFlag) Set(value string) error {
	f.value = value
	
	if value == "" {
		return nil
	}
	
	// Process projects
	projects := strings.Split(value, ",")
	for i, project := range projects {
		project = strings.TrimSpace(project)
		// Remove leading + if present
		if strings.HasPrefix(project, "+") {
			project = project[1:]
		}
		projects[i] = project
	}
	
	// Check if there's already a project modifier
	var projectModifier *core.ProjectModifier
	
	// Look for existing project modifier
	for _, mod := range f.parent.Modifiers.Modifiers {
		if pm, ok := mod.(*core.ProjectModifier); ok {
			projectModifier = pm
			break
		}
	}
	
	// If no project modifier exists, create a new one
	if projectModifier == nil {
		projectModifier = &core.ProjectModifier{}
		f.parent.AddModifier(projectModifier)
	}
	
	// Add projects to the modifier
	projectModifier.AddProjects = append(projectModifier.AddProjects, projects...)
	
	return nil
}

// RemoveProjectFlag implements flag.Value for removing projects
type RemoveProjectFlag struct {
	parent *ModifierComposer
	value  string
}

func (f *RemoveProjectFlag) String() string {
	return f.value
}

func (f *RemoveProjectFlag) Set(value string) error {
	f.value = value
	
	if value == "" {
		return nil
	}
	
	// Process projects
	projects := strings.Split(value, ",")
	for i, project := range projects {
		project = strings.TrimSpace(project)
		// Remove leading + if present
		if strings.HasPrefix(project, "+") {
			project = project[1:]
		}
		projects[i] = project
	}
	
	// Check if there's already a project modifier
	var projectModifier *core.ProjectModifier
	
	// Look for existing project modifier
	for _, mod := range f.parent.Modifiers.Modifiers {
		if pm, ok := mod.(*core.ProjectModifier); ok {
			projectModifier = pm
			break
		}
	}
	
	// If no project modifier exists, create a new one
	if projectModifier == nil {
		projectModifier = &core.ProjectModifier{}
		f.parent.AddModifier(projectModifier)
	}
	
	// Add remove projects to the modifier
	projectModifier.RemoveProjects = append(projectModifier.RemoveProjects, projects...)
	
	return nil
}

// AddContextFlag implements flag.Value for adding contexts
type AddContextFlag struct {
	parent *ModifierComposer
	value  string
}

func (f *AddContextFlag) String() string {
	return f.value
}

func (f *AddContextFlag) Set(value string) error {
	f.value = value
	
	if value == "" {
		return nil
	}
	
	// Process contexts
	contexts := strings.Split(value, ",")
	for i, context := range contexts {
		context = strings.TrimSpace(context)
		// Remove leading @ if present
		if strings.HasPrefix(context, "@") {
			context = context[1:]
		}
		contexts[i] = context
	}
	
	// Check if there's already a context modifier
	var contextModifier *core.ContextModifier
	
	// Look for existing context modifier
	for _, mod := range f.parent.Modifiers.Modifiers {
		if cm, ok := mod.(*core.ContextModifier); ok {
			contextModifier = cm
			break
		}
	}
	
	// If no context modifier exists, create a new one
	if contextModifier == nil {
		contextModifier = &core.ContextModifier{}
		f.parent.AddModifier(contextModifier)
	}
	
	// Add contexts to the modifier
	contextModifier.AddContexts = append(contextModifier.AddContexts, contexts...)
	
	return nil
}

// RemoveContextFlag implements flag.Value for removing contexts
type RemoveContextFlag struct {
	parent *ModifierComposer
	value  string
}

func (f *RemoveContextFlag) String() string {
	return f.value
}

func (f *RemoveContextFlag) Set(value string) error {
	f.value = value
	
	if value == "" {
		return nil
	}
	
	// Process contexts
	contexts := strings.Split(value, ",")
	for i, context := range contexts {
		context = strings.TrimSpace(context)
		// Remove leading @ if present
		if strings.HasPrefix(context, "@") {
			context = context[1:]
		}
		contexts[i] = context
	}
	
	// Check if there's already a context modifier
	var contextModifier *core.ContextModifier
	
	// Look for existing context modifier
	for _, mod := range f.parent.Modifiers.Modifiers {
		if cm, ok := mod.(*core.ContextModifier); ok {
			contextModifier = cm
			break
		}
	}
	
	// If no context modifier exists, create a new one
	if contextModifier == nil {
		contextModifier = &core.ContextModifier{}
		f.parent.AddModifier(contextModifier)
	}
	
	// Add remove contexts to the modifier
	contextModifier.RemoveContexts = append(contextModifier.RemoveContexts, contexts...)
	
	return nil
}

// DueFlag implements flag.Value for setting due date
type DueFlag struct {
	parent *ModifierComposer
	value  string
}

func (f *DueFlag) String() string {
	return f.value
}

func (f *DueFlag) Set(value string) error {
	f.value = value
	
	if value == "" {
		return nil
	}
	
	// Parse due date
	dueDate, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	
	// Create due date modifier
	modifier := &core.DueDateModifier{DueDate: &dueDate}
	f.parent.AddModifier(modifier)
	
	return nil
}

// RemoveDueFlag implements flag.Value for removing due date
type RemoveDueFlag struct {
	parent *ModifierComposer
	value  bool
}

func (f *RemoveDueFlag) String() string {
	return fmt.Sprintf("%v", f.value)
}

func (f *RemoveDueFlag) Set(value string) error {
	// For boolean flags, presence of the flag is enough
	f.value = true
	
	// Create remove due date modifier
	modifier := &core.DueDateModifier{RemoveDueDate: true}
	f.parent.AddModifier(modifier)
	
	return nil
}

// AppendFlag implements flag.Value for appending text
type AppendFlag struct {
	parent *ModifierComposer
	value  string
}

func (f *AppendFlag) String() string {
	return f.value
}

func (f *AppendFlag) Set(value string) error {
	f.value = value
	
	if value == "" {
		return nil
	}
	
	// Check if there's already a text modifier
	var textModifier *core.TextModifier
	
	// Look for existing text modifier
	for _, mod := range f.parent.Modifiers.Modifiers {
		if tm, ok := mod.(*core.TextModifier); ok {
			textModifier = tm
			break
		}
	}
	
	// If no text modifier exists, create a new one
	if textModifier == nil {
		textModifier = &core.TextModifier{}
		f.parent.AddModifier(textModifier)
	}
	
	// Set append text
	textModifier.AppendText = value
	
	return nil
}

// PrependFlag implements flag.Value for prepending text
type PrependFlag struct {
	parent *ModifierComposer
	value  string
}

func (f *PrependFlag) String() string {
	return f.value
}

func (f *PrependFlag) Set(value string) error {
	f.value = value
	
	if value == "" {
		return nil
	}
	
	// Check if there's already a text modifier
	var textModifier *core.TextModifier
	
	// Look for existing text modifier
	for _, mod := range f.parent.Modifiers.Modifiers {
		if tm, ok := mod.(*core.TextModifier); ok {
			textModifier = tm
			break
		}
	}
	
	// If no text modifier exists, create a new one
	if textModifier == nil {
		textModifier = &core.TextModifier{}
		f.parent.AddModifier(textModifier)
	}
	
	// Set prepend text
	textModifier.PrependText = value
	
	return nil
}