package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	LLM       LLMConfig
	Auth      AuthConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type LLMConfig struct {
	Providers map[string]ProviderConfig
	Default   string
}

type ProviderConfig struct {
	Provider    string
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float64
}

type AuthConfig struct {
	Enabled   bool
	JWTSecret string
}

type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	Burst             int
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		},
		LLM: LLMConfig{
			Default: getEnv("LLM_DEFAULT_PROVIDER", "openai"),
			Providers: map[string]ProviderConfig{
				"openai": {
					Provider:    "openai",
					APIKey:      getEnv("OPENAI_API_KEY", ""),
					BaseURL:     getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
					Model:       getEnv("OPENAI_MODEL", "gpt-4"),
					MaxTokens:   getEnvInt("OPENAI_MAX_TOKENS", 2048),
					Temperature: getEnvFloat("OPENAI_TEMPERATURE", 0.7),
				},
				"anthropic": {
					Provider:    "anthropic",
					APIKey:      getEnv("ANTHROPIC_API_KEY", ""),
					BaseURL:     getEnv("ANTHROPIC_BASE_URL", "https://api.anthropic.com/v1"),
					Model:       getEnv("ANTHROPIC_MODEL", "claude-3-sonnet-20240229"),
					MaxTokens:   getEnvInt("ANTHROPIC_MAX_TOKENS", 4096),
					Temperature: getEnvFloat("ANTHROPIC_TEMPERATURE", 0.7),
				},
				"azure": {
					Provider:    "azure",
					APIKey:      getEnv("AZURE_OPENAI_API_KEY", ""),
					BaseURL:     getEnv("AZURE_OPENAI_ENDPOINT", ""),
					Model:       getEnv("AZURE_OPENAI_DEPLOYMENT", "gpt-4"),
					MaxTokens:   getEnvInt("AZURE_OPENAI_MAX_TOKENS", 2048),
					Temperature: getEnvFloat("AZURE_OPENAI_TEMPERATURE", 0.7),
				},
				"ollama": {
					Provider:    "ollama",
					APIKey:      "",
					BaseURL:     getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
					Model:       getEnv("OLLAMA_MODEL", "llama2"),
					MaxTokens:   getEnvInt("OLLAMA_MAX_TOKENS", 2048),
					Temperature: getEnvFloat("OLLAMA_TEMPERATURE", 0.7),
				},
			},
		},
		Auth: AuthConfig{
			Enabled:   getEnvBool("AUTH_ENABLED", false),
			JWTSecret: getEnv("JWT_SECRET", "secret"),
		},
		RateLimit: RateLimitConfig{
			Enabled:           getEnvBool("RATE_LIMIT_ENABLED", true),
			RequestsPerMinute: getEnvInt("RATE_LIMIT_RPM", 60),
			Burst:             getEnvInt("RATE_LIMIT_BURST", 10),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
