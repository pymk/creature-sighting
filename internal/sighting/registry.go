package sighting

import (
	"fmt"
	"sync"
)

// Registry manages thread-safe registration and retrieval of creature generators.
// It maintains a map of category names to their corresponding generators.
type Registry struct {
	mu         sync.RWMutex
	generators map[string]Generator
}

// NewRegistry creates a new empty registry for creature generators.
func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]Generator),
	}
}

// Register adds a generator for a specific category to the registry.
// It returns an error if a generator for the category is already registered.
func (r *Registry) Register(category string, generator Generator) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.generators[category]; exists {
		return fmt.Errorf("generator for category %s already registered", category)
	}

	r.generators[category] = generator
	return nil
}

// Get retrieves a generator for the specified category.
// It returns an error if no generator is found for the category.
func (r *Registry) Get(category string) (Generator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	generator, exists := r.generators[category]
	if !exists {
		return nil, fmt.Errorf("no generator found for category %s", category)
	}

	return generator, nil
}

// Categories returns a list of all registered category names.
func (r *Registry) Categories() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]string, 0, len(r.generators))
	for category := range r.generators {
		categories = append(categories, category)
	}

	return categories
}
