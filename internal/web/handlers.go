// Package web provides HTTP handlers for the web UI interface.
// It handles HTML responses using templ templates and manages sighting display.
package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pymk/creature-sighting/internal/sighting"
	"github.com/pymk/creature-sighting/internal/storage"
	"github.com/pymk/creature-sighting/internal/templates"
)

// Handler provides HTTP handlers for web UI endpoints.
type Handler struct {
	registry *sighting.Registry
	storage  *storage.InMemoryStorage
}

// NewHandler creates a new web handler with the given registry and storage.
func NewHandler(registry *sighting.Registry, storage *storage.InMemoryStorage) *Handler {
	return &Handler{
		registry: registry,
		storage:  storage,
	}
}

// HandleHome renders the home page template.
func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.Home().Render(r.Context(), w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleSightings renders the sightings list page with optional category filtering.
// Accepts "category" query parameter to filter sightings by category.
func (h *Handler) HandleSightings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	var sightings []sighting.Sighting

	if category != "" {
		sightings = h.storage.GetByCategory(category)
	} else {
		sightings = h.storage.GetAll()
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.SightingsList(sightings).Render(r.Context(), w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleSightingDetail renders the detail page for a specific sighting.
// Extracts sighting ID from URL path and handles special "random" case.
func (h *Handler) HandleSightingDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path (e.g., /sighting/kaiju-123 -> kaiju-123)
	path := strings.TrimPrefix(r.URL.Path, "/sighting/")
	if path == "" {
		http.Error(w, "Sighting ID required", http.StatusBadRequest)
		return
	}

	// Handle special case for generating new random sighting
	if path == "random" {
		h.HandleRandomSighting(w, r)
		return
	}

	sighting, exists := h.storage.Get(path)
	if !exists {
		http.Error(w, "Sighting not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.SightingDetail(sighting).Render(r.Context(), w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleRandomSighting generates a new random sighting and redirects to its detail page.
// Accepts optional "category" query parameter, defaults to "kaiju".
func (h *Handler) HandleRandomSighting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get category from query parameter, default to "kaiju"
	category := r.URL.Query().Get("category")
	if category == "" {
		category = "kaiju"
	}

	generator, err := h.registry.Get(category)
	if err != nil {
		http.Error(w, fmt.Sprintf("Category not found: %s", category), http.StatusNotFound)
		return
	}

	sighting, err := generator.Generate()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate sighting: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the generated sighting
	h.storage.Add(*sighting)

	// Redirect to the sighting detail page
	http.Redirect(w, r, fmt.Sprintf("/sighting/%s", sighting.ID), http.StatusSeeOther)
}

// HandleLocations renders the locations list page.
func (h *Handler) HandleLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get unique locations from all sightings
	allSightings := h.storage.GetAll()
	locationMap := make(map[string]sighting.Location)

	for _, s := range allSightings {
		key := fmt.Sprintf("%s,%s", s.Location.City, s.Location.Country)
		locationMap[key] = s.Location
	}

	// Convert map to slice
	locations := make([]sighting.Location, 0, len(locationMap))
	for _, loc := range locationMap {
		locations = append(locations, loc)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.LocationsList(locations).Render(r.Context(), w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleCategories renders the categories list page.
func (h *Handler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories := h.registry.Categories()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.CategoriesList(categories).Render(r.Context(), w); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
