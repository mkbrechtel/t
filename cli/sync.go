package cli

import (
	"fmt"
	"log"
	"strings"
	
	"t5.mkbrechtel.dev/t5/core"
	"t5.mkbrechtel.dev/t5/sync"
	todo "t5.mkbrechtel.dev/t5/sync/todotxt"
)

// SyncTasks synchronizes tasks between various sources and the event store
// Supports multiple sync providers and filtering/modifying tasks during sync
func SyncTasks(ctx *AppContext) {
	// Check if a specific provider was specified
	if len(ctx.Args) > 0 {
		provider := ctx.Args[0]
		
		// Check if it's help flag
		if provider == "--help" || provider == "-h" {
			fmt.Println("Usage:")
			fmt.Println("  t5 sync                        - Sync with the default todo.txt file")
			fmt.Println("  t5 sync [file]                 - Sync with the specified todo.txt file")
			fmt.Println("  t5 sync provider:<name>        - Sync with a specific configured provider")
			fmt.Println("                                   (e.g., provider:github, provider:gitlab)")
			fmt.Println("Configuration for providers should be in the t5 config file.")
			return
		}
		
		// If it's a file path, assume todo.txt sync
		if !strings.HasPrefix(provider, "provider:") {
			// Legacy todo.txt file sync
			todoFile := provider
			result, err := todo.SyncWithRepository(ctx.Repository, todoFile)
			if err != nil {
				log.Fatalf("Sync failed for %s: %v", todoFile, err)
			}
			
			// Print results
			fmt.Printf("Sync completed for %s:\n", todoFile)
			fmt.Printf("  Added: %d\n", result.Added)
			fmt.Printf("  Updated from todo.txt: %d\n", result.FromTodoTxt)
			fmt.Printf("  Updated to todo.txt: %d\n", result.ToTodoTxt)
			fmt.Printf("  Skipped: %d\n", result.Skipped)
			if result.Conflicts > 0 {
				fmt.Printf("  Conflicts detected: %d\n", result.Conflicts)
			}
			return
		}
		
		// Extract provider name from provider:name format
		providerName := strings.TrimPrefix(provider, "provider:")
		
		// Look for the provider in the config
		var providerConfig *sync.ProviderConfig
		for _, pc := range ctx.Config.SyncProviders {
			if pc.Type == providerName {
				providerConfig = &pc
				break
			}
		}
		
		if providerConfig == nil {
			log.Fatalf("Provider '%s' not found in configuration", providerName)
		}
		
		// Create the provider
		syncProvider, err := sync.CreateProvider(*providerConfig)
		if err != nil {
			log.Fatalf("Failed to create provider '%s': %v", providerName, err)
		}
		
		// Determine the sync direction
		direction := sync.SyncDirectionBoth
		switch strings.ToLower(providerConfig.Direction) {
		case "import":
			direction = sync.SyncDirectionImport
		case "export":
			direction = sync.SyncDirectionExport
		}
		
		// Create filters and modifiers
		var importFilter, exportFilter core.TaskFilter
		var importModifier, exportModifier core.TaskModifier
		
		// Import filter/modifier
		if providerConfig.ImportFilter != nil {
			importFilter, err = sync.BuildFilterFromConfig(providerConfig.ImportFilter)
			if err != nil {
				log.Fatalf("Failed to build import filter: %v", err)
			}
		}
		
		if providerConfig.ImportModifier != nil {
			importModifier, err = sync.BuildModifierFromConfig(providerConfig.ImportModifier)
			if err != nil {
				log.Fatalf("Failed to build import modifier: %v", err)
			}
		}
		
		// Export filter/modifier
		if providerConfig.ExportFilter != nil {
			exportFilter, err = sync.BuildFilterFromConfig(providerConfig.ExportFilter)
			if err != nil {
				log.Fatalf("Failed to build export filter: %v", err)
			}
		}
		
		if providerConfig.ExportModifier != nil {
			exportModifier, err = sync.BuildModifierFromConfig(providerConfig.ExportModifier)
			if err != nil {
				log.Fatalf("Failed to build export modifier: %v", err)
			}
		}
		
		// Perform the sync based on direction
		var result *sync.SyncResult
		
		fmt.Printf("Syncing with %s (%s)...\n", syncProvider.Name(), syncProvider.Description())
		
		switch direction {
		case sync.SyncDirectionImport:
			result, err = syncProvider.Import(ctx.Repository, importFilter, importModifier)
			if err != nil {
				log.Fatalf("Import failed: %v", err)
			}
		case sync.SyncDirectionExport:
			result, err = syncProvider.Export(ctx.Repository, exportFilter, exportModifier)
			if err != nil {
				log.Fatalf("Export failed: %v", err)
			}
		case sync.SyncDirectionBoth:
			result, err = syncProvider.Sync(ctx.Repository, nil, nil) // Use provider's own sync method
			if err != nil {
				log.Fatalf("Sync failed: %v", err)
			}
		}
		
		// Print results
		fmt.Printf("Sync completed for %s:\n", syncProvider.Name())
		fmt.Printf("  Added: %d\n", result.Added)
		fmt.Printf("  Updated: %d\n", result.Updated)
		fmt.Printf("  From source: %d\n", result.FromSource)
		fmt.Printf("  To source: %d\n", result.ToSource)
		fmt.Printf("  Skipped: %d\n", result.Skipped)
		if result.Conflicts > 0 {
			fmt.Printf("  Conflicts detected: %d\n", result.Conflicts)
		}
		
		return
	}
	
	// If no provider specified, sync with all enabled providers
	results := make([]*sync.SyncResult, 0)
	
	// First sync with the default todo.txt file
	if ctx.Config.TodoFile != "" {
		result, err := todo.SyncWithRepository(ctx.Repository, ctx.Config.TodoFile)
		if err != nil {
			log.Fatalf("Sync failed for %s: %v", ctx.Config.TodoFile, err)
		}
		
		// Add to results
		syncResult := &sync.SyncResult{
			Added:      result.Added,
			Updated:    result.FromTodoTxt + result.ToTodoTxt,
			Skipped:    result.Skipped,
			Conflicts:  result.Conflicts,
			FromSource: result.FromTodoTxt,
			ToSource:   result.ToTodoTxt,
		}
		results = append(results, syncResult)
		
		// Print individual results
		fmt.Printf("Sync completed for %s:\n", ctx.Config.TodoFile)
		fmt.Printf("  Added: %d\n", result.Added)
		fmt.Printf("  Updated from todo.txt: %d\n", result.FromTodoTxt)
		fmt.Printf("  Updated to todo.txt: %d\n", result.ToTodoTxt)
		fmt.Printf("  Skipped: %d\n", result.Skipped)
		if result.Conflicts > 0 {
			fmt.Printf("  Conflicts detected: %d\n", result.Conflicts)
		}
		fmt.Println()
	}
	
	// Then sync with all enabled providers
	for _, providerConfig := range ctx.Config.SyncProviders {
		if !providerConfig.Enabled {
			continue
		}
		
		// Create the provider
		syncProvider, err := sync.CreateProvider(providerConfig)
		if err != nil {
			fmt.Printf("Failed to create provider '%s': %v\n", providerConfig.Type, err)
			continue
		}
		
		// Determine the sync direction
		direction := sync.SyncDirectionBoth
		switch strings.ToLower(providerConfig.Direction) {
		case "import":
			direction = sync.SyncDirectionImport
		case "export":
			direction = sync.SyncDirectionExport
		}
		
		// Create filters and modifiers
		var importFilter, exportFilter core.TaskFilter
		var importModifier, exportModifier core.TaskModifier
		
		// Import filter/modifier
		if providerConfig.ImportFilter != nil {
			importFilter, err = sync.BuildFilterFromConfig(providerConfig.ImportFilter)
			if err != nil {
				fmt.Printf("Failed to build import filter for %s: %v\n", providerConfig.Type, err)
				continue
			}
		}
		
		if providerConfig.ImportModifier != nil {
			importModifier, err = sync.BuildModifierFromConfig(providerConfig.ImportModifier)
			if err != nil {
				fmt.Printf("Failed to build import modifier for %s: %v\n", providerConfig.Type, err)
				continue
			}
		}
		
		// Export filter/modifier
		if providerConfig.ExportFilter != nil {
			exportFilter, err = sync.BuildFilterFromConfig(providerConfig.ExportFilter)
			if err != nil {
				fmt.Printf("Failed to build export filter for %s: %v\n", providerConfig.Type, err)
				continue
			}
		}
		
		if providerConfig.ExportModifier != nil {
			exportModifier, err = sync.BuildModifierFromConfig(providerConfig.ExportModifier)
			if err != nil {
				fmt.Printf("Failed to build export modifier for %s: %v\n", providerConfig.Type, err)
				continue
			}
		}
		
		// Perform the sync based on direction
		var result *sync.SyncResult
		
		fmt.Printf("Syncing with %s (%s)...\n", syncProvider.Name(), syncProvider.Description())
		
		switch direction {
		case sync.SyncDirectionImport:
			result, err = syncProvider.Import(ctx.Repository, importFilter, importModifier)
			if err != nil {
				fmt.Printf("Import failed for %s: %v\n", providerConfig.Type, err)
				continue
			}
		case sync.SyncDirectionExport:
			result, err = syncProvider.Export(ctx.Repository, exportFilter, exportModifier)
			if err != nil {
				fmt.Printf("Export failed for %s: %v\n", providerConfig.Type, err)
				continue
			}
		case sync.SyncDirectionBoth:
			result, err = syncProvider.Sync(ctx.Repository, nil, nil) // Use provider's own sync method
			if err != nil {
				fmt.Printf("Sync failed for %s: %v\n", providerConfig.Type, err)
				continue
			}
		}
		
		// Add to results
		results = append(results, result)
		
		// Print individual results
		fmt.Printf("Sync completed for %s:\n", syncProvider.Name())
		fmt.Printf("  Added: %d\n", result.Added)
		fmt.Printf("  Updated: %d\n", result.Updated)
		fmt.Printf("  From source: %d\n", result.FromSource)
		fmt.Printf("  To source: %d\n", result.ToSource)
		fmt.Printf("  Skipped: %d\n", result.Skipped)
		if result.Conflicts > 0 {
			fmt.Printf("  Conflicts detected: %d\n", result.Conflicts)
		}
		fmt.Println()
	}
	
	// Print overall results
	if len(results) > 0 {
		mergedResult := sync.MergeSyncResults(results)
		
		fmt.Println("Overall sync results:")
		fmt.Printf("  Providers synced: %d\n", len(results))
		fmt.Printf("  Total added: %d\n", mergedResult.Added)
		fmt.Printf("  Total updated: %d\n", mergedResult.Updated)
		fmt.Printf("  Total from sources: %d\n", mergedResult.FromSource)
		fmt.Printf("  Total to sources: %d\n", mergedResult.ToSource)
		fmt.Printf("  Total skipped: %d\n", mergedResult.Skipped)
		if mergedResult.Conflicts > 0 {
			fmt.Printf("  Total conflicts: %d\n", mergedResult.Conflicts)
		}
	} else {
		fmt.Println("No providers were synced.")
	}
}