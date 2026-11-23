package config

import (
	"github.com/caarlos0/env/v11"
)

// Config represents the application configuration.
type Config struct {
	HTTPPort           int    `env:"ATTENDANCE_PORT" envDefault:"7080"`
	DatabaseURL        string `env:"ATTENDANCE_DB" envRequired:"true"`
	KeycloakIssuerURL  string `env:"KEYCLOAK_ISSUER" envRequired:"true"`
	KeycloakClientID   string `env:"KEYCLOAK_CLIENT_ID" envRequired:"true"`
	KeycloakClientSecret string `env:"KEYCLOAK_CLIENT_SECRET" envRequired:"true"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"info"`
	EnabledModules     []string `env:"ENABLED_MODULES" envSeparator:","`
}

// Load loads configuration from environment variables.
func Load() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		panic("Failed to parse config: " + err.Error())
	}
	return cfg
}
