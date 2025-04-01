package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var configFile string

// TodoFileConfig represents the configuration for a single todo.txt file
type TodoFileConfig struct {
	Path   string            `yaml:"path"`
	Ensure map[string]bool   `yaml:"ensure,omitempty"`
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show your t5 config",
	Long: `t5 config

	With this command you can show your t5 configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print output in a format that matches the tests
		output := yamlStringSettings()
		fmt.Print(output)
	},
}

func yamlStringSettings() string {
	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	return string(bs)
}

// GetTodoFiles retrieves the configured todo.txt files from config
func GetTodoFiles() ([]TodoFileConfig, error) {
	var todoFiles []TodoFileConfig
	
	// Check for legacy single file config first
	legacyFilePath := viper.GetString("todo.file")
	if legacyFilePath != "" {
		// Create a config for the legacy file with default settings
		todoFile := TodoFileConfig{
			Path: legacyFilePath,
			Ensure: map[string]bool{
				"prefershortids":       viper.GetBool("todo.ensure.prefershortids"),
				"enforcecompletiondate": viper.GetBool("todo.ensure.enforcecompletiondate"),
				"enforcecreationdate":   viper.GetBool("todo.ensure.enforcecreationdate"),
			},
		}
		return []TodoFileConfig{todoFile}, nil
	}
	
	// Try to unmarshal the new multi-file config
	if err := viper.UnmarshalKey("todo.files", &todoFiles); err != nil {
		return nil, fmt.Errorf("failed to parse todo files from config: %w", err)
	}
	
	if len(todoFiles) == 0 {
		// If no todo files are configured, use default todo.txt in current directory
		todoFile := TodoFileConfig{
			Path: "todo.txt",
			Ensure: map[string]bool{
				"prefershortids":        true,
				"enforcecompletiondate": true,
				"enforcecreationdate":   true,
			},
		}
		return []TodoFileConfig{todoFile}, nil
	}
	
	return todoFiles, nil
}

// ResetConfig resets the viper configuration for testing purposes
func ResetConfig() {
	viper.Reset()
	initConfig()
}

func init() {
	rootCmd.AddCommand(configCmd)
	
	// Add event store configuration flag
	rootCmd.PersistentFlags().String("eventstore", "", "Event store file path (if not specified, uses in-memory storage)")
	viper.BindPFlag("eventstore.file", rootCmd.PersistentFlags().Lookup("eventstore"))
}

func initConfig() {
	// Don't forget to read config either from configFile or from home directory!
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Use default config file from xdg config dirs
		for _, dir := range xdg.ConfigDirs {
			viper.AddConfigPath(filepath.Join(dir, "t"))
		}
		viper.AddConfigPath(filepath.Join(xdg.ConfigHome, "t"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	viper.ReadInConfig() // Find and read the config file
}
