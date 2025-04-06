package sync

import (
	"fmt"
	"strings"
	"time"
	
	"t5.mkbrechtel.dev/t5/core"
)

// BuildFilterFromConfig creates a task filter from a configuration map
func BuildFilterFromConfig(filterConfig map[string]interface{}) (core.TaskFilter, error) {
	if len(filterConfig) == 0 {
		return nil, nil
	}
	
	composer := core.NewFilterComposer()
	
	// Process each filter type
	for key, value := range filterConfig {
		switch strings.ToLower(key) {
		case "completed":
			if boolVal, ok := value.(bool); ok && boolVal {
				completionFilter := &core.CompletionFilter{UseFlag: true}
				composer.AddFilter(completionFilter)
			}
		case "not_completed", "not-completed":
			if boolVal, ok := value.(bool); ok && boolVal {
				completionFilter := &core.CompletionFilter{NotFlag: true}
				composer.AddFilter(completionFilter)
			}
		case "priority":
			if prioritiesVal, ok := value.(string); ok {
				priorityFilter := &core.PriorityFilter{}
				priorityFilter.Set(prioritiesVal)
				composer.AddFilter(priorityFilter)
			}
		case "project":
			if projectsVal, ok := value.(string); ok {
				projectFilter := &core.ProjectFilter{}
				projectFilter.Set(projectsVal)
				composer.AddFilter(projectFilter)
			}
		case "context":
			if contextsVal, ok := value.(string); ok {
				contextFilter := &core.ContextFilter{}
				contextFilter.Set(contextsVal)
				composer.AddFilter(contextFilter)
			}
		case "due_today", "due-today":
			if boolVal, ok := value.(bool); ok && boolVal {
				// Set both before and after to today's date
				today := time.Now()
				tomorrow := today.Add(24 * time.Hour)
				
				dueDateFilter := &core.DueDateFilter{
					After:  &today,
					Before: &tomorrow,
					Today:  true,
				}
				composer.AddFilter(dueDateFilter)
			}
		case "due_before", "due-before":
			if dateStr, ok := value.(string); ok {
				dueDate, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					return nil, fmt.Errorf("invalid date format for due_before: %v", err)
				}
				
				dueDateFilter := &core.DueDateFilter{
					Before: &dueDate,
				}
				composer.AddFilter(dueDateFilter)
			}
		case "created_after", "created-after":
			if dateStr, ok := value.(string); ok {
				createdDate, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					return nil, fmt.Errorf("invalid date format for created_after: %v", err)
				}
				
				createdFilter := &core.CreatedDateFilter{
					After: &createdDate,
				}
				composer.AddFilter(createdFilter)
			}
		case "has_tag", "has-tag":
			if tagVal, ok := value.(string); ok {
				tagFilter := &core.TagFilter{Tag: tagVal}
				composer.AddFilter(tagFilter)
			}
		case "regex":
			if regexVal, ok := value.(string); ok {
				regexFilter := &core.RegexFilter{}
				regexFilter.Set(regexVal)
				composer.AddFilter(regexFilter)
			}
		}
	}
	
	return composer.ComposeFilter(), nil
}

// BuildModifierFromConfig creates a task modifier from a configuration map
func BuildModifierFromConfig(modifierConfig map[string]interface{}) (core.TaskModifier, error) {
	if len(modifierConfig) == 0 {
		return nil, nil
	}
	
	composer := core.NewModifierComposer()
	
	// Process each modifier type
	for key, value := range modifierConfig {
		switch strings.ToLower(key) {
		case "complete":
			if boolVal, ok := value.(bool); ok && boolVal {
				completionModifier := &core.CompletionModifier{Complete: true}
				composer.AddModifier(completionModifier)
			}
		case "uncomplete":
			if boolVal, ok := value.(bool); ok && boolVal {
				completionModifier := &core.CompletionModifier{Uncomplete: true}
				composer.AddModifier(completionModifier)
			}
		case "priority":
			if priorityVal, ok := value.(string); ok {
				priorityModifier := &core.PriorityModifier{}
				priorityModifier.Set(priorityVal)
				composer.AddModifier(priorityModifier)
			}
		case "remove_priority", "remove-priority":
			if boolVal, ok := value.(bool); ok && boolVal {
				priorityModifier := &core.PriorityModifier{RemovePriority: true}
				composer.AddModifier(priorityModifier)
			}
		case "add_project", "add-project":
			if projectVal, ok := value.(string); ok {
				projectModifier := core.ProjectModifier{}
				addFlag := core.AddProjectFlag{ProjectModifier: &projectModifier}
				addFlag.Set(projectVal)
				composer.AddModifier(&projectModifier)
			}
		case "remove_project", "remove-project":
			if projectVal, ok := value.(string); ok {
				projectModifier := core.ProjectModifier{}
				removeFlag := core.RemoveProjectFlag{ProjectModifier: &projectModifier}
				removeFlag.Set(projectVal)
				composer.AddModifier(&projectModifier)
			}
		case "add_context", "add-context":
			if contextVal, ok := value.(string); ok {
				contextModifier := core.ContextModifier{}
				addFlag := core.AddContextFlag{ContextModifier: &contextModifier}
				addFlag.Set(contextVal)
				composer.AddModifier(&contextModifier)
			}
		case "remove_context", "remove-context":
			if contextVal, ok := value.(string); ok {
				contextModifier := core.ContextModifier{}
				removeFlag := core.RemoveContextFlag{ContextModifier: &contextModifier}
				removeFlag.Set(contextVal)
				composer.AddModifier(&contextModifier)
			}
		case "due":
			if dueVal, ok := value.(string); ok {
				dueModifier := &core.DueDateModifier{}
				dueModifier.Set(dueVal)
				composer.AddModifier(dueModifier)
			}
		case "remove_due", "remove-due":
			if boolVal, ok := value.(bool); ok && boolVal {
				dueModifier := &core.DueDateModifier{RemoveDueDate: true}
				composer.AddModifier(dueModifier)
			}
		case "append":
			if appendVal, ok := value.(string); ok {
				textModifier := &core.TextModifier{AppendText: appendVal}
				composer.AddModifier(textModifier)
			}
		case "prepend":
			if prependVal, ok := value.(string); ok {
				textModifier := &core.TextModifier{PrependText: prependVal}
				composer.AddModifier(textModifier)
			}
		}
	}
	
	return composer.ComposeModifier(), nil
}