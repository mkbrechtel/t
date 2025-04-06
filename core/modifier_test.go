package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCompletionModifier(t *testing.T) {
	// Create test tasks
	completedTask := Task{Completed: true, CompletedDate: time.Now()}
	uncompletedTask := Task{Completed: false, CompletedDate: time.Time{}}

	tests := []struct {
		name     string
		modifier CompletionModifier
		task     Task
		expected Task
	}{
		{
			name:     "Mark task as completed",
			modifier: CompletionModifier{Complete: true},
			task:     uncompletedTask,
			expected: Task{Completed: true},
		},
		{
			name:     "Mark task as uncompleted",
			modifier: CompletionModifier{Uncomplete: true},
			task:     completedTask,
			expected: Task{Completed: false, CompletedDate: time.Time{}},
		},
		{
			name:     "No change when both flags are false",
			modifier: CompletionModifier{},
			task:     uncompletedTask,
			expected: uncompletedTask,
		},
		{
			name:     "Complete wins when both flags are set",
			modifier: CompletionModifier{Complete: true, Uncomplete: true},
			task:     uncompletedTask,
			expected: Task{Completed: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.ModifyTask(tt.task)
			assert.Equal(t, tt.expected.Completed, result.Completed)
			if tt.expected.Completed {
				assert.False(t, result.CompletedDate.IsZero(), "CompletedDate should be set")
			} else {
				assert.True(t, result.CompletedDate.IsZero(), "CompletedDate should be empty")
			}
		})
	}

	// Test IsActive
	assert.True(t, CompletionModifier{Complete: true}.IsActive())
	assert.True(t, CompletionModifier{Uncomplete: true}.IsActive())
	assert.False(t, CompletionModifier{}.IsActive())
}

func TestPriorityModifier(t *testing.T) {
	// Create test tasks
	taskWithPriority := Task{Priority: "A"}
	taskWithoutPriority := Task{Priority: ""}

	tests := []struct {
		name     string
		modifier PriorityModifier
		task     Task
		expected Task
	}{
		{
			name:     "Set priority",
			modifier: PriorityModifier{Priority: "B"},
			task:     taskWithoutPriority,
			expected: Task{Priority: "B"},
		},
		{
			name:     "Change existing priority",
			modifier: PriorityModifier{Priority: "C"},
			task:     taskWithPriority,
			expected: Task{Priority: "C"},
		},
		{
			name:     "Remove priority",
			modifier: PriorityModifier{RemovePriority: true},
			task:     taskWithPriority,
			expected: Task{Priority: ""},
		},
		{
			name:     "RemovePriority takes precedence over Priority",
			modifier: PriorityModifier{Priority: "B", RemovePriority: true},
			task:     taskWithPriority,
			expected: Task{Priority: ""},
		},
		{
			name:     "No change when no priority specified and not removing",
			modifier: PriorityModifier{},
			task:     taskWithPriority,
			expected: taskWithPriority,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.ModifyTask(tt.task)
			assert.Equal(t, tt.expected.Priority, result.Priority)
		})
	}

	// Test IsActive
	assert.True(t, PriorityModifier{Priority: "A"}.IsActive())
	assert.True(t, PriorityModifier{RemovePriority: true}.IsActive())
	assert.False(t, PriorityModifier{}.IsActive())
}

