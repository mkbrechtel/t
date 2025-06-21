package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"t5.mkbrechtel.dev/t5/core"
)

// FilterComposer combines multiple filters and provides flag.Value
// implementations for easy flag integration
type FilterComposer struct {
	Filters          core.AndFilter
	CompletedFlag    CompletionFlag
	NotCompletedFlag CompletionFlag
	PriorityFlag     PriorityFlag
	ProjectFlag      ProjectFlag
	ContextFlag      ContextFlag
	DueTodayFlag     DueTodayFlag
	DueBeforeFlag    DueBeforeFlag
	CreatedAfterFlag CreatedAfterFlag
	HasTagFlag       TagFlag
	RegexFlag        RegexFlag
}

// NewFilterComposer creates a new FilterComposer
func NewFilterComposer() *FilterComposer {
	fc := &FilterComposer{
		Filters: core.AndFilter{Filters: make([]core.TaskFilter, 0)},
	}

	// Initialize flags with parent reference
	fc.CompletedFlag.parent = fc
	fc.CompletedFlag.isCompleted = true

	fc.NotCompletedFlag.parent = fc
	fc.NotCompletedFlag.isCompleted = false

	fc.PriorityFlag.parent = fc
	fc.ProjectFlag.parent = fc
	fc.ContextFlag.parent = fc
	fc.DueTodayFlag.parent = fc
	fc.DueBeforeFlag.parent = fc
	fc.CreatedAfterFlag.parent = fc
	fc.HasTagFlag.parent = fc
	fc.RegexFlag.parent = fc

	return fc
}

// AddFilterFlags adds all filter flags to the provided flagset
func (fc *FilterComposer) AddFilterFlags(flagset *flag.FlagSet) {
	flagset.Var(&fc.CompletedFlag, "completed", "Show only completed tasks")
	flagset.Var(&fc.NotCompletedFlag, "not-completed", "Show only non-completed tasks")
	flagset.Var(&fc.PriorityFlag, "priority", "Filter by priority (e.g., A,B,C)")
	flagset.Var(&fc.ProjectFlag, "project", "Filter by project")
	flagset.Var(&fc.ContextFlag, "context", "Filter by context")
	flagset.Var(&fc.DueTodayFlag, "due-today", "Show tasks due today")
	flagset.Var(&fc.DueBeforeFlag, "due-before", "Show tasks due before a date (YYYY-MM-DD)")
	flagset.Var(&fc.CreatedAfterFlag, "created-after", "Show tasks created after a date (YYYY-MM-DD)")
	flagset.Var(&fc.HasTagFlag, "has-tag", "Show tasks with a specific tag")
	flagset.Var(&fc.RegexFlag, "regex", "Filter tasks using a regular expression")
}

// ProcessOriginalArgs has been removed as it's no longer needed.
// Flag parsing is now handled directly by each command with its own flag set.

// AddFilter adds a filter to the composer
func (fc *FilterComposer) AddFilter(filter core.TaskFilter) {
	if filter != nil {
		fc.Filters.Filters = append(fc.Filters.Filters, filter)
	}
}

// ComposeFilter creates a filter for use with task filtering
func (fc *FilterComposer) ComposeFilter() core.TaskFilter {
	// Filter out inactive filters
	activeFilters := make([]core.TaskFilter, 0, len(fc.Filters.Filters))
	for _, filter := range fc.Filters.Filters {
		if filter.IsActive() {
			activeFilters = append(activeFilters, filter)
		}
	}

	if len(activeFilters) == 0 {
		return nil
	}

	if len(activeFilters) == 1 {
		return activeFilters[0]
	}

	// Create a new AndFilter with only active filters
	return &core.AndFilter{Filters: activeFilters}
}

// CompletionFlag implements flag.Value for completion status filtering
type CompletionFlag struct {
	parent      *FilterComposer
	isCompleted bool
	value       bool
}

func (f *CompletionFlag) String() string {
	return fmt.Sprintf("%v", f.value)
}

func (f *CompletionFlag) Set(value string) error {
	// For boolean flags, presence of the flag is enough
	f.value = true

	// Create and add filter
	var completed *bool = &f.isCompleted
	filter := &core.CompletionFilter{
		Completed: completed,
		UseFlag:   f.isCompleted,
		NotFlag:   !f.isCompleted,
	}
	f.parent.AddFilter(filter)

	return nil
}

// PriorityFlag implements flag.Value for priority filtering
type PriorityFlag struct {
	parent *FilterComposer
	value  string
}

