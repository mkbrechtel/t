package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	uuid "github.com/gofrs/uuid/v5"
	"t5.mkbrechtel.dev/t5/core"
	"t5.mkbrechtel.dev/t5/utils"
)

// StartTask starts time tracking for a task
func StartTask(ctx *AppContext, taskID string, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("start", flag.ContinueOnError)
	noteFlag := fs.String("note", "", "Note about the activity")
	fs.Parse(args)

	// Find the task by ID
	id, err := utils.DecodeUUID(taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid task ID: %s\n", err)
		return
	}

	// Check if the task exists
	state := ctx.Repository.GetAppState()
	task, exists := state.Tasks[id]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: Task with ID %s not found\n", taskID)
		return
	}

	// Create and apply the event
	event := core.NewTaskStartTime(id, *noteFlag)
	err = ctx.Repository.SaveEvent(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting task: %s\n", err)
		return
	}

	fmt.Printf("Started task: %s\n", task.Todo)
}

// StopTask stops time tracking for the active task
func StopTask(ctx *AppContext, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("stop", flag.ContinueOnError)
	noteFlag := fs.String("note", "", "Note about the completed work")
	fs.Parse(args)

	// Get the active task
	state := ctx.Repository.GetAppState()
	if state.ActiveTask == nil {
		fmt.Fprintf(os.Stderr, "Error: No active task to stop\n")
		return
	}

	// Save the TaskID before applying the event
	taskID := state.ActiveTask.TaskID

	// Get the task information before stopping
	task, exists := state.Tasks[taskID]
	taskName := "Unknown task"
	if exists {
		taskName = task.Todo
	}

	// Create and apply the event
	event := core.NewTaskEndTime(taskID, *noteFlag)
	err := ctx.Repository.SaveEvent(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stopping task: %s\n", err)
		return
	}

	fmt.Printf("Stopped task: %s\n", taskName)
}

// PauseTask pauses time tracking for the active task
func PauseTask(ctx *AppContext, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("pause", flag.ContinueOnError)
	reasonFlag := fs.String("reason", "", "Reason for pausing")
	fs.Parse(args)

	// Get the active task
	state := ctx.Repository.GetAppState()
	if state.ActiveTask == nil {
		fmt.Fprintf(os.Stderr, "Error: No active task to pause\n")
		return
	}

	// Check if the task is already paused
	if !state.ActiveTask.LastPaused.IsZero() {
		fmt.Fprintf(os.Stderr, "Error: Task is already paused\n")
		return
	}
	
	// Save the TaskID before applying the event
	taskID := state.ActiveTask.TaskID
	
	// Get the task information before pausing
	task, exists := state.Tasks[taskID]
	taskName := "Unknown task"
	if exists {
		taskName = task.Todo
	}

	// Create and apply the event
	event := core.NewTaskPauseTime(taskID, *reasonFlag)
	err := ctx.Repository.SaveEvent(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pausing task: %s\n", err)
		return
	}

	fmt.Printf("Paused task: %s\n", taskName)
}

// ResumeTask resumes time tracking for the paused task
func ResumeTask(ctx *AppContext, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("resume", flag.ContinueOnError)
	fs.Parse(args)

	// Get the active task
	state := ctx.Repository.GetAppState()
	if state.ActiveTask == nil {
		fmt.Fprintf(os.Stderr, "Error: No active task to resume\n")
		return
	}

	// Check if the task is actually paused
	if state.ActiveTask.LastPaused.IsZero() {
		fmt.Fprintf(os.Stderr, "Error: Task is not paused\n")
		return
	}
	
	// Save the TaskID before applying the event
	taskID := state.ActiveTask.TaskID
	
	// Get the task information before resuming
	task, exists := state.Tasks[taskID]
	taskName := "Unknown task"
	if exists {
		taskName = task.Todo
	}

	// Create and apply the event
	event := core.NewTaskResumeTime(taskID)
	err := ctx.Repository.SaveEvent(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resuming task: %s\n", err)
		return
	}

	fmt.Printf("Resumed task: %s\n", taskName)
}

