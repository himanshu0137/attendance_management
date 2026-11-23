#!/bin/bash

# Setup script to load environment variables from .env file
# Usage: source setup_env.sh

echo "Loading environment variables from .env file..."

if [ -f ".env" ]; then
    # Load the environment variables
    export $(grep -v '^#' .env | xargs)
    echo "Environment variables loaded successfully!"
    echo "ATTENDANCE_PORT=$ATTENDANCE_PORT"
    echo "KEYCLOAK_ISSUER=$KEYCLOAK_ISSUER"
    echo "LOG_LEVEL=$LOG_LEVEL"
else
    echo "Error: .env file not found!"
    exit 1
fi