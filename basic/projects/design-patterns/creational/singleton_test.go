package creational

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingletonInstance(t *testing.T) {
	// Get two instances
	db1 := GetDatabase()
	db2 := GetDatabase()

	// Should be the same instance
	assert.Equal(t, db1, db2, "Both instances should be the same")
	assert.Same(t, db1, db2, "Should point to the same memory address")
}

func TestSingletonConcurrency(t *testing.T) {
	const goroutines = 100
	var wg sync.WaitGroup
	instances := make([]*Database, goroutines)

	// Create instances concurrently
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			instances[index] = GetDatabase()
		}(i)
	}

	wg.Wait()

	// All instances should be the same
	firstInstance := instances[0]
	for i := 1; i < goroutines; i++ {
		assert.Same(t, firstInstance, instances[i], "All instances should be the same")
	}
}

func TestDatabaseConnect(t *testing.T) {
	db := GetDatabase()
	initialConnections := db.GetConnections()

	db.Connect()
	assert.Equal(t, initialConnections+1, db.GetConnections())

	db.Connect()
	assert.Equal(t, initialConnections+2, db.GetConnections())
}

func TestConfigManager(t *testing.T) {
	config1 := GetConfig()
	config2 := GetConfig()

	// Should be the same instance
	assert.Same(t, config1, config2)

	// Test default settings
	appName, ok := config1.Get("app_name")
	assert.True(t, ok)
	assert.Equal(t, "DesignPatterns", appName)

	// Test set and get
	config1.Set("test_key", "test_value")
	value, ok := config2.Get("test_key")
	assert.True(t, ok)
	assert.Equal(t, "test_value", value)
}

func TestConfigManagerConcurrency(t *testing.T) {
	config := GetConfig()
	const goroutines = 50
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			config.Set("key", "value")
		}(i)
	}

	// Concurrent reads
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			config.Get("key")
		}()
	}

	wg.Wait()
}

func BenchmarkGetDatabase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetDatabase()
	}
}

func BenchmarkGetConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetConfig()
	}
}

