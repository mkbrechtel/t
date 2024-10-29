package cmd

import (
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "t/sync/openproject"
)


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
		prefix := viper.GetString("sync.openproject.todo-prefix")

		workPackages,_ := openproject.GetWorkPackages(url, apiKey, queryId)
		openproject.PrintWorkPackages(workPackages)
		tl := openproject.CreateTaskList(workPackages, prefix, url)
		tl.WriteToFile(os.Stdout)
	},
}
