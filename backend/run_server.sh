#!/bin/bash

# Set up environment variables from .env file if present
if [ -f ".env" ]; then
    echo "Loading environment from .env file..."
    export $(grep -v '^#' .env | xargs)
fi

# Set defaults if not already set
: ${ATTENDANCE_PORT:=7080}
: ${LOG_LEVEL:=info}

echo "Starting Attendance Management Server on port $ATTENDANCE_PORT"
echo "Database URL: ${ATTENDANCE_DB:-(not set)}"
echo "Keycloak Issuer: ${KEYCLOAK_ISSUER:-(not set)}"

# Run the server
go run main.go