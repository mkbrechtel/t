package sync

import "t5.mkbrechtel.dev/t5/core"


// ProviderConfig contains configuration for a sync provider
type ProviderBaseConfig struct {
	// Type specifies the type of sync provider (e.g., "todotxt", "github", "gitlab")
	Type string `yaml:"type"`
	
	// Enabled indicates whether this provider is enabled
	Enabled bool `yaml:"enabled"`
	
	// Direction specifies the direction of synchronization
	// Can be "import", "export", or "both"
	Direction string `yaml:"direction"`
		
	// ImportFilter specifies a filter to apply when importing tasks
	ImportFilter core.TaskFilter `yaml:"import_filter"`
	
	// ImportModifier specifies modifiers to apply when importing tasks
	ImportModifier core.TaskModifier `yaml:"import_modifier"`
	
	// ExportFilter specifies a filter to apply when exporting tasks
	ExportFilter core.TaskFilter `yaml:"export_filter"`
	
	// ExportModifier specifies modifiers to apply when exporting tasks
	ExportModifier core.TaskModifier `yaml:"export_modifier"`
}
