package main

import (
	"fmt"

	"hunrx.com/attendance_management/pkg/config"
)

func main() {
	fmt.Println("Testing configuration loading...")
	
	// Load configuration - this will panic if required env vars are missing
	cfg := config.Load()
	
	fmt.Printf("Configuration loaded successfully!\n")
	fmt.Printf("HTTP Port: %d\n", cfg.HTTPPort)
	fmt.Printf("Database URL: %s\n", cfg.DatabaseURL)
	fmt.Printf("Keycloak Issuer: %s\n", cfg.KeycloakIssuerURL)
	fmt.Printf("Keycloak Client ID: %s\n", cfg.KeycloakClientID)
	fmt.Printf("Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("Enabled Modules: %v\n", cfg.EnabledModules)
	
	fmt.Println("\nConfiguration test PASSED!")
}