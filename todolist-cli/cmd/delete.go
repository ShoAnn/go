/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"todolist/internal"
	"os"
	"strconv"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `Delete allows you to remove resources from the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		//validate argument (int)
		taskId, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid input", err)
			os.Exit(1)
		}

		list := &internal.Todolist{}
		if err := list.ReadFromFile("main.tdl.json"); err != nil {
			fmt.Fprintln(os.Stderr, "error reading file", err)
		}

		fmt.Println(taskId)
		//check oob
		if taskId == 0 || taskId > len(*list) {
			fmt.Println("index out of bound")
			os.Exit(1)
		}
		if err := list.Delete(taskId - 1); err != nil {
			fmt.Fprintln(os.Stderr, "task deletion failed", err)
			os.Exit(1)
		}
		if err := list.Store("main.tdl.json"); err != nil {
			fmt.Fprintln(os.Stderr, "error storing to file:", err)
			os.Exit(1)
		}
		fmt.Printf("Task %d deleted succesfully", taskId)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
