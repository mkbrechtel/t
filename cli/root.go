package cli

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
)

// Execute parses flags and runs the appropriate command with the given args
// If stdin is nil, os.Stdin will be used
// If stdout and stderr are nil, os.Stdout and os.Stderr will be used
func Execute(args []string, stdin io.ReadCloser, stdout, stderr io.WriteCloser) error {

	ctx := NewAppContext()

	// Setup custom flag set for parsing the given args
	fs := flag.NewFlagSet("t5", flag.ExitOnError)
	ctx.FlagSet = fs

	// Setup flags on our custom FlagSet
	ctx.SetupFlags()

	// Set custom usage function
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: t5 [flags] <command> [arguments]\n\n")
		fmt.Fprintf(stderr, "t5 is a todo list and time tracker tool\n\n")
		fmt.Fprintf(stderr, "Commands:\n")
		fmt.Fprintf(stderr, "  list todo            List tasks from the event store\n")
		fmt.Fprintf(stderr, "  add todo [task text] Add a new task (reads from stdin if no text provided)\n")
		fmt.Fprintf(stderr, "  modify <task-id>     Modify an existing task with various flags\n")
		fmt.Fprintf(stderr, "  update todo [file]   Update and ensure properties of tasks in your todo list\n")
		fmt.Fprintf(stderr, "  sync todo [file]     Sync tasks with a todo.txt file\n")
		fmt.Fprintf(stderr, "  config               Show the current configuration\n")
		fmt.Fprintf(stderr, "\nFilter flags for list command:\n")
		fmt.Fprintf(stderr, "  --completed         Show only completed tasks\n")
		fmt.Fprintf(stderr, "  --not-completed     Show only non-completed tasks\n")
		fmt.Fprintf(stderr, "  --priority=A        Filter by priority A (can use A,B,C for multiple)\n")
		fmt.Fprintf(stderr, "  --project=Home      Filter tasks with +Home project\n")
		fmt.Fprintf(stderr, "  --context=@phone    Filter tasks with @phone context\n")
		fmt.Fprintf(stderr, "  --due-today         Filter tasks due today\n")
		fmt.Fprintf(stderr, "  --due-before=DATE   Filter tasks due before DATE (YYYY-MM-DD)\n")
		fmt.Fprintf(stderr, "  --created-after=DATE Filter tasks created after DATE (YYYY-MM-DD)\n")
		fmt.Fprintf(stderr, "  --has-tag=tag_name  Filter tasks with specific tag\n")
		fmt.Fprintf(stderr, "  --regex=PATTERN     Filter tasks using a regular expression\n")
		fmt.Fprintf(stderr, "\nModifier flags for modify command:\n")
		fmt.Fprintf(stderr, "  --complete          Mark task as complete\n")
		fmt.Fprintf(stderr, "  --uncomplete        Mark task as incomplete\n")
		fmt.Fprintf(stderr, "  --priority=A        Set priority to A\n")
		fmt.Fprintf(stderr, "  --remove-priority   Remove priority\n")
		fmt.Fprintf(stderr, "  --add-project=Home  Add +Home project to task\n")
		fmt.Fprintf(stderr, "  --remove-project=Home Remove +Home project from task\n")
		fmt.Fprintf(stderr, "  --add-context=phone Add @phone context to task\n")
		fmt.Fprintf(stderr, "  --remove-context=phone Remove @phone context from task\n")
		fmt.Fprintf(stderr, "  --due=DATE          Set due date (YYYY-MM-DD)\n")
		fmt.Fprintf(stderr, "  --remove-due        Remove due date\n")
		fmt.Fprintf(stderr, "  --append=TEXT       Append text to task\n")
		fmt.Fprintf(stderr, "  --prepend=TEXT      Prepend text to task\n")
		fmt.Fprintf(stderr, "\nGlobal flags:\n")
		fs.PrintDefaults()
	}

	// Parse flags - handle empty args case
	if len(args) > 1 {
		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("error parsing flags: %w", err)
		}
	}

	// Load config file after parsing flags
	if err := ctx.LoadConfigFile(); err != nil {
		log.Printf("Warning: failed to load config file: %v", err)
	}

	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		fs.Usage()
		return fmt.Errorf("no command specified")
	}

	command := cmdArgs[0]

	// Initialize repository early so commands can use it
	if err := ctx.InitRepository(); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	switch command {
	case "list":
		if len(cmdArgs) > 1 && cmdArgs[1] == "todo" {
			ListTasks(ctx)
		} else {
			fmt.Fprintf(stderr, "Error: Command 'list' requires 'todo' as the object\n")
			return fmt.Errorf("command 'list' requires 'todo' as the object")
		}
	case "add":
		if len(cmdArgs) > 1 && cmdArgs[1] == "todo" {
			if len(cmdArgs) > 2 {
				// The rest of the arguments form the task text
				taskText := strings.Join(cmdArgs[2:], " ")
				AddTask(ctx, taskText, stdin)
			} else {
				// Read from stdin if no arguments provided
				AddTask(ctx, "", stdin)
			}
		} else {
			fmt.Fprintf(stderr, "Error: Command 'add' requires 'todo' as the object\n")
			return fmt.Errorf("command 'add' requires 'todo' as the object")
		}
	case "modify":
		if len(cmdArgs) > 1 {
			// First argument is the task ID
			taskID := cmdArgs[1]
			ModifyTask(ctx, taskID)
		} else {
			fmt.Fprintf(stderr, "Error: Task ID required for modify command\n")
			return fmt.Errorf("task ID required for modify command")
		}
	case "update":
		if len(cmdArgs) > 1 && cmdArgs[1] == "todo" {
			// Check if a specific file was specified
			if len(cmdArgs) > 2 {
				ctx.Config.TodoFile = cmdArgs[2]
			}
			UpdateTasks(ctx)
		} else {
			fmt.Fprintf(stderr, "Error: Command 'update' requires 'todo' as the object\n")
			return fmt.Errorf("command 'update' requires 'todo' as the object")
		}
	case "todo":
		// If there's a subcommand
		if len(cmdArgs) > 1 {
			// Process the todo subcommand
			subCommand := cmdArgs[1]
			switch subCommand {
			case "add":
				if len(cmdArgs) > 2 {
					// The rest of the arguments form the task text
					taskText := strings.Join(cmdArgs[2:], " ")
					AddTask(ctx, taskText, stdin)
				} else {
					// Read from stdin if no arguments provided
					AddTask(ctx, "", stdin)
				}
			case "modify":
				if len(cmdArgs) > 2 {
					// Second argument is the task ID
					taskID := cmdArgs[2]
					ModifyTask(ctx, taskID)
				} else {
					fmt.Fprintf(stderr, "Error: Task ID required for modify command\n")
					return fmt.Errorf("task ID required for modify command")
				}
			case "update":
				// Check if a specific file was specified
				if len(cmdArgs) > 2 {
					ctx.Config.TodoFile = cmdArgs[2]
				}
				UpdateTasks(ctx)
			default:
				fmt.Fprintf(stderr, "Unknown todo subcommand: %s\n", subCommand)
				fs.Usage()
				return fmt.Errorf("unknown todo subcommand: %s", subCommand)
			}
		} else {
			// Default behavior for 'todo' with no subcommand is to update
			UpdateTasks(ctx)
		}
	case "sync":
		if len(cmdArgs) > 1 && cmdArgs[1] == "todo" {
			// Check if a specific file was specified
			if len(cmdArgs) > 2 {
				ctx.Config.TodoFile = cmdArgs[2]
			}
			SyncTasks(ctx)
		} else {
			fmt.Fprintf(stderr, "Error: Command 'sync' requires 'todo' as the object\n")
			return fmt.Errorf("command 'sync' requires 'todo' as the object")
		}
	case "config":
		ctx.ShowConfig()
	default:
		fmt.Fprintf(stderr, "Unknown command: %s\n", command)
		fs.Usage()
		return fmt.Errorf("unknown command: %s", command)
	}

	// Close the write ends of the pipes to flush output
	stdout.Close()
	stderr.Close()

	return nil
}

// GetAppContextForTesting returns a new app context for testing
func GetAppContextForTesting() *AppContext {
	return NewAppContext()
}
