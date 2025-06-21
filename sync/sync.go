package sync

import (
	"t5.mkbrechtel.dev/t5/core"
)

// SyncDirection represents the direction of synchronization
type SyncDirection int

const (
	// SyncDirectionImport imports data from external source to t5
	SyncDirectionImport SyncDirection = iota
	// SyncDirectionExport exports data from t5 to external source
	SyncDirectionExport
	// SyncDirectionBoth synchronizes data in both directions
	SyncDirectionBoth
)

// ProviderConfig contains configuration parameters for sync providers
type ProviderConfig struct {
	// Params contains provider-specific configuration parameters
	Params map[string]interface{}
}

// SyncProviderFactory is a function that creates a sync provider from configuration
type SyncProviderFactory func(config ProviderConfig) (SyncProvider, error)

// SyncResult contains statistics about a sync operation
type SyncResult struct {
	Added      int
	Updated    int
	Skipped    int
	Conflicts  int
	FromSource int
	ToSource   int
}

// SyncProvider defines the interface for any sync provider
type SyncProvider interface {

	// Import imports tasks from the source into the t5 repository
	// Returns a SyncResult with statistics about the operation
	Import(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*SyncResult, error)

	// Export exports tasks from the t5 repository to the destination
	// Returns a SyncResult with statistics about the operation
	Export(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*SyncResult, error)

	// Sync synchronizes tasks in both directions
	// Returns a SyncResult with statistics about the operation
	Sync(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*SyncResult, error)

	ProvidesImport() bool
	ProvidesExport() bool
}
