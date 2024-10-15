package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"

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
		_, err := openproject.GetWorkPackages(syncOpenProjectUrl,syncOpenProjectApiKey)

		fmt.Print(err.Error())

		// // Check if the response contains work packages
		// if len(opResponse.WorkPackages) == 0 {
		// 	fmt.Println("No work packages found.")
		// 	return
		// }

		// // Process each work package
		// for _, wp := range opResponse.WorkPackages {
		// 	fmt.Printf("Work Package ID: %d\n", wp.ID)
		// 	fmt.Printf("Subject: %s\n", wp.Subject)
		// 	fmt.Printf("Description (Raw): %s\n", wp.Description.Raw)
		// 	// You can add more processing here, like saving to a database or performing actions based on the work package
		// }
		return
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
