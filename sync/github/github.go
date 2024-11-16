package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"strings"

	todo "github.com/1set/todotxt"
)

// Issue represents a GitHub issue
type Issue struct {
	Title    string    `json:"title"`
	HTMLURL  string    `json:"html_url"`
	State    string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	DueOn    *time.Time `json:"due_on"`
}

// GetUserIssues fetches issues assigned to a user from GitHub and returns them
func GetUserIssues(token, baseURL, endpoint string) ([]Issue, error) {
	url := baseURL + endpoint

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "token "+token)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error! status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var issues []Issue
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return issues, nil
	}

// PrintIssues prints the GitHub issues
func PrintIssues(issues []Issue) {
	fmt.Println("Assigned issues:")
	for _, issue := range issues {
		fmt.Printf("Title: %s\nURL: %s\nState: %s\n\n", issue.Title, issue.HTMLURL, issue.State)
}
}

// CreateTaskList creates a todo.TaskList from GitHub issues
func CreateTaskList(issues []Issue, issuePrefix, pullPrefix string) todo.TaskList {
	tl := todo.NewTaskList()

	for _, issue := range issues {
		task := todo.NewTask()

		if task.AdditionalTags == nil {
			task.AdditionalTags = make(map[string]string)
		}

		// Determine the appropriate prefix based on the issue type
		prefix := issuePrefix
		if isPullRequest(issue) {
			prefix = pullPrefix
		}

		task.Todo = prefix + issue.Title
		task.CreatedDate = issue.CreatedAt

		if issue.DueOn != nil {
			task.DueDate = *issue.DueOn
		}

		task.AdditionalTags["url"] = issue.HTMLURL
		task.AdditionalTags["state"] = issue.State

		tl.AddTask(&task)
	}

	return tl
}

// isPullRequest checks if the issue is actually a pull request
func isPullRequest(issue Issue) bool {
	// You may need to adjust this logic based on how GitHub distinguishes
	// between issues and pull requests in the API response
	return issue.HTMLURL != "" && strings.Contains(issue.HTMLURL, "/pull/")
}