// ShowStatus displays the current status of the active task
func ShowStatus(ctx *AppContext, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.Parse(args)

	// Get the app state
	state := ctx.Repository.GetAppState()
	if state.ActiveTask == nil {
		fmt.Println("No active task")
		return
	}

	task, exists := state.Tasks[state.ActiveTask.TaskID]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: Active task with ID %s not found in state\n", state.ActiveTask.TaskID)
		return
	}

	// Calculate how long the task has been running
	var runningTime time.Duration
	now := time.Now()
	today := now.Format("2006-01-02")
	
	// Get time spent today (excluding current session)
	todayTimeSpent := time.Duration(0)
	if task.DailyTimeSpent != nil {
		todayTimeSpent = task.DailyTimeSpent[today]
	}
	
	if state.ActiveTask.LastPaused.IsZero() {
		// Not paused
		runningTime = now.Sub(state.ActiveTask.StartTime) - state.ActiveTask.TotalPaused
		fmt.Printf("Active task: %s\n", task.Todo)
		fmt.Printf("Running for: %s\n", formatDuration(runningTime))
		fmt.Printf("Time today: %s\n", formatDuration(todayTimeSpent+runningTime)) // Add current session
		fmt.Printf("Total time: %s\n", formatDuration(task.UsedTime+runningTime)) // Add current session to total used time
		fmt.Printf("Started at: %s\n", state.ActiveTask.StartTime.Format("2006-01-02 15:04:05"))
	} else {
		// Paused
		runningTime = state.ActiveTask.LastPaused.Sub(state.ActiveTask.StartTime) - state.ActiveTask.TotalPaused
		pausedTime := now.Sub(state.ActiveTask.LastPaused)
		fmt.Printf("Paused task: %s\n", task.Todo)
		fmt.Printf("Active time: %s\n", formatDuration(runningTime))
		fmt.Printf("Time today: %s\n", formatDuration(todayTimeSpent+runningTime)) // Add current session
		fmt.Printf("Total time: %s\n", formatDuration(task.UsedTime+runningTime)) // Add current session to total used time
		fmt.Printf("Paused for: %s\n", formatDuration(pausedTime))
		fmt.Printf("Started at: %s\n", state.ActiveTask.StartTime.Format("2006-01-02 15:04:05"))
	}
}

// BudgetCommand handles the budget command and its subcommands
func BudgetCommand(ctx *AppContext, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Budget command requires a subcommand (set, show, list)\n")
		return
	}

	subcommand := args[0]
	switch subcommand {
	case "set":
		SetBudget(ctx, args[1:])
	case "show":
		ShowBudget(ctx, args[1:])
	case "list":
		ListBudgets(ctx, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown budget subcommand: %s\n", subcommand)
	}
}

// SetBudget sets a time budget for a project
func SetBudget(ctx *AppContext, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("budget set", flag.ContinueOnError)
	descriptionFlag := fs.String("description", "", "Description of the budget")
	fs.Parse(args)

	// Get remaining args
	remainingArgs := fs.Args()
	if len(remainingArgs) < 2 {
		fmt.Fprintf(os.Stderr, "Error: 'budget set' requires project name and budget duration arguments\n")
		return
	}

	projectName := remainingArgs[0]
	budgetStr := remainingArgs[1]

	// Parse the budget duration
	budget, err := time.ParseDuration(budgetStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing budget duration: %s\n", err)
		fmt.Fprintf(os.Stderr, "Use format like '30h', '2d' (2 days = 48h), '1w' (1 week = 168h)\n")
		return
	}

	// Find the project
	state := ctx.Repository.GetAppState()
	project, exists := state.Projects[projectName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: Project '%s' not found\n", projectName)
		return
	}

	// Create and apply the event
	event := core.NewSetTimeBudgetForProject(project.ID, budget, project.UsedTime, *descriptionFlag)
	err = ctx.Repository.SaveEvent(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting budget: %s\n", err)
		return
	}

	fmt.Printf("Set budget of %s for project '%s'\n", formatDuration(budget), projectName)
}

// ShowBudget shows the budget for a project
func ShowBudget(ctx *AppContext, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("budget show", flag.ContinueOnError)
	fs.Parse(args)

	// Get remaining args
	remainingArgs := fs.Args()
	if len(remainingArgs) < 1 {
		fmt.Fprintf(os.Stderr, "Error: 'budget show' requires a project name\n")
		return
	}

	projectName := remainingArgs[0]

	// Find the project
	state := ctx.Repository.GetAppState()
	project, exists := state.Projects[projectName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: Project '%s' not found\n", projectName)
		return
	}

	// Display the budget information
	fmt.Printf("Project: %s\n", projectName)
	fmt.Printf("Budget: %s\n", formatDuration(project.Budget))
	fmt.Printf("Used time: %s (%.1f%%)\n", formatDuration(project.UsedTime), float64(project.UsedTime)/float64(project.Budget)*100)
	fmt.Printf("Remaining: %s\n", formatDuration(project.Budget-project.UsedTime))
	if project.Description != "" {
		fmt.Printf("Description: %s\n", project.Description)
	}
}

