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

        // Get GitLab API endpoints
        issuesEndpoint := viper.GetString("sync.gitlab.issues_endpoint")
        if issuesEndpoint == "" {
            fmt.Println("Error: GitLab API issues endpoint is not configured")
            os.Exit(1)
        }

        mergeRequestsEndpoint := viper.GetString("sync.gitlab.merge_requests_endpoint")
        if mergeRequestsEndpoint == "" {
            fmt.Println("Error: GitLab API merge requests endpoint is not configured")
            os.Exit(1)
        }

        issuePrefix := viper.GetString("sync.gitlab.issue_prefix")
        mergeRequestPrefix := viper.GetString("sync.gitlab.merge_request_prefix")

        // Fetch issues from GitLab
        fmt.Println("Fetching issues from GitLab...")
        issues, err := gitlab.GetUserIssues(token, baseURL, issuesEndpoint)
        if err != nil {
            fmt.Printf("Error fetching GitLab issues: %v\n", err)
            os.Exit(1)
        }

        // Convert GitLab issues to todo list
        fmt.Println("Converting GitLab issues to tasks...")
        issuesSourceList := gitlab.CreateIssueTaskList(issues, issuePrefix)
        
        // Fetch merge requests from GitLab
        fmt.Println("Fetching merge requests from GitLab...")
        mergeRequests, err := gitlab.GetUserMergeRequests(token, baseURL, mergeRequestsEndpoint)
        if err != nil {
            fmt.Printf("Error fetching GitLab merge requests: %v\n", err)
            os.Exit(1)
        }

        // Convert GitLab merge requests to todo list
        fmt.Println("Converting GitLab merge requests to tasks...")
        mergeRequestsSourceList := gitlab.CreateMergeRequestTaskList(mergeRequests, mergeRequestPrefix)
        
        // Load existing todo file
        fmt.Printf("Loading %s...\n", todoFile)
        targetList, err := todotxt.LoadFromPath(todoFile)
        if err != nil {
            fmt.Printf("Error loading todo file: %v\n", err)
            os.Exit(1)
        }

        // Perform sync for issues
        fmt.Println("Syncing issue tasks...")
        updatedList, issueResult, err := todo.SyncTaskLists(targetList, issuesSourceList)
        if err != nil {
            fmt.Printf("Error during issue sync: %v\n", err)
            os.Exit(1)
        }

        // Perform sync for merge requests
        fmt.Println("Syncing merge request tasks...")
        finalList, mrResult, err := todo.SyncTaskLists(updatedList, mergeRequestsSourceList)
        if err != nil {
            fmt.Printf("Error during merge request sync: %v\n", err)
            os.Exit(1)
        }

        // Combine results
        result := todo.SyncResult{
            Added:   issueResult.Added + mrResult.Added,
            Updated: issueResult.Updated + mrResult.Updated,
            Skipped: issueResult.Skipped + mrResult.Skipped,
        }

        // Save the updated list
        fmt.Printf("Saving changes to %s...\n", todoFile)
        if err := finalList.WriteToPath(todoFile); err != nil {
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
    syncGitlabCmd.PersistentFlags().String("issues-endpoint", "/issues", "GitLab API endpoint for issues")
    syncGitlabCmd.PersistentFlags().String("merge-requests-endpoint", "/merge_requests", "GitLab API endpoint for merge requests")

    viper.BindPFlag("sync.gitlab.token", syncGitlabCmd.PersistentFlags().Lookup("token"))
    viper.BindPFlag("sync.gitlab.issue_prefix", syncGitlabCmd.PersistentFlags().Lookup("issue-prefix"))
    viper.BindPFlag("sync.gitlab.merge_request_prefix", syncGitlabCmd.PersistentFlags().Lookup("merge-request-prefix"))
    viper.BindPFlag("sync.gitlab.api_base_url", syncGitlabCmd.PersistentFlags().Lookup("api-base-url"))
    viper.BindPFlag("sync.gitlab.issues_endpoint", syncGitlabCmd.PersistentFlags().Lookup("issues-endpoint"))
    viper.BindPFlag("sync.gitlab.merge_requests_endpoint", syncGitlabCmd.PersistentFlags().Lookup("merge-requests-endpoint"))
}
