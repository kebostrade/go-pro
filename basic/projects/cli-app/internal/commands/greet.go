package commands

import (
	"fmt"

	"github.com/DimaJoyti/go-pro/basic/projects/cli-app/internal/config"
	"github.com/DimaJoyti/go-pro/basic/projects/cli-app/pkg/greeting"
	"github.com/spf13/cobra"
)

var (
	name  string
	times int
)

// NewGreetCmd creates the greet subcommand
func NewGreetCmd() *cobra.Command {
	greetCmd := &cobra.Command{
		Use:   "greet",
		Short: "Greet someone",
		Long:  `Prints a greeting message. Use --name to specify who to greet and --times to specify how many times.`,
		Run:   runGreet,
	}

	greetCmd.Flags().StringVar(&name, "name", "World", "Name of the person to greet")
	greetCmd.Flags().IntVar(&times, "times", 1, "Number of times to repeat the greeting")

	return greetCmd
}

func runGreet(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// Use defaults if config loading fails
		cfg = &config.Config{}
	}

	// Override config with CLI flags if provided
	if cmd.Flags().Changed("name") {
		cfg.Greeting.DefaultName = name
	}
	if cmd.Flags().Changed("times") {
		cfg.Greeting.DefaultTimes = times
	}

	// Use flag values directly if config wasn't loaded
	greetName := name
	greetTimes := times
	if cfg != nil {
		if !cmd.Flags().Changed("name") && cfg.Greeting.DefaultName != "" {
			greetName = cfg.Greeting.DefaultName
		}
		if !cmd.Flags().Changed("times") && cfg.Greeting.DefaultTimes > 0 {
			greetTimes = cfg.Greeting.DefaultTimes
		}
	}

	msg := greeting.Greet(greetName, greetTimes)
	fmt.Print(msg)
}
