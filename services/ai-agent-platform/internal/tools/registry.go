package tools

import (
	"fmt"
	"sync"

	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/errors"
	"github.com/DimaJoyti/go-pro/services/ai-agent-platform/pkg/types"
)

// Registry implements the ToolRegistry interface
type Registry struct {
	tools      map[string]types.Tool
	categories map[string][]string
	mu         sync.RWMutex
}

// NewRegistry creates a new tool registry
func NewRegistry() *Registry {
	return &Registry{
		tools:      make(map[string]types.Tool),
		categories: make(map[string][]string),
	}
}

// Register adds a tool to the registry
func (r *Registry) Register(tool types.Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := tool.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %s already registered", name)
	}

	r.tools[name] = tool
	return nil
}

// RegisterWithCategory adds a tool to the registry with a category
func (r *Registry) RegisterWithCategory(tool types.Tool, category string) error {
	if err := r.Register(tool); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.categories[category] = append(r.categories[category], tool.Name())
	return nil
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (types.Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return nil, errors.NewToolNotFoundError(name)
	}

	return tool, nil
}

// List returns all registered tools
func (r *Registry) List() []types.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]types.Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// Unregister removes a tool from the registry
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[name]; !exists {
		return errors.NewToolNotFoundError(name)
	}

	delete(r.tools, name)

	// Remove from categories
	for category, toolNames := range r.categories {
		for i, toolName := range toolNames {
			if toolName == name {
				r.categories[category] = append(toolNames[:i], toolNames[i+1:]...)
				break
			}
		}
	}

	return nil
}

// GetByCategory returns tools in a specific category
func (r *Registry) GetByCategory(category string) []types.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	toolNames, exists := r.categories[category]
	if !exists {
		return []types.Tool{}
	}

	tools := make([]types.Tool, 0, len(toolNames))
	for _, name := range toolNames {
		if tool, exists := r.tools[name]; exists {
			tools = append(tools, tool)
		}
	}

	return tools
}

// GetCategories returns all categories
func (r *Registry) GetCategories() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]string, 0, len(r.categories))
	for category := range r.categories {
		categories = append(categories, category)
	}

	return categories
}

// Count returns the number of registered tools
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.tools)
}

