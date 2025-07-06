// Package api provides HTTP REST endpoints for generating creature sightings.
// It handles JSON responses and error handling for the API layer.
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pymk/creature-sighting/internal/sighting"
)

// Handler provides HTTP handlers for API endpoints.
type Handler struct {
	registry *sighting.Registry
}

// NewHandler creates a new API handler with the given registry.
func NewHandler(registry *sighting.Registry) *Handler {
	return &Handler{
		registry: registry,
	}
}

// HandleSighting generates and returns a random sighting via GET /api/sighting.
// Accepts optional "category" query parameter, defaults to "kaiju".
func (h *Handler) HandleSighting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "kaiju"
	}

	generator, err := h.registry.Get(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sighting, err := generator.Generate()
	if err != nil {
		log.Printf("Error generating sighting: %v", err)
		http.Error(w, "Failed to generate sighting", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sighting); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// HandleCategories returns all available creature categories via GET /api/categories.
// Returns JSON with "categories" array containing all registered category names.
func (h *Handler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories := h.registry.Categories()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]string{
		"categories": categories,
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
