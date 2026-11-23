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
	fmt.Println("=== Debug Database Connection ===")
	
	// Print all environment variables for debugging
	fmt.Println("Environment Variables:")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
	
	// Load configuration
	cfg := config.Load()
	
	fmt.Printf("\nConfiguration Loaded:\n")
	fmt.Printf("  HTTP Port: %d\n", cfg.HTTPPort)
	fmt.Printf("  Database URL: '%s'\n", cfg.DatabaseURL)
	fmt.Printf("  Keycloak Issuer: %s\n", cfg.KeycloakIssuerURL)
	fmt.Printf("  Keycloak Client ID: %s\n", cfg.KeycloakClientID)
	fmt.Printf("  Log Level: %s\n", cfg.LogLevel)
	
	// Check if database URL is properly formatted
	if cfg.DatabaseURL == "" {
		fmt.Println("\n❌ ERROR: Database URL is empty!")
		fmt.Println("Please set the ATTENDANCE_DB environment variable correctly.")
		fmt.Println("Example: export ATTENDANCE_DB='postgresql://attendance:attendance@localhost:5432/attendance'")
		os.Exit(1)
	}
	
	// Check if database URL contains expected components
	if cfg.DatabaseURL != "" {
		fmt.Printf("\n✅ Database URL detected: %s\n", cfg.DatabaseURL)
		
		// Test connection with more verbose error handling
		ctx := context.Background()
		fmt.Printf("Attempting to connect to database...\n")
		
		pool, err := db.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			fmt.Printf("❌ ERROR: Failed to connect to database: %v\n", err)
			log.Fatal("Database connection failed - please check your configuration")
		}
		
		defer pool.Close()
		
		// Test ping
		err = pool.Ping(ctx)
		if err != nil {
			fmt.Printf("❌ ERROR: Failed to ping database: %v\n", err)
			log.Fatal("Database ping failed")
		}
		
		fmt.Println("✅ SUCCESS: Database connection established!")
	} else {
		fmt.Println("❌ No database URL configured.")
		os.Exit(1)
	}
	
	fmt.Println("\nSystem ready for normal operation.")
}