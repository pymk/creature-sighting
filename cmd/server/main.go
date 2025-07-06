// Package main provides the HTTP server entry point for the Creature Sighting application.
// Creature Sighting generates and displays fictional creature sightings.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pymk/creature-sighting/internal/api"
	"github.com/pymk/creature-sighting/internal/creatures/kaiju"
	"github.com/pymk/creature-sighting/internal/sighting"
	"github.com/pymk/creature-sighting/internal/storage"
	"github.com/pymk/creature-sighting/internal/web"
)

// main is the application entry point that starts the HTTP server.
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run initializes and starts the HTTP server with graceful shutdown handling.
// It sets up creature generators, storage, handlers, and routes before starting the server.
func run() error {
	registry := sighting.NewRegistry()

	kaijuGen := kaiju.NewGenerator()
	if err := registry.Register("kaiju", kaijuGen); err != nil {
		return err
	}

	// Initialize storage with some initial sightings
	store := storage.NewInMemoryStorage()
	store.GenerateInitialSightings(registry)

	// API handlers
	apiHandler := api.NewHandler(registry)

	// Web handlers
	webHandler := web.NewHandler(registry, store)

	mux := http.NewServeMux()

	// Web routes
	mux.HandleFunc("/", webHandler.HandleHome)
	mux.HandleFunc("/sightings", webHandler.HandleSightings)
	mux.HandleFunc("/sighting/", webHandler.HandleSightingDetail)
	mux.HandleFunc("/locations", webHandler.HandleLocations)
	mux.HandleFunc("/categories", webHandler.HandleCategories)

	// API routes
	mux.HandleFunc("/api/sighting", apiHandler.HandleSighting)
	mux.HandleFunc("/api/categories", apiHandler.HandleCategories)

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server starting on %s", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either server error or shutdown signal
	select {
	case err := <-serverErr:
		// Server failed to start or crashed
		return err
	case sig := <-sigChan:
		// Received shutdown signal - perform graceful shutdown
		log.Printf("Received signal: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
}
