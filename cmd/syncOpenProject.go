package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "t/sync/openproject"
    "t/todo"
    todotxt "github.com/1set/todotxt"
)

var syncOpenProjectCmd = &cobra.Command{
    Use:   "openproject",
    Short: "Sync with OpenProject",
    Long: `t sync openproject

    Syncs tasks from OpenProject with your local todo.txt file.
    New tasks will be added and existing tasks will be updated if needed.
    Tasks are matched using their OpenProject URL.
    `,
    Run: func(cmd *cobra.Command, args []string) {
        // Validate input configuration
        url := viper.GetString("sync.openproject.url")
        if url == "" {
            fmt.Println("Error: OpenProject URL is not configured")
            os.Exit(1)
        }

        apiKey := viper.GetString("sync.openproject.api-key")
        if apiKey == "" {
            fmt.Println("Error: OpenProject API key is not configured")
            os.Exit(1)
        }

        queryId := viper.GetString("sync.openproject.query-id")
        if queryId == "" {
            fmt.Println("Error: OpenProject query ID is not configured")
            os.Exit(1)
        }

        prefix := viper.GetString("sync.openproject.todo-prefix")

        // Fetch work packages from OpenProject
        fmt.Println("Fetching tasks from OpenProject...")
        workPackages, err := openproject.GetWorkPackages(url, apiKey, queryId)
        if err != nil {
            fmt.Printf("Error fetching work packages: %v\n", err)
            os.Exit(1)
        }

        // Convert work packages to todo list
        fmt.Println("Converting work packages to tasks...")
        sourceList := openproject.CreateTaskList(workPackages, prefix, url)

        // Load existing todo file
        fmt.Printf("Loading %s...\n", todoFile)
        targetList, err := todotxt.LoadFromPath(todoFile)
        if err != nil {
            fmt.Printf("Error loading todo file: %v\n", err)
            os.Exit(1)
        }

        // Ensure all tasks have proper metadata
        targetList = todo.EnsureProperTasks(targetList)

        // Perform sync
        fmt.Println("Syncing tasks...")
        updatedList, result, err := todo.SyncTaskLists(targetList, sourceList)
        if err != nil {
            fmt.Printf("Error during sync: %v\n", err)
            os.Exit(1)
        }

        // Save the updated list
        fmt.Printf("Saving changes to %s...\n", todoFile)
        if err := updatedList.WriteToPath(todoFile); err != nil {
            fmt.Printf("Error saving todo file: %v\n", err)
            os.Exit(1)
        }

        // Print results
        fmt.Printf("\nSync completed successfully:\n")
        fmt.Printf("  Added: %d tasks\n", result.Added)
        fmt.Printf("  Updated: %d tasks\n", result.Updated)
        fmt.Printf("  Skipped: %d tasks (no changes needed)\n", result.Skipped)
    },
}
