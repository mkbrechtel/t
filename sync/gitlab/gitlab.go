package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	todo "github.com/1set/todotxt"
)
// Issue represents a GitLab issue
type Issue struct {
	Title     string     `json:"title"`
	WebURL    string     `json:"web_url"`
	State     string     `json:"state"`
	CreatedAt time.Time  `json:"created_at"`
	DueDate   *time.Time `json:"due_date"`
}
// GetUserIssues fetches issues assigned to a user from GitLab and returns them
func GetUserIssues(token, baseURL, endpoint string) ([]Issue, error) {
	url := baseURL + endpoint

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	req.Header.Add("Accept", "application/json")

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
// PrintIssues prints the GitLab issues
func PrintIssues(issues []Issue) {
	fmt.Println("Assigned issues:")
	for _, issue := range issues {
		fmt.Printf("Title: %s\nURL: %s\nState: %s\n\n", issue.Title, issue.WebURL, issue.State)
	}
}
// CreateTaskList creates a todo.TaskList from GitLab issues
func CreateIssueTaskList(issues []Issue, prefix string) todo.TaskList {
	tl := todo.NewTaskList()

	for _, issue := range issues {
		task := todo.NewTask()

		if task.AdditionalTags == nil {
			task.AdditionalTags = make(map[string]string)
		}
		// Determine the appropriate prefix based on the issue type
		task.Todo = prefix + issue.Title
		task.CreatedDate = issue.CreatedAt

		if issue.DueDate != nil {
			task.DueDate = *issue.DueDate
		}
		task.AdditionalTags["url"] = issue.WebURL
		task.AdditionalTags["state"] = issue.State

		tl.AddTask(&task)
	}
	return tl
}