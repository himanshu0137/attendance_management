#!/bin/bash

# Setup script for Attendance Management System with Nginx proxy
# This script builds and starts all services with proper Nginx configuration

echo "Setting up Attendance Management System with Nginx proxy..."
echo "=========================================================="

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "Error: docker-compose could not be found"
    echo "Please install Docker Compose and try again."
    exit 1
fi

# Change to backend directory
cd "$(dirname "$0")"

echo "Building and starting all services..."
docker-compose up -d

echo ""
echo "Services Status:"
echo "================"
echo "PostgreSQL: http://localhost:5432"
echo "Keycloak Admin Console: http://localhost:9080"
echo "Keycloak Admin Credentials: admin / admin123"
echo "Backend API: http://localhost/api/"
echo "Web Application: http://localhost/"
echo ""
echo "To stop the system, run: docker-compose down"