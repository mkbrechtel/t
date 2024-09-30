package cmd

import (
	"fmt"
	"log"
	yaml "gopkg.in/yaml.v2"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var configFile string

// todoCmd represents the todo command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show your t0 config",
	Long: `With this command you can show your t0 configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(yamlStringSettings())
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

func init() {
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	// Don't forget to read config either from configFile or from home directory!
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Use default config file from xdg config dirs
		for _, dir := range xdg.ConfigDirs {
			viper.AddConfigPath(filepath.Join(dir, "t0"))
		}
		viper.AddConfigPath(filepath.Join(xdg.ConfigHome, "t0"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	viper.ReadInConfig() // Find and read the config file
}