func (f *PriorityFlag) String() string {
	return f.value
}

func (f *PriorityFlag) Set(value string) error {
	f.value = value

	if value == "" {
		return nil
	}

	// Handle comma-separated list of priorities
	priorities := strings.Split(value, ",")
	priorityFilters := make([]core.TaskFilter, 0, len(priorities))

	for _, p := range priorities {
		p = strings.TrimSpace(p)
		if p != "" {
			priorityFilters = append(priorityFilters, &core.PriorityFilter{Priorities: []string{strings.ToUpper(p)}})
		}
	}

	// If we have multiple priorities, use OR filter to match any of them
	if len(priorityFilters) > 1 {
		f.parent.AddFilter(&core.OrFilter{Filters: priorityFilters})
	} else if len(priorityFilters) == 1 {
		f.parent.AddFilter(priorityFilters[0])
	}

	return nil
}

// ProjectFlag implements flag.Value for project filtering
type ProjectFlag struct {
	parent *FilterComposer
	value  string
}

func (f *ProjectFlag) String() string {
	return f.value
}

func (f *ProjectFlag) Set(value string) error {
	f.value = value

	if value == "" {
		return nil
	}

	projects := strings.Split(value, ",")
	for i, project := range projects {
		project = strings.TrimSpace(project)
		// Remove leading + if present
		project = strings.TrimPrefix(project, "+")
		projects[i] = project
	}

	filter := &core.ProjectFilter{Projects: projects}
	f.parent.AddFilter(filter)

	return nil
}

// ContextFlag implements flag.Value for context filtering
type ContextFlag struct {
	parent *FilterComposer
	value  string
}

func (f *ContextFlag) String() string {
	return f.value
}

func (f *ContextFlag) Set(value string) error {
	f.value = value

	if value == "" {
		return nil
	}

	contexts := strings.Split(value, ",")
	for i, context := range contexts {
		context = strings.TrimSpace(context)
		// Remove leading @ if present
		context = strings.TrimPrefix(context, "@")
		contexts[i] = context
	}

	filter := &core.ContextFilter{Contexts: contexts}
	f.parent.AddFilter(filter)

	return nil
}

// DueTodayFlag implements flag.Value for due-today filtering
type DueTodayFlag struct {
	parent *FilterComposer
	value  bool
}

func (f *DueTodayFlag) String() string {
	return fmt.Sprintf("%v", f.value)
}

func (f *DueTodayFlag) Set(value string) error {
	// For boolean flags, presence of the flag is enough
	f.value = true

	// Create due today filter
	today := time.Now()
	tomorrow := time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, today.Location())

	filter := &core.DueDateFilter{
		Before: &tomorrow,
		After:  &today,
	}
	f.parent.AddFilter(filter)

	return nil
}

// DueBeforeFlag implements flag.Value for due-before filtering
type DueBeforeFlag struct {
	parent *FilterComposer
	value  string
}

func (f *DueBeforeFlag) String() string {
	return f.value
}

func (f *DueBeforeFlag) Set(value string) error {
	f.value = value

	if value == "" {
		return nil
	}

	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	filter := &core.DueDateFilter{Before: &date}
	f.parent.AddFilter(filter)

	return nil
}

// CreatedAfterFlag implements flag.Value for created-after filtering
type CreatedAfterFlag struct {
	parent *FilterComposer
	value  string
}

func (f *CreatedAfterFlag) String() string {
	return f.value
}

func (f *CreatedAfterFlag) Set(value string) error {
	f.value = value

	if value == "" {
		return nil
	}

	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	filter := &core.CreatedDateFilter{After: &date}
	f.parent.AddFilter(filter)

	return nil
}

// TagFlag implements flag.Value for tag filtering
type TagFlag struct {
	parent *FilterComposer
	value  string
}

func (f *TagFlag) String() string {
	return f.value
}

func (f *TagFlag) Set(value string) error {
	f.value = value

	if value == "" {
		return nil
	}

	filter := &core.TagFilter{Tag: value}
	f.parent.AddFilter(filter)

	return nil
}

// RegexFlag implements flag.Value for regex filtering
type RegexFlag struct {
	parent *FilterComposer
	value  string
}

func (f *RegexFlag) String() string {
	return f.value
}

func (f *RegexFlag) Set(value string) error {
	f.value = value

	if value == "" {
		return nil
	}

	filter := &core.RegexFilter{Pattern: value}
	f.parent.AddFilter(filter)

	return nil
}
