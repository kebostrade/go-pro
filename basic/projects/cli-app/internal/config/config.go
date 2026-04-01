package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Greeting GreetingConfig `yaml:"greeting"`
	Server   ServerConfig   `yaml:"server"`
}

// GreetingConfig holds greeting-related settings
type GreetingConfig struct {
	DefaultName  string `yaml:"default_name"`
	DefaultTimes int    `yaml:"default_times"`
}

// ServerConfig holds server-related settings
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.Greeting.DefaultName == "" {
		cfg.Greeting.DefaultName = "World"
	}
	if cfg.Greeting.DefaultTimes == 0 {
		cfg.Greeting.DefaultTimes = 1
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}

	return &cfg, nil
}
