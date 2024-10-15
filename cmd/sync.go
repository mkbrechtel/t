package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"t/sync/openproject"

)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with external services",
	Long: `t sync
	
	With this command you can sync with external services. 
	`,
}

// openproject command
var syncOpenProjectCmd = &cobra.Command{
	Use:   "openproject",
	Short: "Sync with OpenProject",
	Long: `t sync openproject

	With this command you can sync with OpenProject.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Fetch the work packages from OpenProject
		syncOpenProjectUrl := viper.GetString("sync.openproject.url")
		syncOpenProjectApiKey := viper.GetString("sync.openproject.api-key")
		openproject.GetWorkPackages(syncOpenProjectUrl,syncOpenProjectApiKey)
	},
}

func init() {
    rootCmd.AddCommand(syncCmd)
    syncCmd.AddCommand(syncOpenProjectCmd)

    syncOpenProjectCmd.PersistentFlags().String("url", "", "OpenProject URL")
    syncOpenProjectCmd.PersistentFlags().String("api-key", "", "OpenProject API Key")

    viper.BindPFlag("sync.openproject.url", syncOpenProjectCmd.PersistentFlags().Lookup("url"))
    viper.BindPFlag("sync.openproject.api-key", syncOpenProjectCmd.PersistentFlags().Lookup("api-key"))

}
