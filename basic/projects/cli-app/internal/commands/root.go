package commands

import (
	"github.com/spf13/cobra"
)

// configPath is the path to the config file
var configPath string

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "CLI app demonstrating Cobra patterns",
		Long:  `A CLI application demonstrating Go Cobra patterns with subcommands, configuration management, and testing.`,
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yaml", "Path to config file")

	// Add subcommands
	rootCmd.AddCommand(NewGreetCmd())
	rootCmd.AddCommand(NewServeCmd())

	return rootCmd
}

// Execute runs the root command
func Execute() error {
	return NewRootCmd().Execute()
}
