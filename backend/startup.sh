#!/bin/bash

# Startup script for Attendance Management System
# This script sets up and starts all necessary services

echo "Starting Attendance Management System..."
echo "========================================"

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "Error: docker-compose could not be found"
    echo "Please install Docker Compose and try again."
    exit 1
fi

# Change to backend directory
cd "$(dirname "$0")"

echo "Starting database and authentication services..."
docker-compose up -d postgres keycloak

echo "Waiting for services to start..."
sleep 10

# Check if services are running
if docker-compose ps | grep -q "Up.*keycloak"; then
    echo "✅ Keycloak is running"
else
    echo "⚠️  Keycloak may not be running properly"
fi

if docker-compose ps | grep -q "Up.*postgres"; then
    echo "✅ PostgreSQL is running"
else
    echo "⚠️  PostgreSQL may not be running properly"
fi

echo ""
echo "System Status:"
echo "=============="
echo "Backend API: http://localhost:7080"
echo "Keycloak Admin Console: http://localhost:9080"
echo "Keycloak Admin Credentials: admin / admin123"
echo ""
echo "To stop the system, run: docker-compose down"