// ListBudgets lists all projects with budgets
func ListBudgets(ctx *AppContext, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("budget list", flag.ContinueOnError)
	fs.Parse(args)

	// Get the projects with budgets
	state := ctx.Repository.GetAppState()
	hasBudgets := false

	fmt.Println("Project budgets:")
	fmt.Println("---------------|----------|----------|----------|------------")
	fmt.Println("Project        | Budget   | Used     | Remaining| % Used    ")
	fmt.Println("---------------|----------|----------|----------|------------")

	for name, project := range state.Projects {
		if project.Budget > 0 {
			hasBudgets = true
			remaining := project.Budget - project.UsedTime
			percentUsed := float64(project.UsedTime) / float64(project.Budget) * 100
			fmt.Printf("%-15s| %-10s| %-10s| %-10s| %6.1f%%\n",
				name,
				formatDuration(project.Budget),
				formatDuration(project.UsedTime),
				formatDuration(remaining),
				percentUsed,
			)
		}
	}

	if !hasBudgets {
		fmt.Println("No projects with budgets found")
	}
}

// ReportCommand generates time tracking reports
func ReportCommand(ctx *AppContext, args []string) {
	// Setup flags
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	dailyFlag := fs.Bool("daily", false, "Show daily report")
	weeklyFlag := fs.Bool("weekly", false, "Show weekly report")
	monthlyFlag := fs.Bool("monthly", false, "Show monthly report")
	fromFlag := fs.String("from", "", "Start date (YYYY-MM-DD)")
	toFlag := fs.String("to", "", "End date (YYYY-MM-DD)")
	projectFlag := fs.String("project", "", "Filter by project")
	fs.Parse(args)

	// Parse date ranges
	var fromDate, toDate time.Time
	var err error
	
	if *fromFlag != "" {
		fromDate, err = time.Parse("2006-01-02", *fromFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing from date: %s\n", err)
			return
		}
	} else {
		// Default to 7 days ago
		fromDate = time.Now().AddDate(0, 0, -7)
	}
	
	if *toFlag != "" {
		toDate, err = time.Parse("2006-01-02", *toFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing to date: %s\n", err)
			return
		}
		// Set to end of day
		toDate = toDate.Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second)
	} else {
		// Default to now
		toDate = time.Now()
	}

	// Get the time records
	state := ctx.Repository.GetAppState()
	var filteredRecords []core.TimeRecord

	// Filter records by date range and project
	for _, record := range state.TimeRecords {
		if record.StartTime.After(fromDate) && record.EndTime.Before(toDate) {
			if *projectFlag == "" {
				filteredRecords = append(filteredRecords, record)
			} else {
				// Check if the task has the specified project
				task, exists := state.Tasks[record.TaskID]
				if exists {
					for _, proj := range task.Projects {
						if proj == *projectFlag {
							filteredRecords = append(filteredRecords, record)
							break
						}
					}
				}
			}
		}
	}

	// Determine report type and generate it
	if *dailyFlag {
		generateDailyReport(ctx, filteredRecords, fromDate, toDate)
	} else if *weeklyFlag {
		generateWeeklyReport(ctx, filteredRecords, fromDate, toDate)
	} else if *monthlyFlag {
		generateMonthlyReport(ctx, filteredRecords, fromDate, toDate)
	} else {
		// Default to simple list of records
		generateSimpleReport(ctx, filteredRecords, fromDate, toDate, *projectFlag)
	}
}

