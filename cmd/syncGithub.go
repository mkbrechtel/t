package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "t/sync/github"
    "t/todo"
    todotxt "github.com/1set/todotxt"
)

var syncGithubCmd = &cobra.Command{
    Use:   "github",
    Short: "Sync with GitHub",
    Long: `t sync github

    Syncs tasks from GitHub with your local todo.txt file.
    New tasks will be added and existing tasks will be updated if needed.
    Tasks are matched using their GitHub issue URL.
    This syncs issues from all repositories in the user's account.
    `,
    Run: func(cmd *cobra.Command, args []string) {
        // Validate input configuration
        token := viper.GetString("sync.github.token")
        if token == "" {
            fmt.Println("Error: GitHub access token is not configured")
            os.Exit(1)
        }

        // Get GitHub API base URL
        baseURL := viper.GetString("sync.github.api_base_url")
        if baseURL == "" {
            fmt.Println("Error: GitHub API base URL is not configured")
            os.Exit(1)
        }

        // Get GitHub API endpoint
        endpoint := viper.GetString("sync.github.api_endpoint")
        if endpoint == "" {
            fmt.Println("Error: GitHub API endpoint is not configured")
            os.Exit(1)
        }

        issuePrefix := viper.GetString("sync.github.issue_prefix")
        pullPrefix := viper.GetString("sync.github.pull_prefix")
        
        // Fetch issues from GitHub
        fmt.Println("Fetching issues from GitHub...")
        issues, err := github.GetUserIssues(token, baseURL, endpoint)
        if err != nil {
            fmt.Printf("Error fetching GitHub issues: %v\n", err)
            os.Exit(1)
        }

        // Convert GitHub issues to todo list
        fmt.Println("Converting GitHub issues to tasks...")
        sourceList := github.CreateTaskList(issues, issuePrefix, pullPrefix)

        fmt.Print(sourceList)

        // Load existing todo file
        fmt.Printf("Loading %s...\n", todoFile)
        targetList, err := todotxt.LoadFromPath(todoFile)
        if err != nil {
            fmt.Printf("Error loading todo file: %v\n", err)
            os.Exit(1)
        }

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

func init() {
    syncCmd.AddCommand(syncGithubCmd)

    syncGithubCmd.PersistentFlags().String("token", "", "GitHub access token")
    syncGithubCmd.PersistentFlags().String("issue-prefix", "GitHub Issue: ", "Prefix for GitHub issue todos created by t")
    syncGithubCmd.PersistentFlags().String("pull-prefix", "GitHub PR: ", "Prefix for GitHub pull request todos created by t")
    syncGithubCmd.PersistentFlags().String("api-base-url", "https://api.github.com", "GitHub API base URL")
    syncGithubCmd.PersistentFlags().String("api-endpoint", "/issues?filter=assigned&state=all&per_page=1000&pulls=1", "GitHub API endpoint")

    viper.BindPFlag("sync.github.token", syncGithubCmd.PersistentFlags().Lookup("token"))
    viper.BindPFlag("sync.github.issue_prefix", syncGithubCmd.PersistentFlags().Lookup("issue-prefix"))
    viper.BindPFlag("sync.github.pull_prefix", syncGithubCmd.PersistentFlags().Lookup("pull-prefix"))
    viper.BindPFlag("sync.github.api_base_url", syncGithubCmd.PersistentFlags().Lookup("api-base-url"))
    viper.BindPFlag("sync.github.api_endpoint", syncGithubCmd.PersistentFlags().Lookup("api-endpoint"))
}
