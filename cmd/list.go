package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks from the event store",
	Long: `t list

Lists tasks from the event store. This command demonstrates the use of
the event replay system to reconstruct the current state.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize repository with configured event store
		repo, err := InitRepository()
		if err != nil {
			log.Fatalf("Failed to initialize repository: %v", err)
		}

		// List all tasks
		tasks, err := repo.ListTasks(nil)
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
			fmt.Printf("%d. [%s] (%s) %s\n", i+1, status, priority, task.Todo)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}