package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with external services",
	Long: `t sync
	
	With this command you can sync with external services. 
	`,
}


func init() {
    rootCmd.AddCommand(syncCmd)
    syncCmd.AddCommand(syncOpenProjectCmd)

    syncOpenProjectCmd.PersistentFlags().String("url", "", "OpenProject URL")
    syncOpenProjectCmd.PersistentFlags().String("api-key", "", "OpenProject API Key")
    syncOpenProjectCmd.PersistentFlags().String("query-id", "", "OpenProject Query ID")
    syncOpenProjectCmd.PersistentFlags().String("todo-prefix", "", "Prefix for OpenProject todos created by t")

    viper.BindPFlag("sync.openproject.url", syncOpenProjectCmd.PersistentFlags().Lookup("url"))
    viper.BindPFlag("sync.openproject.api-key", syncOpenProjectCmd.PersistentFlags().Lookup("api-key"))
    viper.BindPFlag("sync.openproject.query-id", syncOpenProjectCmd.PersistentFlags().Lookup("query-id"))
    viper.BindPFlag("sync.openproject.todo-prefix", syncOpenProjectCmd.PersistentFlags().Lookup("todo-prefix"))

}
