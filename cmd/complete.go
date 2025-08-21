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

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "complete a  task",
	Run: func(cmd *cobra.Command, args []string) {
		//validate argument (int)
		taskId, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid input", err)
			os.Exit(1)
		}

		//get main list
		list := &internal.Todolist{}
		if err := list.ReadFromFile("main.tdl.json"); err != nil {
			fmt.Fprintln(os.Stderr, "error reading file:", err)
			os.Exit(1)
		}

		//check oob
		if taskId == 0 || taskId >= len(*list) {
			fmt.Println("index out of bound")
			os.Exit(1)
		}

		//check already completed
		if (*list)[taskId-1].Completed {
			fmt.Println("task already completed")
			os.Exit(1)
		}

		if err := list.Complete(taskId-1); err != nil {
			fmt.Fprintln(os.Stderr, "task completion failed", err)
			fmt.Println(taskId)
			os.Exit(1)
		}
		if err := list.Store("main.tdl.json"); err != nil {
			fmt.Fprintln(os.Stderr, "error storing to file:", err)
			os.Exit(1)
		}
		completedTaskName := (*list)[taskId-1].Name
		fmt.Printf("successfully completed task: %s", completedTaskName)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
