// Package kaiju implements a creature generator for giant monster (kaiju) sightings.
// It provides realistic random generation of kaiju creatures with various attributes.
package kaiju

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/pymk/creature-sighting/internal/sighting"
)

// Generator creates random kaiju sightings with predefined sets of names, types, and attributes.
type Generator struct {
	names     []string
	types     []string
	behaviors []string
	sizes     []string
}

// NewGenerator creates a new kaiju generator with predefined creature data.
func NewGenerator() *Generator {
	return &Generator{
		names: []string{
			"Gorgozilla", "Mechataur", "Tsunamius", "Pyroclast",
			"Vortexia", "Thundermaw", "Crystalfang", "Nebulox",
			"Seismodon", "Glacierus", "Plasmoid", "Terracrush",
		},
		types: []string{
			"Aquatic", "Terrestrial", "Aerial", "Subterranean",
			"Amphibious", "Cosmic", "Volcanic", "Arctic",
		},
		behaviors: []string{
			"aggressive", "territorial", "curious", "defensive",
			"migratory", "nocturnal", "predatory", "docile",
		},
		sizes: []string{
			"colossal", "massive", "enormous", "gigantic",
			"titanic", "monstrous", "immense", "gargantuan",
		},
	}
}

// Category returns the creature category this generator handles.
func (g *Generator) Category() string {
	return "kaiju"
}

// Generate creates a random kaiju sighting with randomized attributes and location.
func (g *Generator) Generate() (*sighting.Sighting, error) {
	loc := g.randomLocation()

	name, err := g.randomChoice(g.names)
	if err != nil {
		return nil, fmt.Errorf("failed to generate name: %w", err)
	}

	kaijuType, err := g.randomChoice(g.types)
	if err != nil {
		return nil, fmt.Errorf("failed to generate type: %w", err)
	}

	behavior, err := g.randomChoice(g.behaviors)
	if err != nil {
		return nil, fmt.Errorf("failed to generate behavior: %w", err)
	}

	size, err := g.randomChoice(g.sizes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate size: %w", err)
	}

	height, err := g.randomInt(50, 300)
	if err != nil {
		return nil, fmt.Errorf("failed to generate height: %w", err)
	}

	sighting := &sighting.Sighting{
		ID:          fmt.Sprintf("kaiju-%d", time.Now().UnixNano()),
		Name:        name,
		Type:        kaijuType,
		Category:    g.Category(),
		Location:    loc,
		Description: fmt.Sprintf("A %s %s kaiju displaying %s behavior", size, kaijuType, behavior),
		Timestamp:   time.Now(),
		Attributes: sighting.Attributes{
			"size":     size,
			"behavior": behavior,
			"height":   fmt.Sprintf("%d meters", height),
		},
	}

	return sighting, nil
}

// randomLocation selects a random location from predefined major cities worldwide.
func (g *Generator) randomLocation() sighting.Location {
	locations := []sighting.Location{
		{Latitude: 35.6762, Longitude: 139.6503, City: "Tokyo", Country: "Japan", Region: "Asia"},
		{Latitude: 37.7749, Longitude: -122.4194, City: "San Francisco", Country: "USA", Region: "North America"},
		{Latitude: -33.8688, Longitude: 151.2093, City: "Sydney", Country: "Australia", Region: "Oceania"},
		{Latitude: 51.5074, Longitude: -0.1278, City: "London", Country: "UK", Region: "Europe"},
		{Latitude: -22.9068, Longitude: -43.1729, City: "Rio de Janeiro", Country: "Brazil", Region: "South America"},
		{Latitude: 40.7128, Longitude: -74.0060, City: "New York", Country: "USA", Region: "North America"},
		{Latitude: 1.3521, Longitude: 103.8198, City: "Singapore", Country: "Singapore", Region: "Asia"},
		{Latitude: 64.1466, Longitude: -21.9426, City: "Reykjavik", Country: "Iceland", Region: "Europe"},
		{Latitude: -1.2921, Longitude: 36.8219, City: "Nairobi", Country: "Kenya", Region: "Africa"},
		{Latitude: 19.4326, Longitude: -99.1332, City: "Mexico City", Country: "Mexico", Region: "North America"},
	}

	idx, _ := g.randomInt(0, len(locations)-1)
	return locations[idx]
}

// randomChoice selects a random string from the provided choices slice.
func (g *Generator) randomChoice(choices []string) (string, error) {
	if len(choices) == 0 {
		return "", fmt.Errorf("empty choices")
	}

	idx, err := g.randomInt(0, len(choices)-1)
	if err != nil {
		return "", err
	}

	return choices[idx], nil
}

// randomInt generates a cryptographically secure random integer in the range [min, max].
// Uses crypto/rand for security rather than math/rand for unpredictability.
func (g *Generator) randomInt(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("min cannot be greater than max")
	}

	// Use crypto/rand for security - prevents predictable sequences
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}
