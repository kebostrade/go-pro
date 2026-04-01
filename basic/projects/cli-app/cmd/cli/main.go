package main

import (
	"fmt"
	"os"

	"github.com/DimaJoyti/go-pro/basic/projects/cli-app/internal/commands"
)

func main() {
	rootCmd := commands.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
