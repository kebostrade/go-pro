package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DimaJoyti/go-pro/basic/projects/cli-app/internal/config"
	"github.com/DimaJoyti/go-pro/basic/projects/cli-app/pkg/greeting"
	"github.com/spf13/cobra"
)

var (
	port string
	host string
)

// NewServeCmd creates the serve subcommand
func NewServeCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		Long:  `Starts an HTTP server that returns greeting messages. Use --port and --host to configure.`,
		Run:   runServe,
	}

	serveCmd.Flags().StringVar(&port, "port", "8080", "Port to listen on")
	serveCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind to")

	return serveCmd
}

func runServe(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		// Use defaults if config loading fails
		cfg = &config.Config{}
	}

	// Override config with CLI flags if provided
	if cmd.Flags().Changed("port") {
		cfg.Server.Port = port
	}
	if cmd.Flags().Changed("host") {
		cfg.Server.Host = host
	}

	// Use flag values directly if config wasn't loaded
	listenPort := port
	listenHost := host
	if cfg != nil {
		if !cmd.Flags().Changed("port") && cfg.Server.Port != "" {
			listenPort = cfg.Server.Port
		}
		if !cmd.Flags().Changed("host") && cfg.Server.Host != "" {
			listenHost = cfg.Server.Host
		}
	}

	addr := fmt.Sprintf("%s:%s", listenHost, listenPort)

	// Set up HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go CLI App!")
	})
	http.HandleFunc("/greet/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/greet/"):]
		if name == "" {
			name = "World"
		}
		fmt.Fprintf(w, "%s", greeting.Greet(name, 1))
	})

	log.Printf("Starting server on %s", addr)

	// Start server in goroutine
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")
}
