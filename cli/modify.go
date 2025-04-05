package cli

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"t5.mkbrechtel.dev/core"
	todo "t5.mkbrechtel.dev/sync/todotxt"
)

// TaskModifierFlags holds flag values for modifying tasks
type TaskModifierFlags struct {
	Complete       bool
	Uncomplete     bool
	Priority       string
	RemovePriority bool
	AddProject     string
	RemoveProject  string
	AddContext     string
	RemoveContext  string
	Due            string
	RemoveDue      bool
	Append         string
	Prepend        string
}

// ModifyTask applies modifications to a task identified by ID
// Follows the t5 modify <task-id> pattern where task-id is a direct object
func ModifyTask(ctx *AppContext, taskID string) {
	// Parse modifier flags
	modifierFlags := parseModifierFlags(ctx)
	
	// Find the task by ID (either UUID or short ID)
	var task core.Task
	var err error
	
	// Try to parse as UUID first
	parsedUUID := uuid.FromStringOrNil(taskID)
	if parsedUUID != uuid.Nil {
		task, err = ctx.Repository.GetTask(parsedUUID)
	} else {
		// If not a UUID, search for task with matching short ID
		tasks, err := ctx.Repository.ListTasks(nil)
		if err != nil {
			log.Fatalf("Failed to list tasks: %v", err)
		}
		
		found := false
		for _, t := range tasks {
			if short, ok := t.AdditionalTags["short"]; ok && short == taskID {
				task = t
				found = true
				break
			}
		}
		
		if !found {
			log.Fatalf("Task with ID %s not found", taskID)
		}
	}
	
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}
	
	// Build modifiers based on flags
	modifiers := buildModifiers(modifierFlags)
	
	// Apply modifications
	modified := core.ApplyModifiers(task, modifiers)
	
	// Update the task
	err = ctx.Repository.UpdateTask(modified)
	if err != nil {
		log.Fatalf("Failed to update task: %v", err)
	}
	
	// Read and update todo.txt file
	taskList, err := todo.ReadTodoFile(ctx.Config.TodoFile)
	if err != nil {
		log.Fatalf("Failed to read todo file: %v", err)
	}
	
	// Find and update the task in the task list
	found := false
	for i, t := range taskList {
		if uuidTag, ok := t.AdditionalTags["uuid"]; ok && uuidTag == modified.ID.String() {
			// Convert the modified task to todo.txt format
			newTask, err := todo.ParseTask(modified.ToTodoTxt())
			if err != nil {
				log.Fatalf("Failed to parse modified task: %v", err)
			}
			
			taskList[i] = *newTask
			found = true
			break
		}
	}
	
	if !found {
		log.Printf("Warning: Task not found in todo.txt file. Adding it now.")
		newTask, err := todo.ParseTask(modified.ToTodoTxt())
		if err != nil {
			log.Fatalf("Failed to parse modified task: %v", err)
		}
		taskList = append(taskList, *newTask)
	}
	
	// Apply the configured settings
	config := todo.DefaultEnsureConfig
	config.PreferShortIDs = ctx.Config.PreferShortIds
	config.EnforceCreationDate = ctx.Config.EnforceCreationDate
	config.EnforceCompletionDate = ctx.Config.EnforceCompletionDate
	
	taskList = todo.EnsureTaskListProperties(taskList, config)
	
	// Write updated tasks back to todo.txt file
	if err = todo.WriteTodoFile(taskList, ctx.Config.TodoFile); err != nil {
		log.Fatalf("Failed to write todo file: %v", err)
	}
	
	fmt.Println("Task updated successfully:")
	fmt.Println(modified.ToTodoTxt())
}

