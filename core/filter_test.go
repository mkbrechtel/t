package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProjectFilter(t *testing.T) {
	// Create test tasks with different projects
	taskWithProject := Task{
		Projects: []string{"Home", "Work"},
	}
	taskWithoutProject := Task{
		Projects: []string{"Personal"},
	}
	emptyTask := Task{}

	tests := []struct {
		name     string
		filter   ProjectFilter
		task     Task
		expected bool
	}{
		{
			name:     "Match when project is present",
			filter:   ProjectFilter{Projects: []string{"Home"}},
			task:     taskWithProject,
			expected: true,
		},
		{
			name:     "Match when any project in list is present",
			filter:   ProjectFilter{Projects: []string{"Work", "Other"}},
			task:     taskWithProject,
			expected: true,
		},
		{
			name:     "No match when project is not present",
			filter:   ProjectFilter{Projects: []string{"Other"}},
			task:     taskWithProject,
			expected: false,
		},
		{
			name:     "Match everything when no projects specified",
			filter:   ProjectFilter{Projects: []string{}},
			task:     taskWithoutProject,
			expected: true,
		},
		{
			name:     "Match empty task when no projects specified",
			filter:   ProjectFilter{Projects: []string{}},
			task:     emptyTask,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test IsActive
	assert.True(t, ProjectFilter{Projects: []string{"Home"}}.IsActive())
	assert.False(t, ProjectFilter{Projects: []string{}}.IsActive())
}

func TestContextFilter(t *testing.T) {
	// Create test tasks with different contexts
	taskWithContext := Task{
		Contexts: []string{"phone", "email"},
	}
	taskWithoutContext := Task{
		Contexts: []string{"meeting"},
	}
	emptyTask := Task{}

	tests := []struct {
		name     string
		filter   ContextFilter
		task     Task
		expected bool
	}{
		{
			name:     "Match when context is present",
			filter:   ContextFilter{Contexts: []string{"phone"}},
			task:     taskWithContext,
			expected: true,
		},
		{
			name:     "Match when any context in list is present",
			filter:   ContextFilter{Contexts: []string{"email", "other"}},
			task:     taskWithContext,
			expected: true,
		},
		{
			name:     "No match when context is not present",
			filter:   ContextFilter{Contexts: []string{"other"}},
			task:     taskWithContext,
			expected: false,
		},
		{
			name:     "Match everything when no contexts specified",
			filter:   ContextFilter{Contexts: []string{}},
			task:     taskWithoutContext,
			expected: true,
		},
		{
			name:     "Match empty task when no contexts specified",
			filter:   ContextFilter{Contexts: []string{}},
			task:     emptyTask,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test IsActive
	assert.True(t, ContextFilter{Contexts: []string{"phone"}}.IsActive())
	assert.False(t, ContextFilter{Contexts: []string{}}.IsActive())
}

func TestCompletionFilter(t *testing.T) {
	// Create completed and uncompleted tasks
	completedTask := Task{Completed: true}
	uncompletedTask := Task{Completed: false}

	// Create test cases
	completed := true
	uncompleted := false

	tests := []struct {
		name     string
		filter   CompletionFilter
		task     Task
		expected bool
	}{
		{
			name:     "Match completed task with UseFlag",
			filter:   CompletionFilter{Completed: &completed, UseFlag: true},
			task:     completedTask,
			expected: true,
		},
		{
			name:     "Not match uncompleted task with UseFlag",
			filter:   CompletionFilter{Completed: &completed, UseFlag: true},
			task:     uncompletedTask,
			expected: false,
		},
		{
			name:     "Match uncompleted task with NotFlag",
			filter:   CompletionFilter{Completed: &uncompleted, NotFlag: true},
			task:     uncompletedTask,
			expected: true,
		},
		{
			name:     "Not match completed task with NotFlag",
			filter:   CompletionFilter{Completed: &uncompleted, NotFlag: true},
			task:     completedTask,
			expected: false,
		},
		{
			name:     "Match all tasks when no flags set",
			filter:   CompletionFilter{},
			task:     completedTask,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test IsActive
	assert.True(t, CompletionFilter{Completed: &completed, UseFlag: true}.IsActive())
	assert.True(t, CompletionFilter{Completed: &uncompleted, NotFlag: true}.IsActive())
	assert.False(t, CompletionFilter{}.IsActive())
}

func TestPriorityFilter(t *testing.T) {
	// Create tasks with different priorities
	taskPriorityA := Task{Priority: "A"}
	taskPriorityB := Task{Priority: "B"}
	taskNoPriority := Task{Priority: ""}

	tests := []struct {
		name     string
		filter   PriorityFilter
		task     Task
		expected bool
	}{
		{
			name:     "Match when priority matches exactly",
			filter:   PriorityFilter{Priorities: []string{"A"}},
			task:     taskPriorityA,
			expected: true,
		},
		{
			name:     "Match when any priority in list matches",
			filter:   PriorityFilter{Priorities: []string{"A", "C"}},
			task:     taskPriorityA,
			expected: true,
		},
		{
			name:     "No match when priority doesn't match",
			filter:   PriorityFilter{Priorities: []string{"C", "D"}},
			task:     taskPriorityB,
			expected: false,
		},
		{
			name:     "Match everything when no priorities specified",
			filter:   PriorityFilter{Priorities: []string{}},
			task:     taskPriorityA,
			expected: true,
		},
		{
			name:     "No match when task has no priority",
			filter:   PriorityFilter{Priorities: []string{"A"}},
			task:     taskNoPriority,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test IsActive
	assert.True(t, PriorityFilter{Priorities: []string{"A"}}.IsActive())
	assert.False(t, PriorityFilter{Priorities: []string{}}.IsActive())
}

func TestDueDateFilter(t *testing.T) {
	// Create test dates
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	nextWeek := now.AddDate(0, 0, 7)

	// Create tasks with different due dates
	taskDueTomorrow := Task{DueDate: tomorrow}
	taskDueNextWeek := Task{DueDate: nextWeek}
	taskNoDueDate := Task{DueDate: time.Time{}}

	tests := []struct {
		name     string
		filter   DueDateFilter
		task     Task
		expected bool
	}{
		{
			name:     "Match when due date is before filter date",
			filter:   DueDateFilter{Before: &nextWeek},
			task:     taskDueTomorrow,
			expected: true,
		},
		{
			name:     "No match when due date is after filter date",
			filter:   DueDateFilter{Before: &tomorrow},
			task:     taskDueNextWeek,
			expected: false,
		},
		{
			name:     "Match when due date is after filter date",
			filter:   DueDateFilter{After: &yesterday},
			task:     taskDueTomorrow,
			expected: true,
		},
		{
			name:     "No match when due date is not after filter date",
			filter:   DueDateFilter{After: &tomorrow},
			task:     taskDueTomorrow,
			expected: false,
		},
		{
			name:     "Match when due date is between filter dates",
			filter:   DueDateFilter{Before: &nextWeek, After: &yesterday},
			task:     taskDueTomorrow,
			expected: true,
		},
		{
			name:     "No match when task has no due date",
			filter:   DueDateFilter{Before: &nextWeek},
			task:     taskNoDueDate,
			expected: false,
		},
		{
			name:     "Match everything when no dates specified",
			filter:   DueDateFilter{},
			task:     taskDueTomorrow,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test IsActive
	assert.True(t, DueDateFilter{Before: &tomorrow}.IsActive())
	assert.True(t, DueDateFilter{After: &yesterday}.IsActive())
	assert.True(t, DueDateFilter{Today: true}.IsActive())
	assert.False(t, DueDateFilter{}.IsActive())
}

func TestCreatedDateFilter(t *testing.T) {
	// Create test dates
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -7)

	// Create tasks with different created dates
	taskCreatedYesterday := Task{CreatedDate: yesterday}
	taskCreatedLastWeek := Task{CreatedDate: lastWeek}
	taskNoCreatedDate := Task{CreatedDate: time.Time{}}

	tests := []struct {
		name     string
		filter   CreatedDateFilter
		task     Task
		expected bool
	}{
		{
			name:     "Match when created date is before filter date",
			filter:   CreatedDateFilter{Before: &now},
			task:     taskCreatedYesterday,
			expected: true,
		},
		{
			name:     "No match when created date is after filter date",
			filter:   CreatedDateFilter{Before: &lastWeek},
			task:     taskCreatedYesterday,
			expected: false,
		},
		{
			name:     "Match when created date is after filter date",
			filter:   CreatedDateFilter{After: &lastWeek},
			task:     taskCreatedYesterday,
			expected: true,
		},
		{
			name:     "No match when created date is before filter date",
			filter:   CreatedDateFilter{After: &yesterday},
			task:     taskCreatedLastWeek,
			expected: false,
		},
		{
			name:     "No match when task has no created date",
			filter:   CreatedDateFilter{After: &lastWeek},
			task:     taskNoCreatedDate,
			expected: false,
		},
		{
			name:     "Match everything when no dates specified",
			filter:   CreatedDateFilter{},
			task:     taskCreatedYesterday,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test IsActive
	assert.True(t, CreatedDateFilter{Before: &now}.IsActive())
	assert.True(t, CreatedDateFilter{After: &yesterday}.IsActive())
	assert.False(t, CreatedDateFilter{}.IsActive())
}

func TestTagFilter(t *testing.T) {
	// Create tasks with different tags
	taskWithTags := Task{
		AdditionalTags: map[string]string{
			"test": "value",
			"due":  "2023-05-01",
		},
	}
	taskWithoutTags := Task{}

	tests := []struct {
		name     string
		filter   TagFilter
		task     Task
		expected bool
	}{
		{
			name:     "Match when tag exists",
			filter:   TagFilter{Tag: "test"},
			task:     taskWithTags,
			expected: true,
		},
		{
			name:     "No match when tag doesn't exist",
			filter:   TagFilter{Tag: "nonexistent"},
			task:     taskWithTags,
			expected: false,
		},
		{
			name:     "No match on task without tags",
			filter:   TagFilter{Tag: "test"},
			task:     taskWithoutTags,
			expected: false,
		},
		{
			name:     "Match everything when no tag specified",
			filter:   TagFilter{Tag: ""},
			task:     taskWithTags,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test IsActive
	assert.True(t, TagFilter{Tag: "test"}.IsActive())
	assert.False(t, TagFilter{Tag: ""}.IsActive())
}

func TestRegexFilter(t *testing.T) {
	// Create tasks with different text
	taskWithText := Task{
		Todo: "Buy groceries for dinner party",
	}
	taskWithDifferentText := Task{
		Todo: "Call dentist for appointment",
	}

	tests := []struct {
		name     string
		filter   RegexFilter
		task     Task
		expected bool
	}{
		{
			name:     "Match when regex matches",
			filter:   RegexFilter{Pattern: "groceries|dinner"},
			task:     taskWithText,
			expected: true,
		},
		{
			name:     "No match when regex doesn't match",
			filter:   RegexFilter{Pattern: "dentist|doctor"},
			task:     taskWithText,
			expected: false,
		},
		{
			name:     "Match with case insensitive regex",
			filter:   RegexFilter{Pattern: "(?i)GROCERIES"},
			task:     taskWithText,
			expected: true,
		},
		{
			name:     "Match everything when no pattern specified",
			filter:   RegexFilter{Pattern: ""},
			task:     taskWithDifferentText,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test IsActive
	assert.True(t, RegexFilter{Pattern: "test"}.IsActive())
	assert.False(t, RegexFilter{Pattern: ""}.IsActive())
}

func TestAndFilter(t *testing.T) {
	// Create test task
	task := Task{
		Todo:      "Buy groceries",
		Projects:  []string{"Home"},
		Contexts:  []string{"shopping"},
		Priority:  "A",
		Completed: false,
	}

	// Create individual filters
	projectFilter := ProjectFilter{Projects: []string{"Home"}}
	contextFilter := ContextFilter{Contexts: []string{"shopping"}}
	priorityFilter := PriorityFilter{Priorities: []string{"A"}}
	completed := false
	completionFilter := CompletionFilter{Completed: &completed, NotFlag: true}

	tests := []struct {
		name     string
		filter   AndFilter
		task     Task
		expected bool
	}{
		{
			name: "Match when all filters match",
			filter: AndFilter{Filters: []TaskFilter{
				projectFilter,
				contextFilter,
				priorityFilter,
				completionFilter,
			}},
			task:     task,
			expected: true,
		},
		{
			name: "No match when one filter doesn't match",
			filter: AndFilter{Filters: []TaskFilter{
				projectFilter,
				contextFilter,
				PriorityFilter{Priorities: []string{"B"}}, // Different priority
			}},
			task:     task,
			expected: false,
		},
		{
			name:     "Match when no filters specified",
			filter:   AndFilter{Filters: []TaskFilter{}},
			task:     task,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOrFilter(t *testing.T) {
	// Create test task
	task := Task{
		Todo:      "Buy groceries",
		Projects:  []string{"Home"},
		Contexts:  []string{"shopping"},
		Priority:  "A",
		Completed: false,
	}

	// Create matching and non-matching filters
	matchingFilter := ProjectFilter{Projects: []string{"Home"}}
	nonMatchingFilter := ProjectFilter{Projects: []string{"Work"}}

	tests := []struct {
		name     string
		filter   OrFilter
		task     Task
		expected bool
	}{
		{
			name: "Match when any filter matches",
			filter: OrFilter{Filters: []TaskFilter{
				matchingFilter,
				nonMatchingFilter,
			}},
			task:     task,
			expected: true,
		},
		{
			name: "No match when no filters match",
			filter: OrFilter{Filters: []TaskFilter{
				nonMatchingFilter,
				ProjectFilter{Projects: []string{"School"}},
			}},
			task:     task,
			expected: false,
		},
		{
			name:     "No match when no filters specified",
			filter:   OrFilter{Filters: []TaskFilter{}},
			task:     task,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotFilter(t *testing.T) {
	// Create test task
	task := Task{
		Todo:      "Buy groceries",
		Projects:  []string{"Home"},
		Contexts:  []string{"shopping"},
		Priority:  "A",
		Completed: false,
	}

	// Create matching and non-matching filters
	matchingFilter := ProjectFilter{Projects: []string{"Home"}}
	nonMatchingFilter := ProjectFilter{Projects: []string{"Work"}}

	tests := []struct {
		name     string
		filter   NotFilter
		task     Task
		expected bool
	}{
		{
			name:     "No match when inner filter matches",
			filter:   NotFilter{Filter: matchingFilter},
			task:     task,
			expected: false,
		},
		{
			name:     "Match when inner filter doesn't match",
			filter:   NotFilter{Filter: nonMatchingFilter},
			task:     task,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.FilterTask(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}
}