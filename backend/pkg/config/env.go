package config

import (
	"os"
	"strconv"
	"strings"
)

// ServerConfig is the server-wide configuration.
type ServerConfig struct {
	EnabledModules []string
	Host           string
	Port           int
	DBDSN          string
	Issuer         string
	ClientID       string
	ClientSecret   string
}

// Load parses the environment into ServerConfig.
func Load() ServerConfig {
	cfg := ServerConfig{
		Host:          "0.0.0.0",
		Port:          8080,
		EnabledModules: nil,
	}
	for _, e := range os.Environ() {
		k, v, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		switch k {
		case "ATTENDANCE_HOST":
			cfg.Host = v
		case "ATTENDANCE_PORT":
			p, err := strconv.Atoi(v)
			if err == nil {
				cfg.Port = p
			}
		case "ATTENDANCE_DB":
			cfg.DBDSN = v
		case "ENABLED_MODULES":
			cfg.EnabledModules = strings.Split(v, ",")
		case "KEYCLOAK_ISSUER":
			cfg.Issuer = v
		case "KEYCLOAK_CLIENT_ID":
			cfg.ClientID = v
		case "KEYCLOAK_CLIENT_SECRET":
			cfg.ClientSecret = v
		}
	}
	return cfg
}
