package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "t/sync/gitlab"
    "t/todo"
    todotxt "github.com/1set/todotxt"
)

var syncGitlabCmd = &cobra.Command{
    Use:   "gitlab",
    Short: "Sync with GitLab",
    Long: `t sync gitlab

    Syncs tasks from GitLab with your local todo.txt file.
    New tasks will be added and existing tasks will be updated if needed.
    Tasks are matched using their GitLab issue URL.
    This syncs issues from all projects the user has access to.
    `,
    Run: func(cmd *cobra.Command, args []string) {
        // Validate input configuration
        token := viper.GetString("sync.gitlab.token")
        if token == "" {
            fmt.Println("Error: GitLab access token is not configured")
            os.Exit(1)
        }

        // Get GitLab API base URL
        baseURL := viper.GetString("sync.gitlab.api_base_url")
        if baseURL == "" {
            fmt.Println("Error: GitLab API base URL is not configured")
            os.Exit(1)
        }

        // Get GitHub API endpoint
        endpoint := viper.GetString("sync.gitlab.api_endpoint")
        if endpoint == "" {
            fmt.Println("Error: GitLab API endpoint is not configured")
            os.Exit(1)
        }

        issuePrefix := viper.GetString("sync.gitlab.issue_prefix")
        // mergeRequestPrefix := viper.GetString("sync.gitlab.merge_request_prefix")

        // Fetch issues from GitLab
        fmt.Println("Fetching issues from GitLab...")
        issues, err := gitlab.GetUserIssues(token, baseURL, endpoint)
        if err != nil {
            fmt.Printf("Error fetching GitLab issues: %v\n", err)
            os.Exit(1)
        }

        fmt.Print(issues)


        // Convert GitLab issues to todo list
        fmt.Println("Converting GitLab issues to tasks...")
        sourceList := gitlab.CreateIssueTaskList(issues, issuePrefix)
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
    syncCmd.AddCommand(syncGitlabCmd)

    syncGitlabCmd.PersistentFlags().String("token", "", "GitLab access token")
    syncGitlabCmd.PersistentFlags().String("issue-prefix", "GitLab Issue: ", "Prefix for GitLab issue todos created by t")
    syncGitlabCmd.PersistentFlags().String("merge-request-prefix", "GitLab MR: ", "Prefix for GitLab merge request todos created by t")
    syncGitlabCmd.PersistentFlags().String("api-base-url", "https://gitlab.com/api/v4", "GitLab API base URL")
    syncGitlabCmd.PersistentFlags().String("api-endpoint", "/issues", "GitLab API endpoint")

    viper.BindPFlag("sync.gitlab.token", syncGitlabCmd.PersistentFlags().Lookup("token"))
    viper.BindPFlag("sync.gitlab.issue_prefix", syncGitlabCmd.PersistentFlags().Lookup("issue-prefix"))
    viper.BindPFlag("sync.gitlab.merge_request_prefix", syncGitlabCmd.PersistentFlags().Lookup("merge-request-prefix"))
    viper.BindPFlag("sync.gitlab.api_base_url", syncGitlabCmd.PersistentFlags().Lookup("api-base-url"))
    viper.BindPFlag("sync.gitlab.api_endpoint", syncGitlabCmd.PersistentFlags().Lookup("api-endpoint"))
}
