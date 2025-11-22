package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/luizeduardocarvalho/labflux-functions/docs"
	"github.com/luizeduardocarvalho/labflux-functions/internal"
	"github.com/luizeduardocarvalho/labflux-functions/pkg/database"
)

func main() {
	// Load .env file (ignore error if not found - env vars may be set externally)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	// Setup router
	router := internal.SetupRouter()

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting LabFlux API server on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
