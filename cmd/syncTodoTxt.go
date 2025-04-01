package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"t5.mkbrechtel.dev/core"
	todo "t5.mkbrechtel.dev/sync/todotxt"
)

var todoTxtFilePath string
var syncAllFiles bool

// syncTodoTxtCmd represents the command to sync with a todo.txt file
var syncTodoTxtCmd = &cobra.Command{
	Use:   "todotxt",
	Short: "Sync tasks with a todo.txt file",
	Long: `Sync tasks with a todo.txt formatted file.
	
This command will synchronize tasks between the t5 event store and todo.txt files.
Changes from either side will be merged, with conflict detection.

You can sync a specific file with --file or sync all configured files with --all.
	
Examples:
  t5 sync todotxt --file ~/todo.txt
  t5 sync todotxt --all
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create or open the event store
		// TODO: This should be centralized and reused across commands
		eventStorePath := viper.GetString("eventstore.file")
		if eventStorePath == "" {
			return fmt.Errorf("no event store file configured, use --eventstore flag or configure eventstore.file in config")
		}
		eventStore, err := core.NewAppendLogEventStore(eventStorePath)
		if err != nil {
			return fmt.Errorf("failed to open event store: %w", err)
		}
		
		// Create repository with the event store
		repo := core.NewRepositoryWithEventStore(eventStore)
		
		// Rebuild state from events
		if err := repo.RebuildState(); err != nil {
			return fmt.Errorf("failed to rebuild state: %w", err)
		}
		
		// Determine which files to sync
		if syncAllFiles {
			// Get all configured todo.txt files from config
			todoFiles, err := GetTodoFiles()
			if err != nil {
				return fmt.Errorf("failed to get todo files from config: %w", err)
			}
			
			if len(todoFiles) == 0 {
				return fmt.Errorf("no todo.txt files configured in config")
			}
			
			for _, fileConfig := range todoFiles {
				fmt.Printf("Syncing todo.txt file: %s\n", fileConfig.Path)
				result, err := todo.SyncWithRepository(repo, fileConfig.Path)
				if err != nil {
					return fmt.Errorf("sync failed for %s: %w", fileConfig.Path, err)
				}
				
				// Print results for this file
				fmt.Printf("Sync completed for %s:\n", fileConfig.Path)
				fmt.Printf("  Added: %d\n", result.Added)
				fmt.Printf("  Updated from todo.txt: %d\n", result.FromTodoTxt)
				fmt.Printf("  Updated to todo.txt: %d\n", result.ToTodoTxt)
				fmt.Printf("  Skipped: %d\n", result.Skipped)
				if result.Conflicts > 0 {
					fmt.Printf("  Conflicts detected: %d\n", result.Conflicts)
				}
				fmt.Println()
			}
		} else {
			// Sync with the specified todo.txt file
			result, err := todo.SyncWithRepository(repo, todoTxtFilePath)
			if err != nil {
				return fmt.Errorf("sync failed: %w", err)
			}
			
			// Print results
			fmt.Printf("Sync completed for %s:\n", todoTxtFilePath)
			fmt.Printf("  Added: %d\n", result.Added)
			fmt.Printf("  Updated from todo.txt: %d\n", result.FromTodoTxt)
			fmt.Printf("  Updated to todo.txt: %d\n", result.ToTodoTxt)
			fmt.Printf("  Skipped: %d\n", result.Skipped)
			if result.Conflicts > 0 {
				fmt.Printf("  Conflicts detected: %d\n", result.Conflicts)
			}
		}
		
		return nil
	},
}

func init() {
	syncCmd.AddCommand(syncTodoTxtCmd)
	
	// Add flags
	syncTodoTxtCmd.Flags().StringVarP(&todoTxtFilePath, "file", "f", "todo.txt", "Path to the todo.txt file")
	syncTodoTxtCmd.Flags().BoolVarP(&syncAllFiles, "all", "a", false, "Sync all configured todo.txt files")
}