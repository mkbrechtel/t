package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/adrg/xdg"
	yaml "gopkg.in/yaml.v2"
	"t5.mkbrechtel.dev/core"
)

// AppConfig holds all configuration settings
type AppConfig struct {
	TodoFile             string
	ConfigFile           string
	EventStoreFile       string
	PreferShortIds       bool
	EnforceCreationDate  bool
	EnforceCompletionDate bool
}

// DefaultAppConfig returns a configuration with default values
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		TodoFile:             "todo.txt",
		ConfigFile:           "",
		EventStoreFile:       "",
		PreferShortIds:       true,
		EnforceCreationDate:  true,
		EnforceCompletionDate: true,
	}
}

// AppContext holds application state and configuration
type AppContext struct {
	Config *AppConfig
	Repository *core.Repository
	FlagSet *flag.FlagSet
}

// NewAppContext creates a new application context with default config
func NewAppContext() *AppContext {
	return &AppContext{
		Config: DefaultAppConfig(),
		FlagSet: flag.CommandLine,
	}
}

// SetupFlags initializes command line flags with config values
func (ctx *AppContext) SetupFlags() {
	if ctx.FlagSet == nil {
		// If no FlagSet is provided, use the global flag package
		ctx.FlagSet = flag.CommandLine
	}
	
	// Create flag variables
	ctx.FlagSet.StringVar(&ctx.Config.TodoFile, "todo", ctx.Config.TodoFile, "todo.txt file path")
	ctx.FlagSet.StringVar(&ctx.Config.TodoFile, "t", ctx.Config.TodoFile, "todo.txt file path (shorthand)")
	ctx.FlagSet.StringVar(&ctx.Config.ConfigFile, "config", ctx.Config.ConfigFile, "config file path")
	ctx.FlagSet.StringVar(&ctx.Config.ConfigFile, "c", ctx.Config.ConfigFile, "config file path (shorthand)")
	ctx.FlagSet.StringVar(&ctx.Config.EventStoreFile, "eventstore", ctx.Config.EventStoreFile, "Event store file path (if not specified, uses in-memory storage)")

	// Define flags for task properties
	ctx.FlagSet.BoolVar(&ctx.Config.PreferShortIds, "prefer-short-ids", ctx.Config.PreferShortIds, "Use short form IDs instead of UUIDs")
	ctx.FlagSet.BoolVar(&ctx.Config.EnforceCreationDate, "enforce-creation-date", ctx.Config.EnforceCreationDate, "Ensure tasks have creation dates")
	ctx.FlagSet.BoolVar(&ctx.Config.EnforceCompletionDate, "enforce-completion-date", ctx.Config.EnforceCompletionDate, "Ensure completed tasks have completion dates")
}

// TodoFileConfig represents the configuration for a single todo.txt file
type TodoFileConfig struct {
	Path   string          `yaml:"path"`
	Ensure map[string]bool `yaml:"ensure,omitempty"`
}

// ShowConfig displays the current configuration
func (ctx *AppContext) ShowConfig() {
	fmt.Println("Current Configuration:")
	fmt.Printf("Todo File: %s\n", ctx.Config.TodoFile)
	if ctx.Config.ConfigFile != "" {
		fmt.Printf("Config File: %s\n", ctx.Config.ConfigFile)
	}
	if ctx.Config.EventStoreFile != "" {
		fmt.Printf("Event Store: %s\n", ctx.Config.EventStoreFile)
	} else {
		fmt.Println("Event Store: in-memory (changes will not be persisted)")
	}
	
	fmt.Println("\nTask Properties:")
	fmt.Printf("  Prefer Short IDs: %v\n", ctx.Config.PreferShortIds)
	fmt.Printf("  Enforce Creation Date: %v\n", ctx.Config.EnforceCreationDate)
	fmt.Printf("  Enforce Completion Date: %v\n", ctx.Config.EnforceCompletionDate)
}

