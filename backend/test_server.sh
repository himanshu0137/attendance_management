#!/bin/bash

# Read the .env file and export all variables properly
export $(grep -v '^#' .env | xargs)

# Show what we loaded
echo "Loaded environment variables:"
echo "ATTENDANCE_PORT=$ATTENDANCE_PORT"
echo "ATTENDANCE_DB=$ATTENDANCE_DB"
echo "KEYCLOAK_ISSUER=$KEYCLOAK_ISSUER"
echo "KEYCLOAK_CLIENT_ID=$KEYCLOAK_CLIENT_ID"
echo "LOG_LEVEL=$LOG_LEVEL"

# Run the server in background for a short time
echo "Starting server for 10 seconds..."
timeout 10s ./server