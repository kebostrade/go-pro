package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load valid config", func(t *testing.T) {
		// Create a temporary config file
		content := `
greeting:
  default_name: "TestUser"
  default_times: 3
server:
  port: "9090"
  host: "0.0.0.0"
`
		tmpFile, err := os.CreateTemp("", "config-*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(content)
		assert.NoError(t, err)
		tmpFile.Close()

		cfg, err := LoadConfig(tmpFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "TestUser", cfg.Greeting.DefaultName)
		assert.Equal(t, 3, cfg.Greeting.DefaultTimes)
		assert.Equal(t, "9090", cfg.Server.Port)
		assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	})

	t.Run("file not found", func(t *testing.T) {
		cfg, err := LoadConfig("nonexistent.yaml")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "config-*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("invalid: yaml: content:")
		assert.NoError(t, err)
		tmpFile.Close()

		cfg, err := LoadConfig(tmpFile.Name())
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}
