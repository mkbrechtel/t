package cli

import (
	"log"

	"t5.mkbrechtel.dev/t5/core"
	todo "t5.mkbrechtel.dev/t5/sync/todotxt"
)

// UpdateTasks updates and ensures properties of tasks in your todo list
// Follows the t5 update todo [file] pattern
func UpdateTasks(ctx *AppContext) {
	// Read tasks from todo.txt file
	taskList, err := todo.ReadTodoFile(ctx.Config.TodoFile)
	if err != nil {
		log.Fatalf("Failed to read todo file: %v", err)
	}

	config := todo.DefaultEnsureConfig
	// Override with AppContext config
	config.PreferShortIDs = ctx.Config.PreferShortIds
	config.EnforceCreationDate = ctx.Config.EnforceCreationDate
	config.EnforceCompletionDate = ctx.Config.EnforceCompletionDate
	
	// Set default tags that should be applied to all tasks
	config.DefaultTags = map[string]string{
		// "app":     "t",
		// "version": "1.0",
	}

	taskList = todo.EnsureTaskListProperties(taskList, config)

	// Write updates back to todo.txt file
	if err = todo.WriteTodoFile(taskList, ctx.Config.TodoFile); err != nil {
		log.Fatalf("Failed to write todo file: %v", err)
	}

	// Create TodoTxtTaskUpdate event with the content of the todo.txt file
	content, err := todo.GetTodoFileContent(taskList)
	if err != nil {
		log.Fatalf("Failed to get todo file content: %v", err)
	}

	event := core.NewTodoTxtTaskUpdate(content, ctx.Config.TodoFile)
	err = ctx.Repository.SaveEvent(event)
	if err != nil {
		log.Fatalf("Failed to save event: %v", err)
	}

	log.Println("Task update saved to event store")
}