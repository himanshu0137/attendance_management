#!/bin/bash

# Entry point for Keycloak that runs initialization script before starting Keycloak

# Make the init script executable
chmod +x /opt/keycloak/init/init-keycloak.sh

# Run the initialization script
echo "Running Keycloak initialization..."
/opt/keycloak/init/init-keycloak.sh

# Start Keycloak with the original command
echo "Starting Keycloak..."
exec "$@"