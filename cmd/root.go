package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "t0",
	Short: "t0 is a todo and time tracker",
	Long: `t0 is a todo and time tracker

	manage your todo list with
		t0 todo

	track your time with
		t0 time
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

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&todoFile, "todoFile", "t", "todo.txt", "todo.txt file")
	viper.BindPFlag("todoFile", rootCmd.PersistentFlags().Lookup("todoFile"))
}
