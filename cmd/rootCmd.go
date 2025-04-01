package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "t5",
	Short: "t5 is a todo and time tracker",
	Long: `t5 is a todo and time tracker

	manage your todo list with
		t5 todo

	track your time with
		t5 time
`,
}

var todoFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GetRootCommandForTesting returns a new instance of the root command for testing purposes.
// This allows tests to get a fresh command instance with all subcommands properly initialized.
func GetRootCommandForTesting() *cobra.Command {
	// Create a new root command with the same configuration
	testRootCmd := &cobra.Command{
		Use:   rootCmd.Use,
		Short: rootCmd.Short,
		Long:  rootCmd.Long,
	}
	
	// Add all the same flags
	testRootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	testRootCmd.PersistentFlags().StringVarP(&todoFile, "todoFile", "t", "todo.txt", "todo.txt file")
	testRootCmd.PersistentFlags().String("eventstore", "", "Event store file path (if not specified, uses in-memory storage)")
	
	// Bind flags to viper
	viper.BindPFlag("todo.file", testRootCmd.PersistentFlags().Lookup("todoFile"))
	viper.BindPFlag("eventstore.file", testRootCmd.PersistentFlags().Lookup("eventstore"))
	
	// Add all subcommands
	testRootCmd.AddCommand(configCmd)
	testRootCmd.AddCommand(listCmd)
	testRootCmd.AddCommand(todoCmd)
	testRootCmd.AddCommand(syncCmd)
	
	// Initialize config when command runs
	cobra.OnInitialize(initConfig)
	
	return testRootCmd
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&todoFile, "todoFile", "t", "todo.txt", "todo.txt file")
	viper.BindPFlag("todo.file", rootCmd.PersistentFlags().Lookup("todoFile"))
}
