package cmd

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// Execute parses flags and runs the appropriate command with the given args
func Execute(args []string) error {
	ctx := NewAppContext()
	
	// Setup custom flag set for parsing the given args
	fs := flag.NewFlagSet("t5", flag.ExitOnError)
	ctx.FlagSet = fs
	
	// Setup flags on our custom FlagSet
	ctx.SetupFlags()

	// Set custom usage function
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: t5 [flags] <command> [arguments]\n\n")
		fmt.Fprintf(os.Stderr, "t5 is a todo list and time tracker tool\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  list                List tasks from the event store\n")
		fmt.Fprintf(os.Stderr, "  update [file]       Update and ensure properties of tasks in your todo list\n")
		fmt.Fprintf(os.Stderr, "  sync [file]         Sync tasks with a todo.txt file\n")
		fmt.Fprintf(os.Stderr, "  config              Show the current configuration\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
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
		ListTasks(ctx)
	case "update":
		// Check if a specific file was specified
		if len(cmdArgs) > 1 {
			ctx.Config.TodoFile = cmdArgs[1]
		}
		UpdateTasks(ctx)
	case "todo":
		// If there's a subcommand
		if len(cmdArgs) > 1 {
			// Process the todo subcommand
			subCommand := cmdArgs[1]
			switch subCommand {
			case "update":
				// Check if a specific file was specified
				if len(cmdArgs) > 2 {
					ctx.Config.TodoFile = cmdArgs[2]
				}
				UpdateTasks(ctx)
			default:
				fmt.Fprintf(os.Stderr, "Unknown todo subcommand: %s\n", subCommand)
				fs.Usage()
				return fmt.Errorf("unknown todo subcommand: %s", subCommand)
			}
		} else {
			// Default behavior for 'todo' with no subcommand is to update
			UpdateTasks(ctx)
		}
	case "sync":
		// Check if a specific file was specified
		if len(cmdArgs) > 1 {
			ctx.Config.TodoFile = cmdArgs[1]
		}
		SyncTasks(ctx)
	case "config":
		ctx.ShowConfig()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fs.Usage()
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

// ExecuteWithArgs is a convenience wrapper that allows setting both arguments and output streams
func ExecuteWithArgs(args []string, stdout, stderr io.Writer) error {
	// Save original stdout and stderr
	oldStdout, oldStderr := os.Stdout, os.Stderr
	
	// Create pipes for redirecting output
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	
	// Replace stdout and stderr
	os.Stdout = wOut
	os.Stderr = wErr
	
	// Restore original stdout and stderr when done
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()
	
	// Start copying to provided writers
	go io.Copy(stdout, rOut)
	go io.Copy(stderr, rErr)
	
	// Execute the command
	err := Execute(args)
	
	// Close the write ends of the pipes to flush output
	wOut.Close()
	wErr.Close()
	
	return err
}

// GetAppContextForTesting returns a new app context for testing
func GetAppContextForTesting() *AppContext {
	return NewAppContext()
}