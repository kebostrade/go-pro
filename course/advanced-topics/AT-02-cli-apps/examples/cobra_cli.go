package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taskcli",
	Short: "A task management CLI",
	Long:  "Manage your tasks from the command line",
}

var addCmd = &cobra.Command{
	Use:   "add [task]",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		task := args[0]
		priority, _ := cmd.Flags().GetString("priority")
		fmt.Printf("Added task: %s (priority: %s)\n", task, priority)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		status, _ := cmd.Flags().GetString("status")
		fmt.Printf("Listing tasks (all: %v, status: %s)\n", all, status)
	},
}

var doneCmd = &cobra.Command{
	Use:   "done [task-id]",
	Short: "Mark a task as done",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		fmt.Printf("Task %s marked as done\n", taskID)
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [task-id]",
	Short: "Remove a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Remove task %s? (y/n): ", taskID)
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" {
				fmt.Println("Cancelled")
				return
			}
		}
		fmt.Printf("Removed task %s\n", taskID)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("taskcli v1.0.0")
	},
}

func init() {
	addCmd.Flags().StringP("priority", "p", "medium", "Task priority (low, medium, high)")
	listCmd.Flags().BoolP("all", "a", false, "Show all tasks including completed")
	listCmd.Flags().StringP("status", "s", "pending", "Filter by status")
	removeCmd.Flags().BoolP("force", "f", false, "Force removal without confirmation")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
