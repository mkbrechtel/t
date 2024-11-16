package cmd

import (
    "github.com/spf13/cobra"
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
}
