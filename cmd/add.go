package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"t5.mkbrechtel.dev/core"
	todo "t5.mkbrechtel.dev/sync/todotxt"
)

// AddTask creates a new task from command line or stdin
func AddTask(ctx *AppContext, taskText string, stdin io.ReadCloser) {
	// Check if we should read from stdin (if no taskText provided)
	if taskText == "" {
		reader := bufio.NewReader(stdin)
		var sb strings.Builder
		
		fmt.Println("Enter task text (press Ctrl+D when finished):")
		
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error reading from stdin: %v", err)
			}
			sb.WriteString(line)
		}
		
		taskText = strings.TrimSpace(sb.String())
		if taskText == "" {
			log.Fatalf("No task text provided")
		}
	}

	// Read existing tasks from todo.txt file
	taskList, err := todo.ReadTodoFile(ctx.Config.TodoFile)
	if err != nil {
		log.Fatalf("Failed to read todo file: %v", err)
	}

	// Parse the task text
	task, err := todo.ParseTask(taskText)
	if err != nil {
		log.Fatalf("Failed to parse task: %v", err)
	}

	// Apply the configured settings
	config := todo.DefaultEnsureConfig
	config.PreferShortIDs = ctx.Config.PreferShortIds
	config.EnforceCreationDate = ctx.Config.EnforceCreationDate
	config.EnforceCompletionDate = ctx.Config.EnforceCompletionDate

	// Ensure task has an ID, date, etc.
	taskList = append(taskList, *task)
	taskList = todo.EnsureTaskListProperties(taskList, config)

	// Write updated tasks back to todo.txt file
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

	// Display the added task
	fmt.Println("Task added successfully:")
	fmt.Println(taskList[len(taskList)-1].String())
}