package openproject

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WorkPackage represents the structure of a work package from OpenProject's API
type WorkPackage struct {
	ID                int         `json:"id"`
	Subject           string      `json:"subject"`
	Description       Description `json:"description"`
	DueDate           string      `json:"dueDate"`
	DerivedDueDate    string      `json:"derivedDueDate"`
	DerivedStartDate  string      `json:"derivedStartDate"`
	LockVersion       int         `json:"lockVersion"`
	ScheduleManually  bool        `json:"scheduleManually"`
	StartDate         string      `json:"startDate"`
	EstimatedTime     interface{} `json:"estimatedTime"`
	Duration          string      `json:"duration"`
	PercentageDone    int         `json:"percentageDone"`
	CreatedAt         string      `json:"createdAt"`
	UpdatedAt         string      `json:"updatedAt"`
	Readonly          bool        `json:"readonly"`
	// Add other fields as necessary
}

// Description represents the structure of the work package's description
type Description struct {
	Format string `json:"format"`
	Raw    string `json:"raw"`
	HTML   string `json:"html"`
}

// Embedded represents the embedded data in the API response
type Embedded struct {
	Elements []WorkPackage `json:"elements"`
}

// ResponseData represents the structure of the response from the OpenProject API
type ResponseData struct {
	Embedded Embedded `json:"_embedded"`
	Total    int      `json:"total"`
	Count    int      `json:"count"`
	PageSize int      `json:"pageSize"`
	Offset   int      `json:"offset"`
	// Add other fields as necessary
}

// GetWorkPackages fetches the work packages from OpenProject and unmarshals the JSON response
func GetWorkPackages(syncOpenProjectUrl, syncOpenProjectApiKey string) (*ResponseData, error) {
	// Prepare the URL and API endpoint
	url := fmt.Sprintf("%s/api/v3/work_packages", syncOpenProjectUrl)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set basic authentication with the API key
	req.SetBasicAuth("apikey", syncOpenProjectApiKey)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get work packages. Status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response
	var data ResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Print the details of each work package
	for _, wp := range data.Embedded.Elements {
		fmt.Printf("ID: %d\n", wp.ID)
		fmt.Printf("Subject: %s\n", wp.Subject)
		fmt.Printf("Description: %s\n", wp.Description.Raw)
		fmt.Printf("Due Date: %s\n", wp.DueDate)
		fmt.Printf("Derived Due Date: %s\n", wp.DerivedDueDate)
		fmt.Printf("Start Date: %s\n", wp.StartDate)
		fmt.Printf("Percentage Done: %d%%\n", wp.PercentageDone)
		fmt.Printf("Created At: %s\n", wp.CreatedAt)
		fmt.Printf("Updated At: %s\n", wp.UpdatedAt)
		fmt.Println("---------------------------------------------------")
	}

	// Return the unmarshaled data
	return &data, nil
}
