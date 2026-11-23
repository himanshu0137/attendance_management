package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"hunrx.com/attendance_management/pkg/app"
	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/db"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		slog.Error("Database URL is required")
		os.Exit(1)
	}
	if cfg.KeycloakIssuerURL == "" {
		slog.Error("Keycloak Issuer URL is required")
		os.Exit(1)
	}
	if cfg.KeycloakClientID == "" {
		slog.Error("Keycloak Client ID is required")
		os.Exit(1)
	}
	if cfg.KeycloakClientSecret == "" {
		slog.Error("Keycloak Client Secret is required")
		os.Exit(1)
	}

	// Set up logging
	var logLevel slog.Level
	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	
	slog.SetLogLoggerLevel(logLevel)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		slog.Info("Shutting down server...")
		cancel()
	}()

	// Connect to database
	slog.Info("Connecting to database...", "url", cfg.DatabaseURL)
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Create core context
	core := &app.CoreContext{}

	// Setup the application
	slog.Info("Setting up application...")
	if err := app.Setup(ctx, cfg, pool, core); err != nil {
		slog.Error("Failed to setup application", "error", err)
		os.Exit(1)
	}

	// Create and run server
	server, err := app.New(core)
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Run the server
	slog.Info("Starting server...", "port", cfg.HTTPPort)
	if err := server.Run(ctx, "0.0.0.0", cfg.HTTPPort); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}