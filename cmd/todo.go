/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage your todo list",
	Long: `With this command you can manage your todo list. The todo list is based on the todo.txt format.`,
}


// todoAddCmd represents the todo add command
var todoAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to your todo list",
	Long: `With this command you can add tasks to your todo list.

It either takes the task from the command line argument or multiple tasks in the form of todo.txt lines as standard input.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("todoAdd called")
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// todoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")


	todoCmd.AddCommand(todoAddCmd)
}

