package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "t",
	Short: "t is a todo and time tracker",
	Long: `t is a todo and time tracker

	manage your todo list with
		t todo

	track your time with
		t time
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
