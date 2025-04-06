package core

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type TaskFilter interface {
	FilterTask(Task) bool
	// IsActive returns true if the filter has been initialized with values
	IsActive() bool
}

// ProjectFilter filters tasks that have specific projects
type ProjectFilter struct {
	Projects []string
}

func (f ProjectFilter) FilterTask(task Task) bool {
	if len(f.Projects) == 0 {
		return true
	}
	for _, project := range f.Projects {
		for _, taskProject := range task.Projects {
			if project == taskProject {
				return true
			}
		}
	}
	return false
}

func (f ProjectFilter) IsActive() bool {
	return len(f.Projects) > 0
}

// Implement flag.Value interface
func (f *ProjectFilter) String() string {
	return strings.Join(f.Projects, ",")
}

func (f *ProjectFilter) Set(value string) error {
	if value == "" {
		return nil
	}
	
	projects := strings.Split(value, ",")
	for i, project := range projects {
		project = strings.TrimSpace(project)
		// Remove leading + if present
		if strings.HasPrefix(project, "+") {
			project = project[1:]
		}
		projects[i] = project
	}
	f.Projects = projects
	return nil
}

// ContextFilter filters tasks that have specific contexts
type ContextFilter struct {
	Contexts []string
}

func (f ContextFilter) FilterTask(task Task) bool {
	if len(f.Contexts) == 0 {
		return true
	}
	for _, context := range f.Contexts {
		for _, taskContext := range task.Contexts {
			if context == taskContext {
				return true
			}
		}
	}
	return false
}

func (f ContextFilter) IsActive() bool {
	return len(f.Contexts) > 0
}

// Implement flag.Value interface
func (f *ContextFilter) String() string {
	return strings.Join(f.Contexts, ",")
}

func (f *ContextFilter) Set(value string) error {
	if value == "" {
		return nil
	}
	
	contexts := strings.Split(value, ",")
	for i, context := range contexts {
		context = strings.TrimSpace(context)
		// Remove leading @ if present
		if strings.HasPrefix(context, "@") {
			context = context[1:]
		}
		contexts[i] = context
	}
	f.Contexts = contexts
	return nil
}

// CompletionFilter filters tasks based on completion status
type CompletionFilter struct {
	Completed *bool
	UseFlag   bool // Set when --completed is used
	NotFlag   bool // Set when --not-completed is used
}

func (f CompletionFilter) FilterTask(task Task) bool {
	// If UseFlag is set, only show completed tasks
	if f.UseFlag {
		return task.Completed
	}
	
	// If NotFlag is set, only show non-completed tasks
	if f.NotFlag {
		return !task.Completed
	}
	
	// If neither flag is set, show all tasks
	return true
}

func (f CompletionFilter) IsActive() bool {
	// The filter is active if either the completed or not-completed flag is set
	return f.UseFlag || f.NotFlag
}

// String implements part of the flag.Value interface
func (f *CompletionFilter) String() string {
	if f.UseFlag {
		return "completed:true"
	}
	if f.NotFlag {
		return "not-completed:true"
	}
	return ""
}

// Set implements part of the flag.Value interface
func (f *CompletionFilter) Set(value string) error {
	// This won't be called directly since we're using BoolVar
	return nil
}

// PriorityFilter filters tasks by priority
type PriorityFilter struct {
	Priorities []string
}

func (f PriorityFilter) FilterTask(task Task) bool {
	if len(f.Priorities) == 0 {
		return true
	}
	
	for _, priority := range f.Priorities {
		if task.Priority == priority {
			return true
		}
	}
	return false
}

func (f PriorityFilter) IsActive() bool {
	return len(f.Priorities) > 0
}

// Implement flag.Value interface
func (f *PriorityFilter) String() string {
	return strings.Join(f.Priorities, ",")
}

func (f *PriorityFilter) Set(value string) error {
	if value == "" {
		return nil
	}
	
	priorities := strings.Split(value, ",")
	for i, p := range priorities {
		p = strings.TrimSpace(p)
		if p != "" {
			priorities[i] = strings.ToUpper(p)
		}
	}
	f.Priorities = priorities
	return nil
}

// DueDateFilter filters tasks by due date range
type DueDateFilter struct {
	Before *time.Time
	After  *time.Time
	Today  bool
}

func (f DueDateFilter) FilterTask(task Task) bool {
	if f.Before == nil && f.After == nil {
		return true
	}
	if task.DueDate.IsZero() {
		return false
	}
	if f.Before != nil && task.DueDate.After(*f.Before) {
		return false
	}
	if f.After != nil && !task.DueDate.After(*f.After) {
		return false
	}
	return true
}

func (f DueDateFilter) IsActive() bool {
	return f.Before != nil || f.After != nil || f.Today
}

// Implement flag.Value interface
func (f *DueDateFilter) String() string {
	if f.Before != nil {
		return f.Before.Format("2006-01-02")
	}
	return ""
}

func (f *DueDateFilter) Set(value string) error {
	if value == "" {
		return nil
	}
	
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	f.Before = &date
	return nil
}

// CreatedDateFilter filters tasks by creation date range
type CreatedDateFilter struct {
	Before *time.Time
	After  *time.Time
}

func (f CreatedDateFilter) FilterTask(task Task) bool {
	if f.Before == nil && f.After == nil {
		return true
	}
	if task.CreatedDate.IsZero() {
		return false
	}
	if f.Before != nil && task.CreatedDate.After(*f.Before) {
		return false
	}
	if f.After != nil && !task.CreatedDate.After(*f.After) {
		return false
	}
	return true
}