func TestProjectModifier(t *testing.T) {
	// Create test tasks
	taskWithProjects := Task{Projects: []string{"Home", "Work"}}
	taskWithoutProjects := Task{Projects: []string{}}

	tests := []struct {
		name     string
		modifier ProjectModifier
		task     Task
		expected Task
	}{
		{
			name:     "Add new project to task with existing projects",
			modifier: ProjectModifier{AddProjects: []string{"School"}},
			task:     taskWithProjects,
			expected: Task{Projects: []string{"Home", "Work", "School"}},
		},
		{
			name:     "Add new project to task without projects",
			modifier: ProjectModifier{AddProjects: []string{"Home"}},
			task:     taskWithoutProjects,
			expected: Task{Projects: []string{"Home"}},
		},
		{
			name:     "Don't add duplicate project",
			modifier: ProjectModifier{AddProjects: []string{"Home"}},
			task:     taskWithProjects,
			expected: Task{Projects: []string{"Home", "Work"}},
		},
		{
			name:     "Remove existing project",
			modifier: ProjectModifier{RemoveProjects: []string{"Work"}},
			task:     taskWithProjects,
			expected: Task{Projects: []string{"Home"}},
		},
		{
			name:     "Remove project that doesn't exist (no change)",
			modifier: ProjectModifier{RemoveProjects: []string{"School"}},
			task:     taskWithProjects,
			expected: Task{Projects: []string{"Home", "Work"}},
		},
		{
			name:     "Add and remove different projects",
			modifier: ProjectModifier{AddProjects: []string{"School"}, RemoveProjects: []string{"Work"}},
			task:     taskWithProjects,
			expected: Task{Projects: []string{"Home", "School"}},
		},
		{
			name:     "No change when no projects specified",
			modifier: ProjectModifier{},
			task:     taskWithProjects,
			expected: taskWithProjects,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.ModifyTask(tt.task)
			assert.ElementsMatch(t, tt.expected.Projects, result.Projects)
		})
	}

	// Test IsActive
	assert.True(t, ProjectModifier{AddProjects: []string{"Home"}}.IsActive())
	assert.True(t, ProjectModifier{RemoveProjects: []string{"Work"}}.IsActive())
	assert.False(t, ProjectModifier{}.IsActive())
}

func TestContextModifier(t *testing.T) {
	// Create test tasks
	taskWithContexts := Task{Contexts: []string{"phone", "work"}}
	taskWithoutContexts := Task{Contexts: []string{}}

	tests := []struct {
		name     string
		modifier ContextModifier
		task     Task
		expected Task
	}{
		{
			name:     "Add new context to task with existing contexts",
			modifier: ContextModifier{AddContexts: []string{"email"}},
			task:     taskWithContexts,
			expected: Task{Contexts: []string{"phone", "work", "email"}},
		},
		{
			name:     "Add new context to task without contexts",
			modifier: ContextModifier{AddContexts: []string{"phone"}},
			task:     taskWithoutContexts,
			expected: Task{Contexts: []string{"phone"}},
		},
		{
			name:     "Don't add duplicate context",
			modifier: ContextModifier{AddContexts: []string{"phone"}},
			task:     taskWithContexts,
			expected: Task{Contexts: []string{"phone", "work"}},
		},
		{
			name:     "Remove existing context",
			modifier: ContextModifier{RemoveContexts: []string{"work"}},
			task:     taskWithContexts,
			expected: Task{Contexts: []string{"phone"}},
		},
		{
			name:     "Remove context that doesn't exist (no change)",
			modifier: ContextModifier{RemoveContexts: []string{"email"}},
			task:     taskWithContexts,
			expected: Task{Contexts: []string{"phone", "work"}},
		},
		{
			name:     "Add and remove different contexts",
			modifier: ContextModifier{AddContexts: []string{"email"}, RemoveContexts: []string{"work"}},
			task:     taskWithContexts,
			expected: Task{Contexts: []string{"phone", "email"}},
		},
		{
			name:     "No change when no contexts specified",
			modifier: ContextModifier{},
			task:     taskWithContexts,
			expected: taskWithContexts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.ModifyTask(tt.task)
			assert.ElementsMatch(t, tt.expected.Contexts, result.Contexts)
		})
	}

	// Test IsActive
	assert.True(t, ContextModifier{AddContexts: []string{"phone"}}.IsActive())
	assert.True(t, ContextModifier{RemoveContexts: []string{"work"}}.IsActive())
	assert.False(t, ContextModifier{}.IsActive())
}

