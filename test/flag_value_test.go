package test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"t5.mkbrechtel.dev/t5/cli"
	"t5.mkbrechtel.dev/t5/core"
)

// TestFilterComposerFlags tests the flag.Value implementation for FilterComposer
func TestFilterComposerFlags(t *testing.T) {
	// Create a new FilterComposer
	fc := cli.NewFilterComposer()

	// Create a FlagSet for testing
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	// Register the filter flags
	fs.Var(&fc.CompletedFlag, "completed", "Show only completed tasks")
	fs.Var(&fc.NotCompletedFlag, "not-completed", "Show only non-completed tasks")
	fs.Var(&fc.PriorityFlag, "priority", "Filter by priority")
	fs.Var(&fc.ProjectFlag, "project", "Filter by project")
	fs.Var(&fc.ContextFlag, "context", "Filter by context")
	fs.Var(&fc.DueTodayFlag, "due-today", "Filter tasks due today")
	fs.Var(&fc.DueBeforeFlag, "due-before", "Filter tasks due before date")
	fs.Var(&fc.CreatedAfterFlag, "created-after", "Filter tasks created after date")
	fs.Var(&fc.HasTagFlag, "has-tag", "Filter by tag")
	fs.Var(&fc.RegexFlag, "regex", "Filter by regex")

	// Test parsing various flags
	t.Run("CompletedFlag", func(t *testing.T) {
		// Reset the FlagSet
		fc = cli.NewFilterComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&fc.CompletedFlag, "completed", "Show only completed tasks")

		// Parse the completed flag (boolean flags require a value for tests)
		err := fs.Parse([]string{"--completed=true"})
		assert.NoError(t, err)

		// Verify filter is created correctly
		filter := fc.ComposeFilter()
		assert.NotNil(t, filter)
		
		// Create a completed task
		completeTask := core.Task{Completed: true}
		// Create an incomplete task
		incompleteTask := core.Task{Completed: false}
		
		// The filter should accept completed tasks and reject incomplete ones
		assert.True(t, filter.FilterTask(completeTask))
		assert.False(t, filter.FilterTask(incompleteTask))
	})

	t.Run("NotCompletedFlag", func(t *testing.T) {
		// Reset the FlagSet and FilterComposer
		fc = cli.NewFilterComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&fc.NotCompletedFlag, "not-completed", "Show only non-completed tasks")

		// Parse the not-completed flag
		err := fs.Parse([]string{"--not-completed=true"})
		assert.NoError(t, err)

		// Verify filter is created correctly
		filter := fc.ComposeFilter()
		assert.NotNil(t, filter)
		
		// Create a completed task
		completeTask := core.Task{Completed: true}
		// Create an incomplete task
		incompleteTask := core.Task{Completed: false}
		
		// The filter should reject completed tasks and accept incomplete ones
		assert.False(t, filter.FilterTask(completeTask))
		assert.True(t, filter.FilterTask(incompleteTask))
	})

	t.Run("PriorityFlag", func(t *testing.T) {
		// Reset the FlagSet and FilterComposer
		fc = cli.NewFilterComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&fc.PriorityFlag, "priority", "Filter by priority")

		// Parse the priority flag
		err := fs.Parse([]string{"--priority=A,B"})
		assert.NoError(t, err)

		// Verify filter is created correctly
		filter := fc.ComposeFilter()
		assert.NotNil(t, filter)
		
		// Create tasks with different priorities
		taskA := core.Task{Priority: "A"}
		taskB := core.Task{Priority: "B"}
		taskC := core.Task{Priority: "C"}
		
		// The filter should accept tasks with priority A or B and reject C
		assert.True(t, filter.FilterTask(taskA))
		assert.True(t, filter.FilterTask(taskB))
		assert.False(t, filter.FilterTask(taskC))
	})

	t.Run("ProjectFlag", func(t *testing.T) {
		// Reset the FlagSet and FilterComposer
		fc = cli.NewFilterComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&fc.ProjectFlag, "project", "Filter by project")

		// Parse the project flag
		err := fs.Parse([]string{"--project=Home,Work"})
		assert.NoError(t, err)

		// Verify filter is created correctly
		filter := fc.ComposeFilter()
		assert.NotNil(t, filter)
		
		// Create tasks with different projects
		taskHome := core.Task{Projects: []string{"Home"}}
		taskWork := core.Task{Projects: []string{"Work"}}
		taskOther := core.Task{Projects: []string{"Other"}}
		
		// The filter should accept tasks with Home or Work projects and reject Other
		assert.True(t, filter.FilterTask(taskHome))
		assert.True(t, filter.FilterTask(taskWork))
		assert.False(t, filter.FilterTask(taskOther))
	})

	t.Run("ContextFlag", func(t *testing.T) {
		// Reset the FlagSet and FilterComposer
		fc = cli.NewFilterComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&fc.ContextFlag, "context", "Filter by context")

		// Parse the context flag
		err := fs.Parse([]string{"--context=phone,computer"})
		assert.NoError(t, err)

		// Verify filter is created correctly
		filter := fc.ComposeFilter()
		assert.NotNil(t, filter)
		
		// Create tasks with different contexts
		taskPhone := core.Task{Contexts: []string{"phone"}}
		taskComputer := core.Task{Contexts: []string{"computer"}}
		taskOther := core.Task{Contexts: []string{"other"}}
		
		// The filter should accept tasks with phone or computer contexts and reject other
		assert.True(t, filter.FilterTask(taskPhone))
		assert.True(t, filter.FilterTask(taskComputer))
		assert.False(t, filter.FilterTask(taskOther))
	})
}