// parseModifierFlags extracts modifier flags from the command context
func parseModifierFlags(ctx *AppContext) TaskModifierFlags {
	modifierFlags := TaskModifierFlags{}
	
	// Add modifier flags to the flagset
	ctx.FlagSet.BoolVar(&modifierFlags.Complete, "complete", false, "Mark task as complete")
	ctx.FlagSet.BoolVar(&modifierFlags.Uncomplete, "uncomplete", false, "Mark task as incomplete")
	ctx.FlagSet.StringVar(&modifierFlags.Priority, "priority", "", "Set task priority")
	ctx.FlagSet.BoolVar(&modifierFlags.RemovePriority, "remove-priority", false, "Remove task priority")
	ctx.FlagSet.StringVar(&modifierFlags.AddProject, "add-project", "", "Add project to task")
	ctx.FlagSet.StringVar(&modifierFlags.RemoveProject, "remove-project", "", "Remove project from task")
	ctx.FlagSet.StringVar(&modifierFlags.AddContext, "add-context", "", "Add context to task")
	ctx.FlagSet.StringVar(&modifierFlags.RemoveContext, "remove-context", "", "Remove context from task")
	ctx.FlagSet.StringVar(&modifierFlags.Due, "due", "", "Set due date (YYYY-MM-DD)")
	ctx.FlagSet.BoolVar(&modifierFlags.RemoveDue, "remove-due", false, "Remove due date")
	ctx.FlagSet.StringVar(&modifierFlags.Append, "append", "", "Append text to task")
	ctx.FlagSet.StringVar(&modifierFlags.Prepend, "prepend", "", "Prepend text to task")
	
	return modifierFlags
}

// buildModifiers creates task modifiers based on flag values
func buildModifiers(flags TaskModifierFlags) []core.TaskModifier {
	var modifiers []core.TaskModifier
	
	// Handle completion status
	if flags.Complete {
		modifiers = append(modifiers, core.CompletionModifier{Completed: true})
	} else if flags.Uncomplete {
		modifiers = append(modifiers, core.CompletionModifier{Completed: false})
	}
	
	// Handle priority
	if flags.RemovePriority || flags.Priority != "" {
		modifiers = append(modifiers, core.PriorityModifier{
			Priority:       flags.Priority,
			RemovePriority: flags.RemovePriority,
		})
	}
	
	// Handle projects
	if flags.AddProject != "" || flags.RemoveProject != "" {
		var addProjects []string
		var removeProjects []string
		
		if flags.AddProject != "" {
			addProjects = strings.Split(flags.AddProject, ",")
		}
		
		if flags.RemoveProject != "" {
			removeProjects = strings.Split(flags.RemoveProject, ",")
		}
		
		modifiers = append(modifiers, core.ProjectModifier{
			AddProjects:    addProjects,
			RemoveProjects: removeProjects,
		})
	}
	
	// Handle contexts
	if flags.AddContext != "" || flags.RemoveContext != "" {
		var addContexts []string
		var removeContexts []string
		
		if flags.AddContext != "" {
			addContexts = strings.Split(flags.AddContext, ",")
		}
		
		if flags.RemoveContext != "" {
			removeContexts = strings.Split(flags.RemoveContext, ",")
		}
		
		modifiers = append(modifiers, core.ContextModifier{
			AddContexts:    addContexts,
			RemoveContexts: removeContexts,
		})
	}
	
	// Handle due date
	if flags.RemoveDue || flags.Due != "" {
		var dueDate *time.Time
		
		if flags.Due != "" {
			if date, err := time.Parse("2006-01-02", flags.Due); err == nil {
				dueDate = &date
			} else {
				log.Printf("Warning: Invalid date format for due: %s", flags.Due)
			}
		}
		
		modifiers = append(modifiers, core.DueDateModifier{
			DueDate:      dueDate,
			RemoveDueDate: flags.RemoveDue,
		})
	}
	
	// Handle text modifications
	if flags.Append != "" || flags.Prepend != "" {
		modifiers = append(modifiers, core.TextModifier{
			AppendText:  flags.Append,
			PrependText: flags.Prepend,
		})
	}
	
	return modifiers
}