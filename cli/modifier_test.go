package cli

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"t5.mkbrechtel.dev/t5/core"
)

func TestModifierComposer(t *testing.T) {
	// Test creating a new modifier composer
	mc := NewModifierComposer()
	assert.NotNil(t, mc)
	assert.Empty(t, mc.Modifiers.Modifiers)
	
	// Test adding a modifier
	priorityModifier := &core.PriorityModifier{Priority: "A"}
	mc.AddModifier(priorityModifier)
	assert.Len(t, mc.Modifiers.Modifiers, 1)
	
	// Test composing a single modifier
	composed := mc.ComposeModifier()
	assert.Equal(t, priorityModifier, composed)
	
	// Test composing multiple modifiers
	projectModifier := &core.ProjectModifier{AddProjects: []string{"Home"}}
	mc.AddModifier(projectModifier)
	composed = mc.ComposeModifier()
	assert.IsType(t, &core.ChainModifier{}, composed)
	
	// Test handling inactive modifiers
	mc = NewModifierComposer()
	mc.AddModifier(&core.ProjectModifier{}) // Empty project modifier is inactive
	assert.Nil(t, mc.ComposeModifier())
}

func TestCompleteFlag(t *testing.T) {
	// Test the CompleteFlag
	mc := NewModifierComposer()
	completeFlag := &mc.CompleteFlag
	
	// Test String method
	assert.Equal(t, "false", completeFlag.String())
	
	// Test Set method
	err := completeFlag.Set("true")
	assert.NoError(t, err)
	assert.True(t, completeFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Completed: false}
	modified := composed.ModifyTask(task)
	assert.True(t, modified.Completed)
	assert.False(t, modified.CompletedDate.IsZero())
}

func TestUncompleteFlag(t *testing.T) {
	// Test the UncompleteFlag
	mc := NewModifierComposer()
	uncompleteFlag := &mc.UncompleteFlag
	
	// Test Set method
	err := uncompleteFlag.Set("true")
	assert.NoError(t, err)
	assert.True(t, uncompleteFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	completedDate := time.Now()
	task := core.Task{Completed: true, CompletedDate: completedDate}
	modified := composed.ModifyTask(task)
	assert.False(t, modified.Completed)
	assert.True(t, modified.CompletedDate.IsZero())
}

func TestPriorityModifierFlag(t *testing.T) {
	// Test the PriorityModifierFlag
	mc := NewModifierComposer()
	priorityFlag := &mc.PriorityFlag
	
	// Test String method
	assert.Equal(t, "", priorityFlag.String())
	
	// Test Set method
	err := priorityFlag.Set("A")
	assert.NoError(t, err)
	assert.Equal(t, "A", priorityFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Priority: "B"}
	modified := composed.ModifyTask(task)
	assert.Equal(t, "A", modified.Priority)
}

func TestRemovePriorityFlag(t *testing.T) {
	// Test the RemovePriorityFlag
	mc := NewModifierComposer()
	removePriorityFlag := &mc.RemovePriorityFlag
	
	// Test Set method
	err := removePriorityFlag.Set("true")
	assert.NoError(t, err)
	assert.True(t, removePriorityFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Priority: "A"}
	modified := composed.ModifyTask(task)
	assert.Equal(t, "", modified.Priority)
}

func TestAddProjectFlag(t *testing.T) {
	// Test the AddProjectFlag
	mc := NewModifierComposer()
	addProjectFlag := &mc.AddProjectFlag
	
	// Test String method
	assert.Equal(t, "", addProjectFlag.String())
	
	// Test Set method
	err := addProjectFlag.Set("Home,Work")
	assert.NoError(t, err)
	assert.Equal(t, "Home,Work", addProjectFlag.value)
	
	// Verify a modifier was added and is active
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Projects: []string{"Personal"}}
	modified := composed.ModifyTask(task)
	assert.Len(t, modified.Projects, 3)
	assert.Contains(t, modified.Projects, "Home")
	assert.Contains(t, modified.Projects, "Work")
	assert.Contains(t, modified.Projects, "Personal")
	
	// Test with + prefix in project name
	mc = NewModifierComposer()
	addProjectFlag = &mc.AddProjectFlag
	
	err = addProjectFlag.Set("+School")
	assert.NoError(t, err)
	
	composed = mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	task = core.Task{Projects: []string{"Home"}}
	modified = composed.ModifyTask(task)
	assert.Len(t, modified.Projects, 2)
	assert.Contains(t, modified.Projects, "Home")
	assert.Contains(t, modified.Projects, "School")
}

func TestRemoveProjectFlag(t *testing.T) {
	// Test the RemoveProjectFlag
	mc := NewModifierComposer()
	removeProjectFlag := &mc.RemoveProjectFlag
	
	// Test Set method
	err := removeProjectFlag.Set("Home,Personal")
	assert.NoError(t, err)
	assert.Equal(t, "Home,Personal", removeProjectFlag.value)
	
	// Verify a modifier was added and is active
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Projects: []string{"Home", "Work", "Personal"}}
	modified := composed.ModifyTask(task)
	assert.Len(t, modified.Projects, 1)
	assert.Contains(t, modified.Projects, "Work")
	assert.NotContains(t, modified.Projects, "Home")
	assert.NotContains(t, modified.Projects, "Personal")
}

