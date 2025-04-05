package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"
	
	"t5.mkbrechtel.dev/core"
)

// TaskFilterFlags holds flag values for filtering tasks
type TaskFilterFlags struct {
	Completed     *bool
	NotCompleted  *bool
	Priority      string
	Project       string
	Context       string
	DueToday      bool
	DueBefore     string
	CreatedAfter  string
	HasTag        string
	Regex         string
}

// ListTasks lists tasks from the event store with optional filters
// Follows the t5 list todo pattern
func ListTasks(ctx *AppContext) {
	// Create filter flags
	filterFlags := parseFilterFlags(ctx)
	
	// Build filters based on flags
	filters := buildFilters(filterFlags)
	
	// Get filtered tasks
	tasks, err := ctx.Repository.ListTasks(filters)
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

// parseFilterFlags extracts filter flags from the command context
func parseFilterFlags(ctx *AppContext) TaskFilterFlags {
	filterFlags := TaskFilterFlags{}
	
	// Add filter flags to the flagset
	completed := ctx.FlagSet.Bool("completed", false, "Show only completed tasks")
	notCompleted := ctx.FlagSet.Bool("not-completed", false, "Show only non-completed tasks")
	ctx.FlagSet.StringVar(&filterFlags.Priority, "priority", "", "Filter by priority (e.g., A,B,C)")
	ctx.FlagSet.StringVar(&filterFlags.Project, "project", "", "Filter by project")
	ctx.FlagSet.StringVar(&filterFlags.Context, "context", "", "Filter by context")
	ctx.FlagSet.BoolVar(&filterFlags.DueToday, "due-today", false, "Show tasks due today")
	ctx.FlagSet.StringVar(&filterFlags.DueBefore, "due-before", "", "Show tasks due before a date (YYYY-MM-DD)")
	ctx.FlagSet.StringVar(&filterFlags.CreatedAfter, "created-after", "", "Show tasks created after a date (YYYY-MM-DD)")
	ctx.FlagSet.StringVar(&filterFlags.HasTag, "has-tag", "", "Show tasks with a specific tag")
	ctx.FlagSet.StringVar(&filterFlags.Regex, "regex", "", "Filter tasks using a regular expression")
	
	// Handle the completed/not-completed flags manually to avoid conflicting default values
	if ctx.IsFlagPassed("completed") {
		filterFlags.Completed = completed
	}
	
	if ctx.IsFlagPassed("not-completed") {
		filterFlags.NotCompleted = notCompleted
	}
	
	return filterFlags
}

// buildFilters creates task filters based on flag values
func buildFilters(flags TaskFilterFlags) []core.TaskFilter {
	var filters []core.TaskFilter
	
	// Completion status filter
	if flags.Completed != nil && *flags.Completed {
		completed := true
		filters = append(filters, core.CompletionFilter{Completed: &completed})
	}
	
	if flags.NotCompleted != nil && *flags.NotCompleted {
		completed := false
		filters = append(filters, core.CompletionFilter{Completed: &completed})
	}
	
	// Priority filter
	if flags.Priority != "" {
		priorities := strings.Split(flags.Priority, ",")
		priorityFilters := make([]core.TaskFilter, 0, len(priorities))
		
		for _, p := range priorities {
			p = strings.TrimSpace(p)
			if p != "" {
				priorityFilters = append(priorityFilters, core.PriorityFilter{Priority: p})
			}
		}
		
		// If we have multiple priorities, use OR filter to match any of them
		if len(priorityFilters) > 1 {
			filters = append(filters, core.OrFilter{Filters: priorityFilters})
		} else if len(priorityFilters) == 1 {
			filters = append(filters, priorityFilters[0])
		}
	}
	
	// Project filter
	if flags.Project != "" {
		projects := strings.Split(flags.Project, ",")
		filters = append(filters, core.ProjectFilter{Projects: projects})
	}
	
	// Context filter
	if flags.Context != "" {
		contexts := strings.Split(flags.Context, ",")
		filters = append(filters, core.ContextFilter{Contexts: contexts})
	}
	
	// Due date filters
	if flags.DueToday {
		today := time.Now()
		// Set time to end of day
		tomorrow := time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, today.Location())
		
		filters = append(filters, core.DueDateFilter{
			Before: &tomorrow,
			After:  &today,
		})
	}
	
	if flags.DueBefore != "" {
		if date, err := time.Parse("2006-01-02", flags.DueBefore); err == nil {
			filters = append(filters, core.DueDateFilter{
				Before: &date,
			})
		} else {
			log.Printf("Warning: Invalid date format for due-before: %s", flags.DueBefore)
		}
	}
	
	// Created date filter
	if flags.CreatedAfter != "" {
		if date, err := time.Parse("2006-01-02", flags.CreatedAfter); err == nil {
			filters = append(filters, core.CreatedDateFilter{
				After: &date,
			})
		} else {
			log.Printf("Warning: Invalid date format for created-after: %s", flags.CreatedAfter)
		}
	}
	
	// Tag filter
	if flags.HasTag != "" {
		filters = append(filters, core.TagFilter{Tag: flags.HasTag})
	}
	
	// Regex filter
	if flags.Regex != "" {
		filters = append(filters, core.RegexFilter{Pattern: flags.Regex})
	}
	
	// If we have multiple filters, wrap them in an AND filter
	if len(filters) > 1 {
		return []core.TaskFilter{core.AndFilter{Filters: filters}}
	}
	
	return filters
}