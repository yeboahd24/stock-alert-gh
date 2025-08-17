package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"shares-alert-backend/internal/app"
	"shares-alert-backend/internal/config"
)

func main() {
	// Load environment variables from .env file
	// Try multiple paths to find the .env file
	envPaths := []string{".env", "../.env", "../../.env"}
	var envLoaded bool
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded environment variables from: %s", path)
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		log.Printf("Warning: .env file not found in any of the expected locations")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize and start the application
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Starting server on port %s", port)
	if err := application.Start(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}