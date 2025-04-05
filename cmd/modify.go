package cmd

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
	
	// Apply modifications
	modified := applyModifications(task, modifierFlags)
	
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

// applyModifications applies the requested modifications to a task
func applyModifications(task core.Task, flags TaskModifierFlags) core.Task {
	// Handle completion status
	if flags.Complete {
		task.Completed = true
		task.CompletedDate = time.Now()
	} else if flags.Uncomplete {
		task.Completed = false
		task.CompletedDate = time.Time{}
	}
	
	// Handle priority
	if flags.RemovePriority {
		task.Priority = ""
	} else if flags.Priority != "" {
		task.Priority = strings.ToUpper(flags.Priority)
	}
	
	// Handle projects
	if flags.AddProject != "" {
		projects := strings.Split(flags.AddProject, ",")
		for _, project := range projects {
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
	}
	
	if flags.RemoveProject != "" {
		projects := strings.Split(flags.RemoveProject, ",")
		for _, project := range projects {
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
	}
	
	// Handle contexts
	if flags.AddContext != "" {
		contexts := strings.Split(flags.AddContext, ",")
		for _, context := range contexts {
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
	}
	
	if flags.RemoveContext != "" {
		contexts := strings.Split(flags.RemoveContext, ",")
		for _, context := range contexts {
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
	}
	
	// Handle due date
	if flags.RemoveDue {
		task.DueDate = time.Time{}
	} else if flags.Due != "" {
		if date, err := time.Parse("2006-01-02", flags.Due); err == nil {
			task.DueDate = date
		} else {
			log.Printf("Warning: Invalid date format for due: %s", flags.Due)
		}
	}
	
	// Handle text modifications
	if flags.Append != "" {
		task.Todo = task.Todo + " " + flags.Append
	}
	
	if flags.Prepend != "" {
		task.Todo = flags.Prepend + " " + task.Todo
	}
	
	return task
}