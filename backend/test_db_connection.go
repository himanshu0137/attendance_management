package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/db"
)

func main() {
	fmt.Println("=== Testing Database Connection ===")
	
	// Load configuration
	cfg := config.Load()
	fmt.Printf("Loaded configuration:\n")
	fmt.Printf("  HTTP Port: %d\n", cfg.HTTPPort)
	fmt.Printf("  Database URL: %s\n", cfg.DatabaseURL)
	fmt.Printf("  Keycloak Issuer: %s\n", cfg.KeycloakIssuerURL)
	
	// Test if database URL is set
	if cfg.DatabaseURL == "" {
		fmt.Println("ERROR: Database URL is empty! Please check your environment variables.")
		os.Exit(1)
	}
	
	// Try to connect to database
	ctx := context.Background()
	fmt.Printf("\nAttempting to connect to database at: %s\n", cfg.DatabaseURL)
	
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Printf("ERROR: Failed to connect to database: %v\n", err)
		log.Fatal("Database connection failed")
	}
	
	defer pool.Close()
	
	// Test ping
	err = pool.Ping(ctx)
	if err != nil {
		fmt.Printf("ERROR: Failed to ping database: %v\n", err)
		log.Fatal("Database ping failed")
	}
	
	fmt.Println("SUCCESS: Database connection established!")
	fmt.Println("The system is ready to run with proper database connectivity.")
}