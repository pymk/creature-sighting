// Package storage provides in-memory storage for creature sightings.
// It maintains thread-safe access and preserves insertion order for chronological display.
package storage

import (
	"sync"
	"time"

	"github.com/pymk/creature-sighting/internal/sighting"
)

// InMemoryStorage provides thread-safe in-memory storage for sightings.
// It maintains both a map for fast lookups and a slice for insertion order.
type InMemoryStorage struct {
	mu        sync.RWMutex
	sightings map[string]sighting.Sighting
	order     []string // maintain insertion order
}

// NewInMemoryStorage creates a new empty in-memory storage instance.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		sightings: make(map[string]sighting.Sighting),
		order:     make([]string, 0),
	}
}

// Add stores a sighting in the storage, maintaining insertion order.
func (s *InMemoryStorage) Add(sighting sighting.Sighting) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sightings[sighting.ID] = sighting
	s.order = append(s.order, sighting.ID)
}

// Get retrieves a sighting by ID, returning the sighting and whether it exists.
func (s *InMemoryStorage) Get(id string) (sighting.Sighting, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sighting, exists := s.sightings[id]
	return sighting, exists
}

// GetAll returns all sightings in reverse chronological order (most recent first).
func (s *InMemoryStorage) GetAll() []sighting.Sighting {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]sighting.Sighting, 0, len(s.order))
	// Return in reverse order (most recent first)
	for i := len(s.order) - 1; i >= 0; i-- {
		id := s.order[i]
		if sighting, exists := s.sightings[id]; exists {
			result = append(result, sighting)
		}
	}
	return result
}

// GetByCategory returns all sightings for a specific category in reverse chronological order.
func (s *InMemoryStorage) GetByCategory(category string) []sighting.Sighting {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]sighting.Sighting, 0)
	// Return in reverse order (most recent first)
	for i := len(s.order) - 1; i >= 0; i-- {
		id := s.order[i]
		if sighting, exists := s.sightings[id]; exists && sighting.Category == category {
			result = append(result, sighting)
		}
	}
	return result
}

// Count returns the total number of stored sightings.
func (s *InMemoryStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.sightings)
}

// Clear removes all sightings from storage.
func (s *InMemoryStorage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sightings = make(map[string]sighting.Sighting)
	s.order = make([]string, 0)
}

// GenerateInitialSightings creates demo sightings spread across recent days.
// This populates the storage with sample data for demonstration purposes.
func (s *InMemoryStorage) GenerateInitialSightings(registry *sighting.Registry) {
	// Generate 5 initial kaiju sightings
	for i := 0; i < 5; i++ {
		if gen, err := registry.Get("kaiju"); err == nil {
			sighting, err := gen.Generate()
			if err != nil {
				continue
			}
			// Spread timestamps across last few days (6 hours apart)
			sighting.Timestamp = time.Now().Add(-time.Duration(i*6) * time.Hour)
			s.Add(*sighting)
		}
	}
}