func (f CreatedDateFilter) IsActive() bool {
	return f.Before != nil || f.After != nil
}

// Implement flag.Value interface
func (f *CreatedDateFilter) String() string {
	if f.After != nil {
		return f.After.Format("2006-01-02")
	}
	return ""
}

func (f *CreatedDateFilter) Set(value string) error {
	if value == "" {
		return nil
	}
	
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	f.After = &date
	return nil
}

// TagFilter filters tasks by the presence of a specific tag
type TagFilter struct {
	Tag string
}

func (f TagFilter) FilterTask(task Task) bool {
	if f.Tag == "" {
		return true
	}
	
	// First check if it's a key-only tag
	if _, exists := task.AdditionalTags[f.Tag]; exists {
		return true
	}
	
	// Then check for any tag key that starts with the given tag
	for key := range task.AdditionalTags {
		if strings.HasPrefix(key, f.Tag) {
			return true
		}
	}
	
	return false
}

func (f TagFilter) IsActive() bool {
	return f.Tag != ""
}

// Implement flag.Value interface
func (f *TagFilter) String() string {
	return f.Tag
}

func (f *TagFilter) Set(value string) error {
	f.Tag = value
	return nil
}

// RegexFilter filters tasks matching a regular expression pattern
type RegexFilter struct {
	Pattern string
	regex   *regexp.Regexp
}

func (f RegexFilter) FilterTask(task Task) bool {
	if f.Pattern == "" {
		return true
	}
	
	// Compile the regex pattern
	pattern, err := regexp.Compile(f.Pattern)
	if err != nil {
		// If pattern is invalid, log an error and return false
		log.Printf("Invalid regex pattern: %v", err)
		return false
	}
	
	// Check task text
	if pattern.MatchString(task.Todo) {
		return true
	}
	
	// Check projects
	for _, project := range task.Projects {
		if pattern.MatchString(project) {
			return true
		}
	}
	
	// Check contexts
	for _, context := range task.Contexts {
		if pattern.MatchString(context) {
			return true
		}
	}
	
	// Check tag keys and values
	for key, value := range task.AdditionalTags {
		if pattern.MatchString(key) || pattern.MatchString(value) {
			return true
		}
	}
	
	return false
}

func (f RegexFilter) IsActive() bool {
	return f.Pattern != ""
}

// Implement flag.Value interface
func (f *RegexFilter) String() string {
	return f.Pattern
}

func (f *RegexFilter) Set(value string) error {
	if value == "" {
		return nil
	}
	
	// Validate the regex pattern
	_, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %v", err)
	}
	
	f.Pattern = value
	return nil
}

// AndFilter combines multiple filters with AND logic
type AndFilter struct {
	Filters []TaskFilter
}

func (f AndFilter) FilterTask(task Task) bool {
	if len(f.Filters) == 0 {
		return true
	}
	
	for _, filter := range f.Filters {
		if !filter.FilterTask(task) {
			return false
		}
	}
	
	return true
}

func (f AndFilter) IsActive() bool {
	return len(f.Filters) > 0
}

// AddFilter adds a filter to the AndFilter if it's active
func (f *AndFilter) AddFilter(filter TaskFilter) {
	if filter != nil && filter.IsActive() {
		f.Filters = append(f.Filters, filter)
	}
}

// OrFilter combines multiple filters with OR logic
type OrFilter struct {
	Filters []TaskFilter
}

func (f OrFilter) FilterTask(task Task) bool {
	if len(f.Filters) == 0 {
		return false // Empty OR returns false (nothing to match)
	}
	
	for _, filter := range f.Filters {
		if filter.FilterTask(task) {
			return true
		}
	}
	
	return false
}

func (f OrFilter) IsActive() bool {
	return len(f.Filters) > 0
}

// AddFilter adds a filter to the OrFilter if it's active
func (f *OrFilter) AddFilter(filter TaskFilter) {
	if filter != nil && filter.IsActive() {
		f.Filters = append(f.Filters, filter)
	}
}

// NotFilter negates the result of a filter
type NotFilter struct {
	Filter TaskFilter
}

func (f NotFilter) FilterTask(task Task) bool {
	return !f.Filter.FilterTask(task)
}

func (f NotFilter) IsActive() bool {
	if f.Filter == nil {
		return false
	}
	return f.Filter.IsActive()
}

// FilterComposer manages a collection of filters and builds an AndFilter
// It integrates with flag parsing to only include filters that are actually selected
type FilterComposer struct {
	filters []TaskFilter
}

// NewFilterComposer creates a new FilterComposer
func NewFilterComposer() *FilterComposer {
	return &FilterComposer{filters: make([]TaskFilter, 0)}
}

// AddFilter adds a filter to the composer
// Only active filters will be included in the final composition
func (fc *FilterComposer) AddFilter(filter TaskFilter) {
	if filter != nil {
		fc.filters = append(fc.filters, filter)
	}
}

// ComposeFilter creates an AndFilter containing all active filters
// If no filters are active, it returns nil
// If only one filter is active, it returns that filter directly
func (fc *FilterComposer) ComposeFilter() TaskFilter {
	activeFilters := make([]TaskFilter, 0)
	
	for _, filter := range fc.filters {
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
	
	return &AndFilter{Filters: activeFilters}
}
