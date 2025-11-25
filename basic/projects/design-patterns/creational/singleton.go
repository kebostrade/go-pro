package creational

import (
	"fmt"
	"sync"
)

/*
SINGLETON PATTERN

Purpose: Ensure a class has only one instance and provide a global point of access to it.

Use Cases:
- Database connections
- Configuration managers
- Logger instances
- Cache managers

Go-Specific Implementation:
- Use sync.Once for thread-safe lazy initialization
- Private struct with public accessor function
*/

// Database represents a singleton database connection
type Database struct {
	ConnectionString string
	mu               sync.Mutex
	connections      int
}

var (
	dbInstance *Database
	once       sync.Once
)

// GetDatabase returns the singleton database instance
// Thread-safe using sync.Once
func GetDatabase() *Database {
	once.Do(func() {
		fmt.Println("Creating database instance...")
		dbInstance = &Database{
			ConnectionString: "localhost:5432",
			connections:      0,
		}
	})
	return dbInstance
}

// Connect simulates a database connection
func (db *Database) Connect() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.connections++
	fmt.Printf("Connected to database. Total connections: %d\n", db.connections)
}

// GetConnections returns the number of connections
func (db *Database) GetConnections() int {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.connections
}

// Alternative: Eager Initialization Singleton
// Initialized at package load time

var eagerInstance = &Database{
	ConnectionString: "localhost:5432",
	connections:      0,
}

// GetEagerDatabase returns the eagerly initialized singleton
func GetEagerDatabase() *Database {
	return eagerInstance
}

// ConfigManager is another singleton example for configuration
type ConfigManager struct {
	settings map[string]string
	mu       sync.RWMutex
}

var (
	configInstance *ConfigManager
	configOnce     sync.Once
)

// GetConfig returns the singleton configuration manager
func GetConfig() *ConfigManager {
	configOnce.Do(func() {
		configInstance = &ConfigManager{
			settings: make(map[string]string),
		}
		// Load default settings
		configInstance.settings["app_name"] = "DesignPatterns"
		configInstance.settings["version"] = "1.0.0"
	})
	return configInstance
}

// Get retrieves a configuration value
func (c *ConfigManager) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.settings[key]
	return val, ok
}

// Set updates a configuration value
func (c *ConfigManager) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.settings[key] = value
}

// GetAll returns all configuration settings
func (c *ConfigManager) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// Return a copy to prevent external modification
	copy := make(map[string]string, len(c.settings))
	for k, v := range c.settings {
		copy[k] = v
	}
	return copy
}

