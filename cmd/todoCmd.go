package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"t5.mkbrechtel.dev/core"
	todo "t5.mkbrechtel.dev/sync/todotxt"
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage your todo list",
	Long: `t todo
	
	With this command you can manage your todo list. 
	
	The todo list is based on the todo.txt format. See http://todotxt.org
	`,
}

// todoUpdateCmd represents the todo update command
var todoUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update and ensure properties of tasks in your todo list",
	Long: `t todo update

	With this command you can update your todo list tasks to ensure they have all required properties:
	- Completion dates for completed tasks
	- Unique identifiers (short or long form)
	- Default tags

	It processes the entire todo.txt file and updates task properties according to configuration.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read tasks from todo.txt file
		taskList, err := todo.ReadTodoFile(todoFile)
		if err != nil {
			log.Fatalf("Failed to read todo file: %v", err)
		}

		config := todo.DefaultEnsureConfig
		// Set default tags that should be applied to all tasks
		config.DefaultTags = map[string]string{
			// "app":     "t",
			// "version": "1.0",
		}

		taskList = todo.EnsureTaskListProperties(taskList, config)

		// Write updates back to todo.txt file
		if err = todo.WriteTodoFile(taskList, todoFile); err != nil {
			log.Fatalf("Failed to write todo file: %v", err)
		}

		// Initialize repository with configured event store
		repo, err := InitRepository()
		if err != nil {
			log.Fatalf("Failed to initialize repository: %v", err)
		}

		// Create TodoTxtTaskUpdate event with the content of the todo.txt file
		content, err := todo.GetTodoFileContent(taskList)
		if err != nil {
			log.Fatalf("Failed to get todo file content: %v", err)
		}

		event := core.NewTodoTxtTaskUpdate(content, todoFile)
		err = repo.SaveEvent(event)
		if err != nil {
			log.Fatalf("Failed to save event: %v", err)
		}

		log.Println("Task update saved to event store")
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)

	// Add todo subcommands
	todoCmd.AddCommand(todoUpdateCmd)

	// Add configuration flags
	todoCmd.PersistentFlags().Bool("prefer-short-ids", true, "Use short form IDs instead of UUIDs")
    todoCmd.PersistentFlags().Bool("enforce-creation-date", true, "Ensure tasks have creation dates")
	todoCmd.PersistentFlags().Bool("enforce-completion-date", true, "Ensure completed tasks have completion dates")
	
	// Bind flags to viper configuration
	viper.BindPFlag("todo.ensure.preferShortIds", todoCmd.PersistentFlags().Lookup("prefer-short-ids"))
    viper.BindPFlag("todo.ensure.enforceCreationDate", todoCmd.PersistentFlags().Lookup("enforce-creation-date"))
	viper.BindPFlag("todo.ensure.enforceCompletionDate", todoCmd.PersistentFlags().Lookup("enforce-completion-date"))
}
