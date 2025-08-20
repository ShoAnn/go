/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"todolist/internal"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: "create task",
	Run: func(cmd *cobra.Command, args []string) {
		list := &internal.Todolist{}
		filePath := "main.tdl.json"
		if err := list.ReadFromFile(filePath); err == nil {
			fmt.Fprintln(os.Stderr, "error reading file:", err)
			os.Exit(1)
		}
		taskName := strings.Join(args, "")
		list.Create(taskName)
		if err := list.Store(filePath); err != nil {
			fmt.Fprintln(os.Stderr, "error storing to file:", err)
			os.Exit(1)
		}
		fmt.Printf("successfully added task: %s", taskName)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