func TestDueDateModifier(t *testing.T) {
	// Create test dates
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)

	// Create test tasks
	taskWithDueDate := Task{DueDate: now}
	taskWithoutDueDate := Task{DueDate: time.Time{}}

	tests := []struct {
		name     string
		modifier DueDateModifier
		task     Task
		expected Task
	}{
		{
			name:     "Set due date on task without due date",
			modifier: DueDateModifier{DueDate: &tomorrow},
			task:     taskWithoutDueDate,
			expected: Task{DueDate: tomorrow},
		},
		{
			name:     "Change existing due date",
			modifier: DueDateModifier{DueDate: &tomorrow},
			task:     taskWithDueDate,
			expected: Task{DueDate: tomorrow},
		},
		{
			name:     "Remove due date",
			modifier: DueDateModifier{RemoveDueDate: true},
			task:     taskWithDueDate,
			expected: Task{DueDate: time.Time{}},
		},
		{
			name:     "RemoveDueDate takes precedence over DueDate",
			modifier: DueDateModifier{DueDate: &tomorrow, RemoveDueDate: true},
			task:     taskWithDueDate,
			expected: Task{DueDate: time.Time{}},
		},
		{
			name:     "No change when no due date specified and not removing",
			modifier: DueDateModifier{},
			task:     taskWithDueDate,
			expected: taskWithDueDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.ModifyTask(tt.task)
			assert.Equal(t, tt.expected.DueDate, result.DueDate)
		})
	}

	// Test IsActive
	assert.True(t, DueDateModifier{DueDate: &tomorrow}.IsActive())
	assert.True(t, DueDateModifier{RemoveDueDate: true}.IsActive())
	assert.False(t, DueDateModifier{}.IsActive())
}

func TestTextModifier(t *testing.T) {
	// Create test tasks
	task := Task{Todo: "Buy groceries"}

	tests := []struct {
		name     string
		modifier TextModifier
		task     Task
		expected Task
	}{
		{
			name:     "Append text",
			modifier: TextModifier{AppendText: " for dinner"},
			task:     task,
			expected: Task{Todo: "Buy groceries for dinner"},
		},
		{
			name:     "Prepend text",
			modifier: TextModifier{PrependText: "Today: "},
			task:     task,
			expected: Task{Todo: "Today: Buy groceries"},
		},
		{
			name:     "Both append and prepend",
			modifier: TextModifier{AppendText: " for dinner", PrependText: "Today: "},
			task:     task,
			expected: Task{Todo: "Today: Buy groceries for dinner"},
		},
		{
			name:     "No change when no text specified",
			modifier: TextModifier{},
			task:     task,
			expected: task,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.ModifyTask(tt.task)
			assert.Equal(t, tt.expected.Todo, result.Todo)
		})
	}

	// Test IsActive
	assert.True(t, TextModifier{AppendText: " for dinner"}.IsActive())
	assert.True(t, TextModifier{PrependText: "Today: "}.IsActive())
	assert.False(t, TextModifier{}.IsActive())
}

func TestChainModifier(t *testing.T) {
	// Create test task
	task := Task{
		Todo:      "Buy groceries",
		Projects:  []string{"Home"},
		Contexts:  []string{"shopping"},
		Priority:  "B",
		Completed: false,
	}

	// Create individual modifiers
	priorityModifier := PriorityModifier{Priority: "A"}
	projectModifier := ProjectModifier{AddProjects: []string{"Personal"}}
	textModifier := TextModifier{AppendText: " for dinner"}

	// Create test cases
	tests := []struct {
		name     string
		modifier ChainModifier
		task     Task
		expected Task
	}{
		{
			name: "Apply all modifiers in chain",
			modifier: ChainModifier{Modifiers: []TaskModifier{
				priorityModifier,
				projectModifier,
				textModifier,
			}},
			task: task,
			expected: Task{
				Todo:      "Buy groceries for dinner",
				Projects:  []string{"Home", "Personal"},
				Contexts:  []string{"shopping"},
				Priority:  "A",
				Completed: false,
			},
		},
		{
			name:     "No change when no modifiers in chain",
			modifier: ChainModifier{Modifiers: []TaskModifier{}},
			task:     task,
			expected: task,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.modifier.ModifyTask(tt.task)
			
			// Check each expected field
			assert.Equal(t, tt.expected.Todo, result.Todo)
			assert.ElementsMatch(t, tt.expected.Projects, result.Projects)
			assert.ElementsMatch(t, tt.expected.Contexts, result.Contexts)
			assert.Equal(t, tt.expected.Priority, result.Priority)
			assert.Equal(t, tt.expected.Completed, result.Completed)
		})
	}
}