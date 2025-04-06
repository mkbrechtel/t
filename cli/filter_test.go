package cli

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"t5.mkbrechtel.dev/t5/core"
)

func TestFilterComposer(t *testing.T) {
	// Test creating a new filter composer
	fc := NewFilterComposer()
	assert.NotNil(t, fc)
	assert.Empty(t, fc.Filters.Filters)
	
	// Test adding a filter
	projectFilter := &core.ProjectFilter{Projects: []string{"Home"}}
	fc.AddFilter(projectFilter)
	assert.Len(t, fc.Filters.Filters, 1)
	
	// Test composing a single filter
	composed := fc.ComposeFilter()
	assert.Equal(t, projectFilter, composed)
	
	// Test composing multiple filters
	contextFilter := &core.ContextFilter{Contexts: []string{"phone"}}
	fc.AddFilter(contextFilter)
	composed = fc.ComposeFilter()
	assert.IsType(t, &core.AndFilter{}, composed)
	
	// Test composing with inactive filters
	fc = NewFilterComposer()
	fc.AddFilter(&core.ProjectFilter{}) // Empty project filter is inactive
	assert.Nil(t, fc.ComposeFilter())
}

func TestCompletionFlag(t *testing.T) {
	// Test the CompletionFlag for completed tasks
	fc := NewFilterComposer()
	completedFlag := &fc.CompletedFlag
	
	// Test String method
	assert.Equal(t, "false", completedFlag.String())
	
	// Test Set method
	err := completedFlag.Set("true")
	assert.NoError(t, err)
	assert.True(t, completedFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test the filter's behavior
	completedTask := core.Task{Completed: true}
	incompleteTask := core.Task{Completed: false}
	
	assert.True(t, composed.FilterTask(completedTask))
	assert.False(t, composed.FilterTask(incompleteTask))
}

func TestNotCompletedFlag(t *testing.T) {
	// Test the CompletionFlag for not-completed tasks
	fc := NewFilterComposer()
	notCompletedFlag := &fc.NotCompletedFlag
	
	// Test Set method
	err := notCompletedFlag.Set("true")
	assert.NoError(t, err)
	assert.True(t, notCompletedFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test the filter's behavior
	completedTask := core.Task{Completed: true}
	incompleteTask := core.Task{Completed: false}
	
	assert.False(t, composed.FilterTask(completedTask))
	assert.True(t, composed.FilterTask(incompleteTask))
}

func TestPriorityFlag(t *testing.T) {
	// Test the PriorityFlag
	fc := NewFilterComposer()
	priorityFlag := &fc.PriorityFlag
	
	// Test String method
	assert.Equal(t, "", priorityFlag.String())
	
	// Test Set method with single priority
	err := priorityFlag.Set("A")
	assert.NoError(t, err)
	assert.Equal(t, "A", priorityFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test the filter's behavior
	taskPriorityA := core.Task{Priority: "A"}
	taskPriorityB := core.Task{Priority: "B"}
	
	assert.True(t, composed.FilterTask(taskPriorityA))
	assert.False(t, composed.FilterTask(taskPriorityB))
	
	// Test with multiple priorities
	fc = NewFilterComposer()
	priorityFlag = &fc.PriorityFlag
	
	err = priorityFlag.Set("A,B")
	assert.NoError(t, err)
	
	composed = fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Should match either A or B
	assert.True(t, composed.FilterTask(taskPriorityA))
	assert.True(t, composed.FilterTask(taskPriorityB))
	assert.False(t, composed.FilterTask(core.Task{Priority: "C"}))
}

func TestProjectFlag(t *testing.T) {
	// Test the ProjectFlag
	fc := NewFilterComposer()
	projectFlag := &fc.ProjectFlag
	
	// Test String method
	assert.Equal(t, "", projectFlag.String())
	
	// Test Set method
	err := projectFlag.Set("Home,Work")
	assert.NoError(t, err)
	assert.Equal(t, "Home,Work", projectFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test the filter's behavior
	taskHomeProject := core.Task{Projects: []string{"Home"}}
	taskWorkProject := core.Task{Projects: []string{"Work"}}
	taskOtherProject := core.Task{Projects: []string{"Other"}}
	
	assert.True(t, composed.FilterTask(taskHomeProject))
	assert.True(t, composed.FilterTask(taskWorkProject))
	assert.False(t, composed.FilterTask(taskOtherProject))
	
	// Test with + prefix in project name
	fc = NewFilterComposer()
	projectFlag = &fc.ProjectFlag
	
	err = projectFlag.Set("+Home")
	assert.NoError(t, err)
	
	composed = fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	assert.True(t, composed.FilterTask(taskHomeProject))
}

func TestContextFlag(t *testing.T) {
	// Test the ContextFlag
	fc := NewFilterComposer()
	contextFlag := &fc.ContextFlag
	
	// Test String method
	assert.Equal(t, "", contextFlag.String())
	
	// Test Set method
	err := contextFlag.Set("phone,email")
	assert.NoError(t, err)
	assert.Equal(t, "phone,email", contextFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test the filter's behavior
	taskPhoneContext := core.Task{Contexts: []string{"phone"}}
	taskEmailContext := core.Task{Contexts: []string{"email"}}
	taskOtherContext := core.Task{Contexts: []string{"meeting"}}
	
	assert.True(t, composed.FilterTask(taskPhoneContext))
	assert.True(t, composed.FilterTask(taskEmailContext))
	assert.False(t, composed.FilterTask(taskOtherContext))
	
	// Test with @ prefix in context name
	fc = NewFilterComposer()
	contextFlag = &fc.ContextFlag
	
	err = contextFlag.Set("@phone")
	assert.NoError(t, err)
	
	composed = fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	assert.True(t, composed.FilterTask(taskPhoneContext))
}

func TestDueTodayFlag(t *testing.T) {
	// Test the DueTodayFlag
	fc := NewFilterComposer()
	dueTodayFlag := &fc.DueTodayFlag
	
	// Test String method
	assert.Equal(t, "false", dueTodayFlag.String())
	
	// Test Set method
	err := dueTodayFlag.Set("true")
	assert.NoError(t, err)
	assert.True(t, dueTodayFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test the filter's behavior
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)
	yesterday := today.AddDate(0, 0, -1)
	
	taskDueToday := core.Task{DueDate: today}
	taskDueTomorrow := core.Task{DueDate: tomorrow}
	taskDueYesterday := core.Task{DueDate: yesterday}
	
	assert.True(t, composed.FilterTask(taskDueToday))
	assert.False(t, composed.FilterTask(taskDueTomorrow))
	assert.False(t, composed.FilterTask(taskDueYesterday))
}

func TestDueBeforeFlag(t *testing.T) {
	// Test the DueBeforeFlag
	fc := NewFilterComposer()
	dueBeforeFlag := &fc.DueBeforeFlag
	
	// Test String method
	assert.Equal(t, "", dueBeforeFlag.String())
	
	// Test Set method
	err := dueBeforeFlag.Set("2023-12-31")
	assert.NoError(t, err)
	assert.Equal(t, "2023-12-31", dueBeforeFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test invalid date format
	fc = NewFilterComposer()
	dueBeforeFlag = &fc.DueBeforeFlag
	
	err = dueBeforeFlag.Set("invalid-date")
	assert.Error(t, err)
}

func TestCreatedAfterFlag(t *testing.T) {
	// Test the CreatedAfterFlag
	fc := NewFilterComposer()
	createdAfterFlag := &fc.CreatedAfterFlag
	
	// Test String method
	assert.Equal(t, "", createdAfterFlag.String())
	
	// Test Set method
	err := createdAfterFlag.Set("2023-01-01")
	assert.NoError(t, err)
	assert.Equal(t, "2023-01-01", createdAfterFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test invalid date format
	fc = NewFilterComposer()
	createdAfterFlag = &fc.CreatedAfterFlag
	
	err = createdAfterFlag.Set("invalid-date")
	assert.Error(t, err)
}

func TestTagFlag(t *testing.T) {
	// Test the TagFlag
	fc := NewFilterComposer()
	tagFlag := &fc.HasTagFlag
	
	// Test String method
	assert.Equal(t, "", tagFlag.String())
	
	// Test Set method
	err := tagFlag.Set("due")
	assert.NoError(t, err)
	assert.Equal(t, "due", tagFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test the filter's behavior
	taskWithTag := core.Task{AdditionalTags: map[string]string{"due": "2023-12-31"}}
	taskWithoutTag := core.Task{AdditionalTags: map[string]string{"other": "value"}}
	
	assert.True(t, composed.FilterTask(taskWithTag))
	assert.False(t, composed.FilterTask(taskWithoutTag))
}

func TestRegexFlag(t *testing.T) {
	// Test the RegexFlag
	fc := NewFilterComposer()
	regexFlag := &fc.RegexFlag
	
	// Test String method
	assert.Equal(t, "", regexFlag.String())
	
	// Test Set method
	err := regexFlag.Set("groceries|shopping")
	assert.NoError(t, err)
	assert.Equal(t, "groceries|shopping", regexFlag.value)
	
	// Verify a filter was added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	
	// Test the filter's behavior
	taskMatching := core.Task{Todo: "Buy groceries"}
	taskNotMatching := core.Task{Todo: "Call dentist"}
	
	assert.True(t, composed.FilterTask(taskMatching))
	assert.False(t, composed.FilterTask(taskNotMatching))
}

func TestFilterFlagIntegration(t *testing.T) {
	// Test integrating with the flag package
	fc := NewFilterComposer()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	
	// Add all filter flags
	fs.Var(&fc.CompletedFlag, "completed", "Show only completed tasks")
	fs.Var(&fc.NotCompletedFlag, "not-completed", "Show only non-completed tasks")
	fs.Var(&fc.PriorityFlag, "priority", "Filter by priority")
	fs.Var(&fc.ProjectFlag, "project", "Filter by project")
	fs.Var(&fc.ContextFlag, "context", "Filter by context")
	
	// Parse some flags
	args := []string{"--completed=true", "--project=Home"}
	err := fs.Parse(args)
	assert.NoError(t, err)
	
	// Verify flags were set correctly
	assert.True(t, fc.CompletedFlag.value)
	assert.Equal(t, "Home", fc.ProjectFlag.value)
	
	// Verify filters were added
	composed := fc.ComposeFilter()
	assert.NotNil(t, composed)
	assert.IsType(t, &core.AndFilter{}, composed)
	
	// Test the combined filter behavior
	taskCompletedHome := core.Task{Completed: true, Projects: []string{"Home"}}
	taskCompletedWork := core.Task{Completed: true, Projects: []string{"Work"}}
	taskIncompleteHome := core.Task{Completed: false, Projects: []string{"Home"}}
	
	assert.True(t, composed.FilterTask(taskCompletedHome))
	assert.False(t, composed.FilterTask(taskCompletedWork))
	assert.False(t, composed.FilterTask(taskIncompleteHome))
}