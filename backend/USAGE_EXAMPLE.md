# Usage Example - Attendance Management System

## Complete Setup and Testing Guide

### 1. Environment Setup

First, let's make sure we have the proper environment variables set:

```bash
cd /home/carbon/prog/attendance_management/backend

# Check current environment
echo "Current ATTENDANCE_DB: $ATTENDANCE_DB"
echo "Current KEYCLOAK_ISSUER: $KEYCLOAK_ISSUER"

# If environment is not set, use the .env file
if [ -f ".env" ]; then
    echo "Loading from .env file..."
    export $(grep -v '^#' .env | xargs)
fi

echo "After loading:"
echo "ATTENDANCE_DB=$ATTENDANCE_DB"
echo "KEYCLOAK_ISSUER=$KEYCLOAK_ISSUER"
```

### 2. Test Database Connection

Let's test our database connection with the proper environment:

```bash
cd /home/carbon/prog/attendance_management/backend

# Load environment variables
export $(grep -v '^#' .env | xargs)

# Test database connection
echo "Testing database connection..."
go run test_db_connection.go
```

### 3. Run the Server

If you want to run the server properly, make sure Docker services are running:

```bash
cd /home/carbon/prog/attendance_management/backend

# Start Docker services (if not already running)
docker compose up -d

# Wait a moment for services to start
sleep 5

# Load environment variables
export $(grep -v '^#' .env | xargs)

# Build and run server
go build -o server cmd/server/main.go
echo "Starting server on port $ATTENDANCE_PORT..."
./server
```

### 4. Environment Variable Verification

Here's a complete verification script:

```bash
cd /home/carbon/prog/attendance_management/backend

# Verify all required environment variables are set
echo "=== Environment Variables Check ==="
echo "ATTENDANCE_PORT: $ATTENDANCE_PORT"
echo "ATTENDANCE_DB: $ATTENDANCE_DB" 
echo "KEYCLOAK_ISSUER: $KEYCLOAK_ISSUER"
echo "KEYCLOAK_CLIENT_ID: $KEYCLOAK_CLIENT_ID"
echo "LOG_LEVEL: $LOG_LEVEL"

# Check if all required variables are set
if [ -z "$ATTENDANCE_DB" ]; then
    echo "ERROR: ATTENDANCE_DB is not set!"
    exit 1
fi

if [ -z "$KEYCLOAK_ISSUER" ]; then
    echo "ERROR: KEYCLOAK_ISSUER is not set!"
    exit 1
fi

echo "All required environment variables are set correctly."
```

### 5. Docker Service Verification

Check that Docker services are running properly:

```bash
cd /home/carbon/prog/attendance_management/backend

# Check if Docker containers are running
echo "=== Docker Services Status ==="
docker ps

# If containers are not running, start them
if [ -z "$(docker ps -q)" ]; then
    echo "Starting Docker services..."
    docker compose up -d
    sleep 5
fi

# Check PostgreSQL specifically  
echo "=== PostgreSQL Service ==="
docker logs $(docker ps -q --filter "name=postgres") 2>/dev/null || echo "PostgreSQL container not found"

# Check Keycloak specifically
echo "=== Keycloak Service ===" 
docker logs $(docker ps -q --filter "name=keycloak") 2>/dev/null || echo "Keycloak container not found"
```

### 6. Expected Output

When everything is properly configured, you should see output like:

```
=== Testing Database Connection ===
Loaded configuration:
  HTTP Port: 7080
  Database URL: postgresql://attendance:attendance@localhost:5432/attendance
  Keycloak Issuer: http://localhost:9080/realms/attendance

Attempting to connect to database at: postgresql://attendance:attendance@localhost:5432/attendance
SUCCESS: Database connection established!
The system is ready to run with proper database connectivity.
```

### Common Issues and Solutions

1. **Database Connection Error**:
   ```
   ERROR: Failed to connect to database: failed to connect to `user=carbon database=`: /tmp/.s.PGSQL.5432 (/tmp): dial error: dial unix /tmp/.s.PGSQL.5432: connect: no such file or directory
   ```
   
   **Solution**: Make sure Docker services are running and environment variables are properly set.

2. **Missing Environment Variables**:
   ```
   ERROR: Database URL is empty! Please check your environment variables.
   ```
   
   **Solution**: Ensure all environment variables are loaded from the .env file or set manually.

3. **Docker Services Not Running**:
   ```
   Error response from daemon: pull access denied for postgres, repository does not exist
   ```
   
   **Solution**: Run `docker compose up -d` to start services.

### Testing Your Setup

Run these commands to test your complete setup:

```bash
cd /home/carbon/prog/attendance_management/backend

# 1. Load environment
export $(grep -v '^#' .env | xargs)

# 2. Verify environment variables are set
echo "Environment loaded: ATTENDANCE_DB=$ATTENDANCE_DB"

# 3. Test database connection
go run test_db_connection.go

# 4. If successful, build and start server
go build -o server cmd/server/main.go
echo "Server built successfully!"
```

This setup will allow you to properly test the system once Docker services are running.