func TestAddContextFlag(t *testing.T) {
	// Test the AddContextFlag
	mc := NewModifierComposer()
	addContextFlag := &mc.AddContextFlag
	
	// Test String method
	assert.Equal(t, "", addContextFlag.String())
	
	// Test Set method
	err := addContextFlag.Set("phone,email")
	assert.NoError(t, err)
	assert.Equal(t, "phone,email", addContextFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Contexts: []string{"meeting"}}
	modified := composed.ModifyTask(task)
	assert.Len(t, modified.Contexts, 3)
	assert.Contains(t, modified.Contexts, "phone")
	assert.Contains(t, modified.Contexts, "email")
	assert.Contains(t, modified.Contexts, "meeting")
	
	// Test with @ prefix in context name
	mc = NewModifierComposer()
	addContextFlag = &mc.AddContextFlag
	
	err = addContextFlag.Set("@computer")
	assert.NoError(t, err)
	
	composed = mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	task = core.Task{Contexts: []string{"phone"}}
	modified = composed.ModifyTask(task)
	assert.Len(t, modified.Contexts, 2)
	assert.Contains(t, modified.Contexts, "phone")
	assert.Contains(t, modified.Contexts, "computer")
}

func TestRemoveContextFlag(t *testing.T) {
	// Test the RemoveContextFlag
	mc := NewModifierComposer()
	removeContextFlag := &mc.RemoveContextFlag
	
	// Test Set method
	err := removeContextFlag.Set("phone,meeting")
	assert.NoError(t, err)
	assert.Equal(t, "phone,meeting", removeContextFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Contexts: []string{"phone", "email", "meeting"}}
	modified := composed.ModifyTask(task)
	assert.Len(t, modified.Contexts, 1)
	assert.Contains(t, modified.Contexts, "email")
	assert.NotContains(t, modified.Contexts, "phone")
	assert.NotContains(t, modified.Contexts, "meeting")
}

func TestDueFlag(t *testing.T) {
	// Test the DueFlag
	mc := NewModifierComposer()
	dueFlag := &mc.DueFlag
	
	// Test String method
	assert.Equal(t, "", dueFlag.String())
	
	// Test Set method
	err := dueFlag.Set("2023-12-31")
	assert.NoError(t, err)
	assert.Equal(t, "2023-12-31", dueFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{}
	modified := composed.ModifyTask(task)
	assert.False(t, modified.DueDate.IsZero())
	
	// Test invalid date format
	mc = NewModifierComposer()
	dueFlag = &mc.DueFlag
	
	err = dueFlag.Set("invalid-date")
	assert.Error(t, err)
}

func TestRemoveDueFlag(t *testing.T) {
	// Test the RemoveDueFlag
	mc := NewModifierComposer()
	removeDueFlag := &mc.RemoveDueFlag
	
	// Test Set method
	err := removeDueFlag.Set("true")
	assert.NoError(t, err)
	assert.True(t, removeDueFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	dueDate := time.Now()
	task := core.Task{DueDate: dueDate}
	modified := composed.ModifyTask(task)
	assert.True(t, modified.DueDate.IsZero())
}

func TestAppendFlag(t *testing.T) {
	// Test the AppendFlag
	mc := NewModifierComposer()
	appendFlag := &mc.AppendFlag
	
	// Test String method
	assert.Equal(t, "", appendFlag.String())
	
	// Test Set method
	err := appendFlag.Set(" for dinner")
	assert.NoError(t, err)
	assert.Equal(t, " for dinner", appendFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Todo: "Buy groceries"}
	modified := composed.ModifyTask(task)
	assert.Equal(t, "Buy groceries for dinner", modified.Todo)
}

func TestPrependFlag(t *testing.T) {
	// Test the PrependFlag
	mc := NewModifierComposer()
	prependFlag := &mc.PrependFlag
	
	// Test String method
	assert.Equal(t, "", prependFlag.String())
	
	// Test Set method
	err := prependFlag.Set("Today: ")
	assert.NoError(t, err)
	assert.Equal(t, "Today: ", prependFlag.value)
	
	// Verify a modifier was added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the modifier's behavior
	task := core.Task{Todo: "Buy groceries"}
	modified := composed.ModifyTask(task)
	assert.Equal(t, "Today: Buy groceries", modified.Todo)
}

func TestCombiningTextModifiers(t *testing.T) {
	// Test that append and prepend flags share a modifier
	mc := NewModifierComposer()
	
	// Set both flags
	mc.AppendFlag.Set(" for dinner")
	mc.PrependFlag.Set("Today: ")
	
	// Verify a single modifier handles both
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test combined behavior
	task := core.Task{Todo: "Buy groceries"}
	modified := composed.ModifyTask(task)
	assert.Equal(t, "Today: Buy groceries for dinner", modified.Todo)
}

func TestModifierFlagIntegration(t *testing.T) {
	// Test integrating with the flag package
	mc := NewModifierComposer()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	
	// Add some modifier flags
	fs.Var(&mc.CompleteFlag, "complete", "Mark task as complete")
	fs.Var(&mc.PriorityFlag, "priority", "Set task priority")
	fs.Var(&mc.AddProjectFlag, "add-project", "Add project to task")
	
	// Parse some flags
	args := []string{"--complete=true", "--priority=A", "--add-project=Home"}
	err := fs.Parse(args)
	assert.NoError(t, err)
	
	// Verify flags were set correctly
	assert.True(t, mc.CompleteFlag.value)
	assert.Equal(t, "A", mc.PriorityFlag.value)
	assert.Equal(t, "Home", mc.AddProjectFlag.value)
	
	// Verify modifiers were added
	composed := mc.ComposeModifier()
	assert.NotNil(t, composed)
	
	// Test the combined modifier behavior
	task := core.Task{
		Completed: false,
		Priority:  "B",
		Projects:  []string{"Work"},
	}
	
	modified := composed.ModifyTask(task)
	assert.True(t, modified.Completed)
	assert.Equal(t, "A", modified.Priority)
	assert.Contains(t, modified.Projects, "Home")
	assert.Contains(t, modified.Projects, "Work")
}