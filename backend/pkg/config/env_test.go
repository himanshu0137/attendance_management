package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Set up test environment
	os.Setenv("ATTENDANCE_PORT", "9090")
	os.Setenv("ATTENDANCE_DB", "postgresql://test:pass@localhost/testdb")
	os.Setenv("KEYCLOAK_ISSUER", "http://keycloak.test/realms/test")
	os.Setenv("KEYCLOAK_CLIENT_ID", "test-client")
	os.Setenv("KEYCLOAK_CLIENT_SECRET", "test-secret")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENABLED_MODULES", "org,attendance")

	cfg := Load()
	
	assert.Equal(t, 9090, cfg.HTTPPort)
	assert.Equal(t, "postgresql://test:pass@localhost/testdb", cfg.DatabaseURL)
	assert.Equal(t, "http://keycloak.test/realms/test", cfg.KeycloakIssuerURL)
	assert.Equal(t, "test-client", cfg.KeycloakClientID)
	assert.Equal(t, "test-secret", cfg.KeycloakClientSecret)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, []string{"org", "attendance"}, cfg.EnabledModules)
	
	// Clean up
	os.Unsetenv("ATTENDANCE_PORT")
	os.Unsetenv("ATTENDANCE_DB")
	os.Unsetenv("KEYCLOAK_ISSUER")
	os.Unsetenv("KEYCLOAK_CLIENT_ID")
	os.Unsetenv("KEYCLOAK_CLIENT_SECRET")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("ENABLED_MODULES")
}