func generateSimpleReport(ctx *AppContext, records []core.TimeRecord, fromDate, toDate time.Time, projectFilter string) {
	if len(records) == 0 {
		fmt.Println("No time records found in the specified date range")
		return
	}

	state := ctx.Repository.GetAppState()
	
	// Print header
	fmt.Println("Time tracking report")
	fmt.Printf("Period: %s to %s\n", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	if projectFilter != "" {
		fmt.Printf("Project: %s\n", projectFilter)
	}
	fmt.Println()
	
	fmt.Println("Date       | Start    | End      | Duration | Task")
	fmt.Println("-----------|----------|----------|----------|--------------------")
	
	// Track total duration
	var totalDuration time.Duration
	
	// Print each record
	for _, record := range records {
		task, exists := state.Tasks[record.TaskID]
		taskName := "Unknown task"
		if exists {
			taskName = task.Todo
			if len(taskName) > 40 {
				taskName = taskName[:37] + "..."
			}
		}
		
		totalDuration += record.Duration
		
		fmt.Printf("%s | %s | %s | %-8s | %s\n",
			record.StartTime.Format("2006-01-02"),
			record.StartTime.Format("15:04:05"),
			record.EndTime.Format("15:04:05"),
			formatDuration(record.Duration),
			taskName,
		)
	}
	
	fmt.Println("-----------|----------|----------|----------|--------------------")
	fmt.Printf("Total      |          |          | %-8s |\n", formatDuration(totalDuration))
}

func generateDailyReport(ctx *AppContext, records []core.TimeRecord, fromDate, toDate time.Time) {
	if len(records) == 0 {
		fmt.Println("No time records found in the specified date range")
		return
	}

	state := ctx.Repository.GetAppState()
	
	// Group records by day
	dailyTimes := make(map[string]time.Duration)
	dailyTasks := make(map[string]map[uuid.UUID]time.Duration)
	
	for _, record := range records {
		day := record.StartTime.Format("2006-01-02")
		dailyTimes[day] += record.Duration
		
		if dailyTasks[day] == nil {
			dailyTasks[day] = make(map[uuid.UUID]time.Duration)
		}
		dailyTasks[day][record.TaskID] += record.Duration
	}
	
	// Print header
	fmt.Println("Daily time tracking report")
	fmt.Printf("Period: %s to %s\n", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	fmt.Println()
	
	fmt.Println("Date       | Total    | Top Tasks")
	fmt.Println("-----------|----------|--------------------")
	
	// Sort days
	var days []string
	for day := range dailyTimes {
		days = append(days, day)
	}
	sort.Strings(days)
	
	// Print each day's summary
	totalDuration := time.Duration(0)
	for _, day := range days {
		duration := dailyTimes[day]
		totalDuration += duration
		
		// Find top 3 tasks for this day
		type taskTime struct {
			id       uuid.UUID
			duration time.Duration
		}
		var taskTimes []taskTime
		
		for taskID, taskDuration := range dailyTasks[day] {
			taskTimes = append(taskTimes, taskTime{taskID, taskDuration})
		}
		
		sort.Slice(taskTimes, func(i, j int) bool {
			return taskTimes[i].duration > taskTimes[j].duration
		})
		
		// Format top tasks
		var topTasks string
		for i, tt := range taskTimes {
			if i >= 3 {
				break
			}
			
			task, exists := state.Tasks[tt.id]
			taskName := "Unknown task"
			if exists {
				taskName = task.Todo
				if len(taskName) > 30 {
					taskName = taskName[:27] + "..."
				}
			}
			
			if i > 0 {
				topTasks += ", "
			}
			topTasks += fmt.Sprintf("%s (%s)", taskName, formatDuration(tt.duration))
		}
		
		fmt.Printf("%s | %-8s | %s\n", day, formatDuration(duration), topTasks)
	}
	
	fmt.Println("-----------|----------|--------------------")
	fmt.Printf("Total      | %-8s |\n", formatDuration(totalDuration))
}

func generateWeeklyReport(ctx *AppContext, records []core.TimeRecord, fromDate, toDate time.Time) {
	if len(records) == 0 {
		fmt.Println("No time records found in the specified date range")
		return
	}

	state := ctx.Repository.GetAppState()
	
	// Group records by week
	weeklyTimes := make(map[string]time.Duration)
	weeklyTasks := make(map[string]map[uuid.UUID]time.Duration)
	
	for _, record := range records {
		// Get the start of the week (Monday)
		year, week := record.StartTime.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		
		weeklyTimes[weekKey] += record.Duration
		
		if weeklyTasks[weekKey] == nil {
			weeklyTasks[weekKey] = make(map[uuid.UUID]time.Duration)
		}
		weeklyTasks[weekKey][record.TaskID] += record.Duration
	}
	
	// Print header
	fmt.Println("Weekly time tracking report")
	fmt.Printf("Period: %s to %s\n", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	fmt.Println()
	
	fmt.Println("Week       | Total    | Top Tasks")
	fmt.Println("-----------|----------|--------------------")
	
	// Sort weeks
	var weeks []string
	for week := range weeklyTimes {
		weeks = append(weeks, week)
	}
	sort.Strings(weeks)
	
	// Print each week's summary
	totalDuration := time.Duration(0)
	for _, week := range weeks {
		duration := weeklyTimes[week]
		totalDuration += duration
		
		// Find top 3 tasks for this week
		type taskTime struct {
			id       uuid.UUID
			duration time.Duration
		}
		var taskTimes []taskTime
		
		for taskID, taskDuration := range weeklyTasks[week] {
			taskTimes = append(taskTimes, taskTime{taskID, taskDuration})
		}
		
		sort.Slice(taskTimes, func(i, j int) bool {
			return taskTimes[i].duration > taskTimes[j].duration
		})
		
		// Format top tasks
		var topTasks string
		for i, tt := range taskTimes {
			if i >= 3 {
				break
			}
			
			task, exists := state.Tasks[tt.id]
			taskName := "Unknown task"
			if exists {
				taskName = task.Todo
				if len(taskName) > 30 {
					taskName = taskName[:27] + "..."
				}
			}
			
			if i > 0 {
				topTasks += ", "
			}
			topTasks += fmt.Sprintf("%s (%s)", taskName, formatDuration(tt.duration))
		}
		
		fmt.Printf("%s | %-8s | %s\n", week, formatDuration(duration), topTasks)
	}
	
	fmt.Println("-----------|----------|--------------------")
	fmt.Printf("Total      | %-8s |\n", formatDuration(totalDuration))
}

func generateMonthlyReport(ctx *AppContext, records []core.TimeRecord, fromDate, toDate time.Time) {
	if len(records) == 0 {
		fmt.Println("No time records found in the specified date range")
		return
	}

	state := ctx.Repository.GetAppState()
	
	// Group records by month
	monthlyTimes := make(map[string]time.Duration)
	monthlyTasks := make(map[string]map[uuid.UUID]time.Duration)
	
	for _, record := range records {
		month := record.StartTime.Format("2006-01")
		monthlyTimes[month] += record.Duration
		
		if monthlyTasks[month] == nil {
			monthlyTasks[month] = make(map[uuid.UUID]time.Duration)
		}
		monthlyTasks[month][record.TaskID] += record.Duration
	}
	
	// Print header
	fmt.Println("Monthly time tracking report")
	fmt.Printf("Period: %s to %s\n", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	fmt.Println()
	
	fmt.Println("Month      | Total    | Top Tasks")
	fmt.Println("-----------|----------|--------------------")
	
	// Sort months
	var months []string
	for month := range monthlyTimes {
		months = append(months, month)
	}
	sort.Strings(months)
	
	// Print each month's summary
	totalDuration := time.Duration(0)
	for _, month := range months {
		duration := monthlyTimes[month]
		totalDuration += duration
		
		// Find top 3 tasks for this month
		type taskTime struct {
			id       uuid.UUID
			duration time.Duration
		}
		var taskTimes []taskTime
		
		for taskID, taskDuration := range monthlyTasks[month] {
			taskTimes = append(taskTimes, taskTime{taskID, taskDuration})
		}
		
		sort.Slice(taskTimes, func(i, j int) bool {
			return taskTimes[i].duration > taskTimes[j].duration
		})
		
		// Format top tasks
		var topTasks string
		for i, tt := range taskTimes {
			if i >= 3 {
				break
			}
			
			task, exists := state.Tasks[tt.id]
			taskName := "Unknown task"
			if exists {
				taskName = task.Todo
				if len(taskName) > 30 {
					taskName = taskName[:27] + "..."
				}
			}
			
			if i > 0 {
				topTasks += ", "
			}
			topTasks += fmt.Sprintf("%s (%s)", taskName, formatDuration(tt.duration))
		}
		
		fmt.Printf("%s | %-8s | %s\n", month, formatDuration(duration), topTasks)
	}
	
	fmt.Println("-----------|----------|--------------------")
	fmt.Printf("Total      | %-8s |\n", formatDuration(totalDuration))
}

// Helper function to format duration in a readable way
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh%02dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm%02ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}