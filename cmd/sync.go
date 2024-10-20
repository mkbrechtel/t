package cmd

import (
	"os"

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
		url := viper.GetString("sync.openproject.url")
		apiKey := viper.GetString("sync.openproject.api-key")
		queryId := viper.GetString("sync.openproject.query-id")

		workPackages,_ := openproject.GetWorkPackages(url, apiKey, queryId)
		openproject.PrintWorkPackages(workPackages)
		tl := openproject.CreateTaskList(workPackages)
		tl.WriteToFile(os.Stdout)
	},
}

func init() {
    rootCmd.AddCommand(syncCmd)
    syncCmd.AddCommand(syncOpenProjectCmd)

    syncOpenProjectCmd.PersistentFlags().String("url", "", "OpenProject URL")
    syncOpenProjectCmd.PersistentFlags().String("api-key", "", "OpenProject API Key")
    syncOpenProjectCmd.PersistentFlags().String("query-id", "", "OpenProject Query ID")

    viper.BindPFlag("sync.openproject.url", syncOpenProjectCmd.PersistentFlags().Lookup("url"))
    viper.BindPFlag("sync.openproject.api-key", syncOpenProjectCmd.PersistentFlags().Lookup("api-key"))
    viper.BindPFlag("sync.openproject.query-id", syncOpenProjectCmd.PersistentFlags().Lookup("query-id"))

}
