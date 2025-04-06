package sync

import (
	"t5.mkbrechtel.dev/t5/core"
)

// SyncDirection indicates the direction of synchronization
type SyncDirection int

const (
	// SyncDirectionImport imports tasks from source to t5
	SyncDirectionImport SyncDirection = iota
	// SyncDirectionExport exports tasks from t5 to destination
	SyncDirectionExport
	// SyncDirectionBoth synchronizes in both directions
	SyncDirectionBoth
)

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
	// Name returns the name of the sync provider
	Name() string
	
	// Description returns a description of the sync provider
	Description() string
	
	// SupportedDirections returns the sync directions supported by this provider
	SupportedDirections() []SyncDirection
	
	// Import imports tasks from the source into the t5 repository
	// Returns a SyncResult with statistics about the operation
	Import(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*SyncResult, error)
	
	// Export exports tasks from the t5 repository to the destination
	// Returns a SyncResult with statistics about the operation
	Export(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*SyncResult, error)
	
	// Sync synchronizes tasks in both directions
	// Returns a SyncResult with statistics about the operation
	Sync(repo *core.Repository, filter core.TaskFilter, modifier core.TaskModifier) (*SyncResult, error)
}

// ProviderConfig contains configuration for a sync provider
type ProviderConfig struct {
	// Type specifies the type of sync provider (e.g., "todotxt", "github", "gitlab")
	Type string `yaml:"type"`
	
	// Enabled indicates whether this provider is enabled
	Enabled bool `yaml:"enabled"`
	
	// Direction specifies the direction of synchronization
	// Can be "import", "export", or "both"
	Direction string `yaml:"direction"`
	
	// Params contains provider-specific parameters
	Params map[string]interface{} `yaml:"params"`
	
	// ImportFilter specifies a filter to apply when importing tasks
	ImportFilter map[string]interface{} `yaml:"import_filter"`
	
	// ImportModifier specifies modifiers to apply when importing tasks
	ImportModifier map[string]interface{} `yaml:"import_modifier"`
	
	// ExportFilter specifies a filter to apply when exporting tasks
	ExportFilter map[string]interface{} `yaml:"export_filter"`
	
	// ExportModifier specifies modifiers to apply when exporting tasks
	ExportModifier map[string]interface{} `yaml:"export_modifier"`
}

// SyncProviderFactory creates a new sync provider from a configuration
type SyncProviderFactory func(config ProviderConfig) (SyncProvider, error)

// MergeSyncResults combines multiple sync results into one
func MergeSyncResults(results []*SyncResult) *SyncResult {
	merged := &SyncResult{}
	
	for _, result := range results {
		if result == nil {
			continue
		}
		
		merged.Added += result.Added
		merged.Updated += result.Updated
		merged.Skipped += result.Skipped
		merged.Conflicts += result.Conflicts
		merged.FromSource += result.FromSource
		merged.ToSource += result.ToSource
	}
	
	return merged
}