/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"todolist/internal"
	"os"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "display whole list",
	Run: func(cmd *cobra.Command, args []string) {
		todolist := &internal.Todolist{}
		if err := todolist.ReadFromFile("main.tdl.json"); err != nil {
			fmt.Fprintln(os.Stderr, "error reading file:", err)
			os.Exit(1)
		}
		if len(*todolist) == 0 {
			fmt.Println("You currently got nothing to do.")
		}
		fmt.Println("You gotta do:")
		for i, task := range *todolist {
			status := ""
			if task.Completed {
				status = "--done--"
			}
			fmt.Printf("%d | %s %s\n", i+1, task.Name, status)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
