// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package container provides functionality for the GO-PRO Learning Platform.
package container

import (
	"context"
	"fmt"
	"time"

	"go-pro-backend/internal/agents"
	"go-pro-backend/internal/cache"
	"go-pro-backend/internal/config"
	"go-pro-backend/internal/executor"
	"go-pro-backend/internal/messaging"
	"go-pro-backend/internal/repository"
	"go-pro-backend/internal/repository/postgres"
	"go-pro-backend/internal/service"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// Container holds all application dependencies.
type Container struct {
	// Configuration.
	Config *config.Config

	// Infrastructure.
	Logger    logger.Logger
	Validator validator.Validator
	Cache     cache.CacheManager
	Messaging *messaging.Service

	// Repositories.
	Repositories *repository.Repositories

	// Services.
	Services *service.Services

	// AI Agent Pool for intelligent code analysis and assistance.
	AgentPool *agents.AgentPool

	// Lifecycle management.
	shutdownFuncs []func() error
}

// ContainerConfig holds configuration for container initialization.
type ContainerConfig struct {
	Config        *config.Config
	Logger        logger.Logger
	SkipMessaging bool
	SkipCache     bool
}

// NewContainer creates a new dependency injection container.
func NewContainer(cfg *ContainerConfig) (*Container, error) {
	if cfg == nil {
		return nil, fmt.Errorf("container config is required")
	}

	container := &Container{
		Config:        cfg.Config,
		Logger:        cfg.Logger,
		shutdownFuncs: make([]func() error, 0),
	}

	// Initialize components in dependency order.
	if err := container.initializeValidator(); err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	if !cfg.SkipCache {
		if err := container.initializeCache(); err != nil {
			return nil, fmt.Errorf("failed to initialize cache: %w", err)
		}
	}

	if !cfg.SkipMessaging {
		if err := container.initializeMessaging(); err != nil {
			return nil, fmt.Errorf("failed to initialize messaging: %w", err)
		}
	}

	if err := container.initializeRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	if err := container.initializeServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	return container, nil
}

// initializeValidator initializes the validator.
func (c *Container) initializeValidator() error {
	c.Validator = validator.New()
	return nil
}

// initializeCache initializes the cache manager.
func (c *Container) initializeCache() error {
	cacheManager, err := cache.NewManagerFromEnv()
	if err != nil {
		c.Logger.Warn(context.Background(), "Failed to initialize cache manager, using no-op cache", "error", err)
		// Use a no-op cache implementation.
		c.Cache = &NoOpCacheManager{}

		return nil
	}

	c.Cache = cacheManager
	c.addShutdownFunc(cacheManager.Close)

	c.Logger.Info(context.Background(), "Cache manager initialized successfully")

	return nil
}

// initializeMessaging initializes the messaging service.
func (c *Container) initializeMessaging() error {
	messagingConfig := messaging.LoadConfigFromEnv()
	messagingService, err := messaging.NewService(messagingConfig)
	if err != nil {
		c.Logger.Warn(context.Background(), "Failed to initialize messaging service", "error", err)
		// Use a no-op messaging service.
		c.Messaging = &messaging.Service{}

		return nil
	}

	c.Messaging = messagingService
	c.addShutdownFunc(messagingService.Close)

	// Start consumer in background if enabled.
	if messagingService.IsEnabled() {
		go func() {
			if err := messagingService.StartConsumer(context.Background()); err != nil {
				c.Logger.Error(context.Background(), "Messaging consumer failed", "error", err)
			}
		}()
		c.Logger.Info(context.Background(), "Messaging service initialized successfully")
	}

	return nil
}

// initializeRepositories initializes the repository layer.
func (c *Container) initializeRepositories() error {
	if c.Config.Database.Driver == "postgres" {
		pgRepos, err := postgres.NewRepositoriesFromEnv()
		if err != nil {
			c.Logger.Error(context.Background(), "Failed to initialize PostgreSQL repositories, falling back to memory", "error", err)
			c.Repositories = repository.NewRepositoriesSimple()

			return nil
		}

		// PostgreSQL repositories initialized successfully.

		c.Repositories = &repository.Repositories{
			Course:   pgRepos.Course,
			Lesson:   pgRepos.Lesson,
			Exercise: pgRepos.Exercise,
			Progress: pgRepos.Progress,
		}

		c.addShutdownFunc(func() error {
			return pgRepos.Close()
		})

		c.Logger.Info(context.Background(), "PostgreSQL repositories initialized successfully")
	} else {
		c.Repositories = repository.NewRepositoriesSimple()
		c.Logger.Info(context.Background(), "Using in-memory repositories")
	}

	return nil
}

// initializeServices initializes the service layer.
func (c *Container) initializeServices() error {
	serviceConfig := &service.Config{
		Logger:    c.Logger,
		Validator: c.Validator,
		Cache:     c.Cache,
		Messaging: c.Messaging,
	}

	services, err := service.NewServices(c.Repositories, serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create services: %w", err)
	}

	// Replace mock executor with real Docker executor for code execution
	services.Executor = executor.NewDockerExecutor()
	c.Logger.Info(context.Background(), "Docker executor initialized for code execution")

	c.Services = services
	c.Logger.Info(context.Background(), "Services initialized successfully")

	// Initialize AI Agent Pool for intelligent code analysis
	c.AgentPool = agents.NewAgentPool()
	c.Logger.Info(context.Background(), "AI Agent Pool initialized for playground")

	return nil
}

// addShutdownFunc adds a function to be called during shutdown.
func (c *Container) addShutdownFunc(fn func() error) {
	c.shutdownFuncs = append(c.shutdownFuncs, fn)
}

// Shutdown gracefully shuts down all components.
func (c *Container) Shutdown(ctx context.Context) error {
	c.Logger.Info(ctx, "Shutting down container...")

	var lastErr error
	for i := len(c.shutdownFuncs) - 1; i >= 0; i-- {
		if err := c.shutdownFuncs[i](); err != nil {
			c.Logger.Error(ctx, "Error during shutdown", "error", err)
			lastErr = err
		}
	}

	c.Logger.Info(ctx, "Container shutdown completed")

	return lastErr
}

// HealthCheck performs health checks on all components.
func (c *Container) HealthCheck(ctx context.Context) error {
	// Check cache.
	if c.Cache != nil {
		if err := c.Cache.HealthCheck(ctx); err != nil {
			return fmt.Errorf("cache health check failed: %w", err)
		}
	}

	// Check repositories - basic check that they exist.
	if c.Repositories == nil {
		return fmt.Errorf("repositories not initialized")
	}

	// Check messaging - basic check that it's enabled if configured.
	if c.Messaging != nil && c.Messaging.IsEnabled() {
		// Messaging is enabled and available.
	}

	return nil
}

// GetStats returns statistics about the container components.
func (c *Container) GetStats(ctx context.Context) map[string]interface{} {
	stats := make(map[string]interface{})

	// Cache stats - basic info.
	if c.Cache != nil {
		stats["cache"] = map[string]interface{}{
			"enabled": true,
			"type":    "redis", // or "no-op" for NoOpCacheManager
		}
	}

	// Add more stats as needed.
	stats["components"] = map[string]bool{
		"cache":     c.Cache != nil,
		"messaging": c.Messaging != nil && c.Messaging.IsEnabled(),
		"database":  c.Repositories != nil,
	}

	return stats
}

// NoOpCacheManager is a no-op implementation of CacheManager.
type NoOpCacheManager struct{}

// Cache interface implementation.
func (n *NoOpCacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (n *NoOpCacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	return cache.ErrCacheMiss
}

func (n *NoOpCacheManager) Delete(ctx context.Context, key string) error {
	return nil
}

func (n *NoOpCacheManager) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (n *NoOpCacheManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (n *NoOpCacheManager) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

func (n *NoOpCacheManager) Increment(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (n *NoOpCacheManager) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (n *NoOpCacheManager) Decrement(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (n *NoOpCacheManager) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (n *NoOpCacheManager) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return false, nil
}

func (n *NoOpCacheManager) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	return "", nil
}

func (n *NoOpCacheManager) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return make([]interface{}, len(keys)), nil
}

func (n *NoOpCacheManager) MSet(ctx context.Context, pairs ...interface{}) error {
	return nil
}

func (n *NoOpCacheManager) Keys(ctx context.Context, pattern string) ([]string, error) {
	return []string{}, nil
}

func (n *NoOpCacheManager) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return []string{}, 0, nil
}

func (n *NoOpCacheManager) FlushDB(ctx context.Context) error {
	return nil
}

func (n *NoOpCacheManager) HealthCheck(ctx context.Context) error {
	return nil
}

func (n *NoOpCacheManager) Close() error {
	return nil
}

// SessionStore interface implementation.
func (n *NoOpCacheManager) CreateSession(ctx context.Context, sessionID string, data map[string]interface{}, expiration time.Duration) error {
	return nil
}

func (n *NoOpCacheManager) GetSession(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	return nil, cache.ErrCacheNotFound
}

func (n *NoOpCacheManager) UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
	return nil
}

