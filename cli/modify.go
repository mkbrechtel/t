package cli

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"
	"t5.mkbrechtel.dev/t5/core"
	todo "t5.mkbrechtel.dev/t5/sync/todotxt"
)

// ModifyTask applies modifications to a task identified by ID
// Follows the t5 modify <task-id> pattern where task-id is a direct object
func ModifyTask(ctx *AppContext, taskID string) {
	// Create a modifier composer
	modifier := NewModifierComposer()
	
	// Add modifier flags to the flagset
	modifier.AddModifierFlags(ctx.FlagSet)
	
	// Process original args for flags parsed at higher level
	modifier.ProcessOriginalArgs(ctx.OriginalArgs)
	
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
	
	// Compose the final modifier
	taskModifier := modifier.ComposeModifier()
	
	// Apply modifications
	var modified core.Task
	if taskModifier != nil {
		modified = taskModifier.ModifyTask(task)
	} else {
		modified = task
	}
	
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