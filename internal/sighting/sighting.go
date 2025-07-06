// Package sighting provides core domain types and interfaces for creature sightings.
// It defines the structure of sightings, locations, and the generator interface
// used by creature-specific implementations.
package sighting

import (
	"time"
)

// Sighting represents a fictional creature sighting with all relevant details.
// It includes identification, classification, location, and custom attributes.
type Sighting struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Category    string     `json:"category"`
	Location    Location   `json:"location"`
	Description string     `json:"description"`
	Timestamp   time.Time  `json:"timestamp"`
	Attributes  Attributes `json:"attributes"`
}

// Location represents the geographic location of a sighting.
// It includes coordinates and optional human-readable location details.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
}

// Attributes provides flexible key-value storage for creature-specific properties.
// This allows different creature types to store custom data without schema changes.
type Attributes map[string]any

// Generator defines the interface for creature sighting generators.
// Implementations must be able to generate random sightings and identify their category.
type Generator interface {
	Generate() (*Sighting, error)
	Category() string
}
