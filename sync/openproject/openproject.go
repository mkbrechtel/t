package openproject

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WorkPackage represents the structure of a work package from OpenProject's API
type WorkPackage struct {
    DerivedStartDate   string `json:"derivedStartDate"`
    DerivedDueDate     string `json:"derivedDueDate"`
    Type               string `json:"type"`
    Status             string `json:"status"`
    Project            string `json:"project"`
    ID                 int    `json:"id"`
    LockVersion        int    `json:"lockVersion"`
    Subject            string `json:"subject"`
    Description        struct {
        Format string `json:"format"`
        Raw    string `json:"raw"`
        HTML   string `json:"html"`
    } `json:"description"`
    ScheduleManually   bool `json:"scheduleManually"`
    StartDate          string `json:"startDate"`
    DueDate            string `json:"dueDate"`
    EstimatedTime      *int   `json:"estimatedTime,omitempty"`
    DerivedEstimatedTime *int `json:"derivedEstimatedTime,omitempty"`
    Duration           string `json:"duration"`
    IgnoreNonWorkingDays bool `json:"ignoreNonWorkingDays"`
    PercentageDone     int    `json:"percentageDone"`
    CreatedAt          string `json:"createdAt"`
    UpdatedAt          string `json:"updatedAt"`
    Readonly           bool   `json:"readonly"`
    CustomFields       map[string]interface{} `json:"customFields"` // Flexible for various custom fields

    Links              struct {
        Type struct {
            Href  string `json:"href"`
            Title string `json:"title"`
        } `json:"type"`
        Status struct {
            Href  string `json:"href"`
            Title string `json:"title"`
        } `json:"status"`
        Project struct {
            Href  string `json:"href"`
            Title string `json:"title"`
        } `json:"project"`
    } `json:"_links"` // Correctly formatted JSON tag for the _links field
}

// Description represents the structure of the work package's description
type Description struct {
	Format string `json:"format"`
	Raw    string `json:"raw"`
	HTML   string `json:"html"`
}

// GetWorkPackages fetches work packages from OpenProject using a query ID
func GetWorkPackages(baseUrl, apiKey, queryId string) ([]WorkPackage, error) {
	// Prepare the URL for the query endpoint
	queryUrl := fmt.Sprintf("%s/api/v3/queries/%s", baseUrl, queryId)

	// Create a new HTTP request for the query
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set basic authentication with the API key
	req.SetBasicAuth("apikey", apiKey)

	// Perform the request for the query
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get query. Status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response for the query
	var queryResponse struct {
		Embedded struct {
			Results struct {
				Embedded struct {
					Elements []WorkPackage `json:"elements"`
				} `json:"_embedded"`
			} `json:"results"`
		} `json:"_embedded"`
	}

	err = json.Unmarshal(body, &queryResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Return the unmarshaled work packages
	return queryResponse.Embedded.Results.Embedded.Elements, nil
}

// PrintWorkPackages prints the work packages to the console
func PrintWorkPackages(workPackages []WorkPackage) {
	for _, wp := range workPackages {
		fmt.Printf("ID: %d\n", wp.ID)
		fmt.Printf("Subject: %s\n", wp.Subject)
		fmt.Printf("Type: %s\n", wp.Links.Type.Title)
		fmt.Printf("Project: %s\n", wp.Links.Project.Title)
		fmt.Printf("Status: %s\n", wp.Links.Status.Title)
		fmt.Printf("Due Date: %s\n", wp.DueDate)
		fmt.Printf("Derived Due Date: %s\n", wp.DerivedDueDate)
		fmt.Printf("Start Date: %s\n", wp.StartDate)
		fmt.Printf("Percentage Done: %d%%\n", wp.PercentageDone)
		fmt.Println("-----")
	}
}
