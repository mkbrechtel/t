package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
	"t5.mkbrechtel.dev/core"
)

// InitRepository initializes a repository with the appropriate event store
// based on configuration
func InitRepository() (*core.Repository, error) {
	// Check if an event store file is specified
	eventStorePath := viper.GetString("eventstore.file")
	if eventStorePath == "" {
		// Use in-memory event store
		return core.NewRepository(), nil
	}

	// If the path is not absolute, make it relative to the working directory or data home
	// Paths starting with ./ are relative to the working directory
	if !filepath.IsAbs(eventStorePath) {
		if len(eventStorePath) >= 2 && eventStorePath[0:2] == "./" {
			// Keep as relative path for testing
			wd, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get current working directory: %w", err)
			}
			eventStorePath = filepath.Join(wd, eventStorePath[2:])
		} else {
			// Otherwise use the data home directory
			eventStorePath = filepath.Join(xdg.DataHome, "t5", eventStorePath)
		}
	}

	// Ensure the directory exists
	err := os.MkdirAll(filepath.Dir(eventStorePath), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory for event store: %w", err)
	}

	// Create the append log event store
	store, err := core.NewAppendLogEventStore(eventStorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create append log event store: %w", err)
	}

	// Create a repository with the append log store
	repo := core.NewRepositoryWithEventStore(store)

	// Rebuild the state from the event store
	err = repo.RebuildState()
	if err != nil {
		return nil, fmt.Errorf("failed to rebuild state from event store: %w", err)
	}

	return repo, nil
}