// GetTodoFiles retrieves the configured todo.txt files
func (ctx *AppContext) GetTodoFiles() []TodoFileConfig {
	todoFile := TodoFileConfig{
		Path: ctx.Config.TodoFile,
		Ensure: map[string]bool{
			"prefershortids":       ctx.Config.PreferShortIds,
			"enforcecompletiondate": ctx.Config.EnforceCompletionDate,
			"enforcecreationdate":   ctx.Config.EnforceCreationDate,
		},
	}
	return []TodoFileConfig{todoFile}
}

// IsFlagPassed checks if a flag was explicitly passed on the command line
func (ctx *AppContext) IsFlagPassed(name string) bool {
	found := false
	ctx.FlagSet.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// LoadConfigFile loads the YAML config file if specified
func (ctx *AppContext) LoadConfigFile() error {
	cfg := ctx.Config
	if cfg.ConfigFile == "" {
		// Try default locations
		homeConfig := filepath.Join(xdg.ConfigHome, "t5", "config.yaml")
		if _, err := os.Stat(homeConfig); err == nil {
			cfg.ConfigFile = homeConfig
		} else {
			// No config file found, use defaults
			return nil
		}
	}

	data, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML config
	var config struct {
		Todo struct {
			File   string `yaml:"file"`
			Ensure struct {
				PreferShortIds        bool `yaml:"prefershortids"`
				EnforceCreationDate   bool `yaml:"enforcecreationdate"`
				EnforceCompletionDate bool `yaml:"enforcecompletiondate"`
			} `yaml:"ensure"`
		} `yaml:"todo"`
		EventStore struct {
			File string `yaml:"file"`
		} `yaml:"eventstore"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Only override command-line values if they're specified in config and flags weren't explicitly set
	if config.Todo.File != "" && !ctx.IsFlagPassed("todo") && !ctx.IsFlagPassed("t") {
		cfg.TodoFile = config.Todo.File
	}
	
	if config.EventStore.File != "" && !ctx.IsFlagPassed("eventstore") {
		cfg.EventStoreFile = config.EventStore.File
	}

	// Apply ensure settings if they're present and flags weren't explicitly set
	if !ctx.IsFlagPassed("prefer-short-ids") {
		cfg.PreferShortIds = config.Todo.Ensure.PreferShortIds
	}
	
	if !ctx.IsFlagPassed("enforce-creation-date") {
		cfg.EnforceCreationDate = config.Todo.Ensure.EnforceCreationDate
	}
	
	if !ctx.IsFlagPassed("enforce-completion-date") {
		cfg.EnforceCompletionDate = config.Todo.Ensure.EnforceCompletionDate
	}

	return nil
}

// InitRepository initializes a repository with the appropriate event store
func (ctx *AppContext) InitRepository() error {
	cfg := ctx.Config
	
	// Check if an event store file is specified
	eventStoreFile := cfg.EventStoreFile
	if eventStoreFile == "" {
		// Use the default location in the data home directory
		eventStoreFile = filepath.Join(xdg.DataHome, "t5", "events.jsonl")
	}

	// If the path is not absolute, make it relative to the working directory or data home
	if !filepath.IsAbs(eventStoreFile) {
		if len(eventStoreFile) >= 2 && eventStoreFile[0:2] == "./" {
			// Keep as relative path for testing
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}
			eventStoreFile = filepath.Join(wd, eventStoreFile[2:])
		} else {
			// Otherwise use the data home directory
			eventStoreFile = filepath.Join(xdg.DataHome, "t5", eventStoreFile)
		}
	}

	// Ensure the directory exists
	err := os.MkdirAll(filepath.Dir(eventStoreFile), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory for event store: %w", err)
	}

	// Create the append log event store
	store, err := core.NewAppendLogEventStore(eventStoreFile)
	if err != nil {
		return fmt.Errorf("failed to create append log event store: %w", err)
	}

	// Create a repository with the append log store
	repo := core.NewRepositoryWithEventStore(store)

	// Rebuild the state from the event store
	err = repo.RebuildState()
	if err != nil {
		return fmt.Errorf("failed to rebuild state from event store: %w", err)
	}

	ctx.Repository = repo
	return nil
}