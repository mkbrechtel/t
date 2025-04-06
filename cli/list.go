package cli

import (
	"flag"
	"fmt"
	"log"
	
	"t5.mkbrechtel.dev/t5/core"
)

// ListTasks lists tasks from the event store with optional filters
// Follows the t5 list todo pattern
func ListTasks(ctx *AppContext, args []string) {
	// Create a dedicated FlagSet for list command
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	
	// Create a filter composer
	filterer := NewFilterComposer()
	
	// Add filter flags to the flagset
	filterer.AddFilterFlags(fs)
	
	// Parse the args for this specific command
	fs.Parse(args)
	
	// Compose the final filter
	taskFilter := filterer.ComposeFilter()
	
	// Convert to slice for the ListTasks function
	var taskFilters []core.TaskFilter
	if taskFilter != nil {
		taskFilters = []core.TaskFilter{taskFilter}
	}
	
	// Get filtered tasks
	tasks, err := ctx.Repository.ListTasks(taskFilters)
	if err != nil {
		log.Fatalf("Failed to list tasks: %v", err)
	}

	// Display tasks
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Printf("Found %d tasks:\n", len(tasks))
	for i, task := range tasks {
		priority := task.Priority
		if priority == "" {
			priority = "-"
		}
		status := " "
		if task.Completed {
			status = "x"
		}
		
		// Format projects and contexts
		projects := ""
		for _, proj := range task.Projects {
			projects += " +" + proj
		}
		
		contexts := ""
		for _, ctx := range task.Contexts {
			contexts += " @" + ctx
		}
		
		// Add any due date if present
		dueDate := ""
		if !task.DueDate.IsZero() {
			dueDate = " t:" + task.DueDate.Format("2006-01-02")
		}
		
		// Print task with optional metadata
		fmt.Printf("%d. [%s] (%s) %s%s%s%s\n", 
			i+1, status, priority, task.Todo, projects, contexts, dueDate)
	}
}