func (n *NoOpCacheManager) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}

func (n *NoOpCacheManager) RefreshSession(ctx context.Context, sessionID string, expiration time.Duration) error {
	return nil
}

func (n *NoOpCacheManager) ListUserSessions(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}

func (n *NoOpCacheManager) DeleteUserSessions(ctx context.Context, userID string) error {
	return nil
}

func (n *NoOpCacheManager) CleanupExpiredSessions(ctx context.Context) error {
	return nil
}

// DistributedLock interface implementation.
func (n *NoOpCacheManager) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return false, nil
}

func (n *NoOpCacheManager) Unlock(ctx context.Context, key string) error {
	return nil
}

func (n *NoOpCacheManager) Extend(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}

func (n *NoOpCacheManager) IsLocked(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (n *NoOpCacheManager) GetLockTTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

// RateLimiter interface implementation.
func (n *NoOpCacheManager) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	return true, nil
}

func (n *NoOpCacheManager) AllowN(ctx context.Context, key string, count, limit int64, window time.Duration) (bool, error) {
	return true, nil
}

func (n *NoOpCacheManager) Remaining(ctx context.Context, key string, limit int64, window time.Duration) (int64, error) {
	return limit, nil
}

func (n *NoOpCacheManager) Reset(ctx context.Context, key string) error {
	return nil
}

// PubSub interface implementation.
func (n *NoOpCacheManager) Publish(ctx context.Context, channel string, message interface{}) error {
	return nil
}

func (n *NoOpCacheManager) Subscribe(ctx context.Context, channels ...string) (<-chan cache.Message, error) {
	ch := make(chan cache.Message)
	close(ch) // Return closed channel

	return ch, nil
}

func (n *NoOpCacheManager) Unsubscribe(ctx context.Context, channels ...string) error {
	return nil
}

func (n *NoOpCacheManager) PSubscribe(ctx context.Context, patterns ...string) (<-chan cache.Message, error) {
	ch := make(chan cache.Message)
	close(ch) // Return closed channel

	return ch, nil
}

func (n *NoOpCacheManager) PUnsubscribe(ctx context.Context, patterns ...string) error {
	return nil
}
