package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadConfig attempts to load configuration from a .env file.
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found: %v. Proceeding with system environment variables.", err)
	} else {
		log.Println("System: Configuration loaded from environment.")
	}
}
