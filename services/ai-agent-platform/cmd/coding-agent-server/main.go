package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/agent"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/api"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/languages/common"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/languages/golang"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/llm"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/internal/tools/programming"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

const (
	defaultPort    = "8080"
	defaultTimeout = 60 * time.Second
)

func main() {
	fmt.Println("🤖 Coding Agent API Server")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println()

	// Load configuration
	config := loadConfig()

	// Initialize components
	llmProvider, err := initializeLLM(config)
	if err != nil {
		log.Fatalf("❌ Failed to initialize LLM: %v", err)
	}
	fmt.Println("✅ LLM provider initialized")

	languageRegistry := initializeLanguages()
	fmt.Println("✅ Language registry initialized")

	tools := initializeTools(languageRegistry)
	fmt.Println("✅ Programming tools initialized")

	codingAgent := initializeAgent(llmProvider, tools)
	fmt.Println("✅ Coding agent initialized")

	// Create API server
	server := api.NewCodingAgentServer(api.ServerConfig{
		Port:             config.Port,
		Agent:            codingAgent,
		LanguageRegistry: languageRegistry,
		ReadTimeout:      30 * time.Second,
		WriteTimeout:     60 * time.Second,
		MaxRequestSize:   1024 * 1024, // 1MB
	})

	// Start server
	fmt.Printf("\n🚀 Server starting on port %s\n", config.Port)
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  POST /api/v1/coding/ask       - Ask programming question")
	fmt.Println("  POST /api/v1/coding/analyze   - Analyze code")
	fmt.Println("  POST /api/v1/coding/execute   - Execute code")
	fmt.Println("  POST /api/v1/coding/debug     - Debug code")
	fmt.Println("  GET  /api/v1/health           - Health check")
	fmt.Println("  GET  /api/v1/languages        - List supported languages")
	fmt.Println()

	// Graceful shutdown
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server stopped gracefully")
}

// Config holds server configuration
type Config struct {
	Port      string
	LLMAPIKey string
	LLMModel  string
}

// loadConfig loads configuration from environment
func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ OPENAI_API_KEY environment variable not set")
	}

	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "gpt-4"
	}

	return Config{
		Port:      port,
		LLMAPIKey: apiKey,
		LLMModel:  model,
	}
}

// initializeLLM creates the LLM provider
func initializeLLM(config Config) (*llm.ProviderManager, error) {
	providerManager := llm.NewProviderManager()

	openAIProvider, err := llm.NewOpenAIProvider(llm.OpenAIConfig{
		APIKey: config.LLMAPIKey,
		Model:  config.LLMModel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI provider: %w", err)
	}
	if err := providerManager.Register("openai", openAIProvider); err != nil {
		return nil, fmt.Errorf("failed to register OpenAI provider: %w", err)
	}
	log.Println("Registered OpenAI provider")

	if vertexConfig, ok := llm.GetVertexConfigFromEnv(); ok {
		vertexProvider, err := llm.NewVertexAIProvider(vertexConfig)
		if err != nil {
			log.Printf("Warning: failed to create Vertex AI provider: %v", err)
		} else {
			if err := providerManager.Register("vertex", vertexProvider); err != nil {
				log.Printf("Warning: failed to register Vertex AI provider: %v", err)
			} else {
				log.Println("Registered Vertex AI provider")
			}
		}
	}

	return providerManager, nil
}

// initializeLanguages sets up language support
func initializeLanguages() *common.LanguageRegistry {
	registry := common.NewLanguageRegistry()

	// Register Go language
	goProvider := golang.NewProvider()
	if err := registry.Register(goProvider); err != nil {
		log.Printf("⚠️  Failed to register Go provider: %v", err)
	}

	// TODO: Register other languages as they are implemented
	// pythonProvider := python.NewProvider()
	// registry.Register(pythonProvider)

	return registry
}

// initializeTools creates programming tools
func initializeTools(registry *common.LanguageRegistry) []types.Tool {
	return []types.Tool{
		programming.NewCodeAnalysisTool(registry),
		programming.NewCodeExecutionTool(registry),
		programming.NewDocumentationSearchTool(),
		programming.NewStackOverflowSearchTool(),
		programming.NewGitHubSearchTool(),
	}
}

// initializeAgent creates the coding expert agent
func initializeAgent(llmProvider types.LLMProvider, tools []types.Tool) *agent.CodingExpertAgent {
	return agent.NewCodingExpertAgent(agent.CodingExpertConfig{
		Name:           "CodingExpertAPI",
		Description:    "API-based coding expert agent",
		LLM:            llmProvider,
		Tools:          tools,
		MaxSteps:       5,
		Temperature:    0.7,
		Timeout:        defaultTimeout,
		VerboseLogging: true,
		SupportedLangs: []string{"Go", "Python", "JavaScript", "TypeScript", "Rust"},
	})
}