// TestModifierComposerFlags tests the flag.Value implementation for ModifierComposer
func TestModifierComposerFlags(t *testing.T) {
	// Create a new ModifierComposer
	mc := cli.NewModifierComposer()

	// Create a FlagSet for testing
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	// Register the modifier flags
	fs.Var(&mc.CompleteFlag, "complete", "Mark task as complete")
	fs.Var(&mc.UncompleteFlag, "uncomplete", "Mark task as incomplete")
	fs.Var(&mc.PriorityFlag, "priority", "Set priority")
	fs.Var(&mc.RemovePriorityFlag, "remove-priority", "Remove priority")
	fs.Var(&mc.AddProjectFlag, "add-project", "Add project")
	fs.Var(&mc.RemoveProjectFlag, "remove-project", "Remove project")
	fs.Var(&mc.AddContextFlag, "add-context", "Add context")
	fs.Var(&mc.RemoveContextFlag, "remove-context", "Remove context")
	fs.Var(&mc.DueFlag, "due", "Set due date")
	fs.Var(&mc.RemoveDueFlag, "remove-due", "Remove due date")
	fs.Var(&mc.AppendFlag, "append", "Append text")
	fs.Var(&mc.PrependFlag, "prepend", "Prepend text")

	// Test parsing various flags
	t.Run("CompleteFlag", func(t *testing.T) {
		// Reset the FlagSet and ModifierComposer
		mc = cli.NewModifierComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&mc.CompleteFlag, "complete", "Mark task as complete")

		// Parse the complete flag
		err := fs.Parse([]string{"--complete=true"})
		assert.NoError(t, err)

		// Verify modifier is created correctly
		modifier := mc.ComposeModifier()
		assert.NotNil(t, modifier)
		
		// Create an incomplete task
		task := core.Task{Completed: false}
		
		// The modifier should mark the task as complete
		modifiedTask := modifier.ModifyTask(task)
		assert.True(t, modifiedTask.Completed)
	})

	t.Run("UncompleteFlag", func(t *testing.T) {
		// Reset the FlagSet and ModifierComposer
		mc = cli.NewModifierComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&mc.UncompleteFlag, "uncomplete", "Mark task as incomplete")

		// Parse the uncomplete flag
		err := fs.Parse([]string{"--uncomplete=true"})
		assert.NoError(t, err)

		// Verify modifier is created correctly
		modifier := mc.ComposeModifier()
		assert.NotNil(t, modifier)
		
		// Create a completed task
		task := core.Task{Completed: true}
		
		// The modifier should mark the task as incomplete
		modifiedTask := modifier.ModifyTask(task)
		assert.False(t, modifiedTask.Completed)
	})

	t.Run("PriorityFlag", func(t *testing.T) {
		// Reset the FlagSet and ModifierComposer
		mc = cli.NewModifierComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&mc.PriorityFlag, "priority", "Set priority")

		// Parse the priority flag
		err := fs.Parse([]string{"--priority=A"})
		assert.NoError(t, err)

		// Verify modifier is created correctly
		modifier := mc.ComposeModifier()
		assert.NotNil(t, modifier)
		
		// Create a task with no priority
		task := core.Task{Priority: ""}
		
		// The modifier should set the task priority to A
		modifiedTask := modifier.ModifyTask(task)
		assert.Equal(t, "A", modifiedTask.Priority)
	})

	t.Run("AddProjectFlag", func(t *testing.T) {
		// Reset the FlagSet and ModifierComposer
		mc = cli.NewModifierComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&mc.AddProjectFlag, "add-project", "Add project")

		// Parse the add-project flag
		err := fs.Parse([]string{"--add-project=Home"})
		assert.NoError(t, err)

		// Check if any modifiers were added
		t.Logf("Number of modifiers: %d", len(mc.Modifiers.Modifiers))
		for i, mod := range mc.Modifiers.Modifiers {
			t.Logf("Modifier %d: %T, IsActive: %v", i, mod, mod.IsActive())
		}

		// Verify modifier is created correctly
		modifier := mc.ComposeModifier()
		if !assert.NotNil(t, modifier) {
			t.FailNow()
		}
		
		// Create a task with no projects
		task := core.Task{Projects: []string{}}
		
		// The modifier should add the Home project
		modifiedTask := modifier.ModifyTask(task)
		assert.Contains(t, modifiedTask.Projects, "Home")
	})

	t.Run("RemoveProjectFlag", func(t *testing.T) {
		// Reset the FlagSet and ModifierComposer
		mc = cli.NewModifierComposer()
		fs = flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Var(&mc.RemoveProjectFlag, "remove-project", "Remove project")

		// Parse the remove-project flag
		err := fs.Parse([]string{"--remove-project=Home"})
		assert.NoError(t, err)

		// Check if any modifiers were added
		t.Logf("Number of modifiers: %d", len(mc.Modifiers.Modifiers))
		for i, mod := range mc.Modifiers.Modifiers {
			t.Logf("Modifier %d: %T, IsActive: %v", i, mod, mod.IsActive())
		}

		// Verify modifier is created correctly
		modifier := mc.ComposeModifier()
		if !assert.NotNil(t, modifier) {
			t.FailNow()
		}
		
		// Create a task with Home project
		task := core.Task{Projects: []string{"Home", "Work"}}
		
		// The modifier should remove the Home project
		modifiedTask := modifier.ModifyTask(task)
		assert.NotContains(t, modifiedTask.Projects, "Home")
		assert.Contains(t, modifiedTask.Projects, "Work")
	})
}