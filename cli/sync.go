package cli

import (
	"fmt"

	"t5.mkbrechtel.dev/t5/core"
)

type SyncGroupsConfig []SyncGroupConfig

type SyncGroupConfig struct {
	Name         string   `yaml:"name"`
	TodoTxtFiles []string `yaml:"todotxt_files"`
}

// ProviderConfig contains configuration for a sync provider
type SyncProviderBaseConfig struct {
	Direction string `yaml:"direction"`

	// ImportFilter specifies a filter to apply when importing tasks
	ImportFilter core.FilterComposer `yaml:"import_filter"`

	// ImportModifier specifies modifiers to apply when importing tasks
	ImportModifier core.ModifierComposer `yaml:"import_modifier"`

	// ExportFilter specifies a filter to apply when exporting tasks
	ExportFilter core.FilterComposer `yaml:"export_filter"`

	// ExportModifier specifies modifiers to apply when exporting tasks
	ExportModifier core.ModifierComposer `yaml:"export_modifier"`
}

// SyncTasks synchronizes tasks with external sources
func SyncTasks(ctx *AppContext, todoFilePath string) error {
	// TODO: implement sync functionality
	return fmt.Errorf("sync functionality not yet implemented")
}
