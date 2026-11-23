# Attendance Management System - Backend

## Overview

This is the backend server for the Attendance Management System, built with Go using Gin framework and PostgreSQL database.

## Prerequisites

- Go 1.27 or later
- PostgreSQL database
- Docker (for running services)

## Running the Server

### Method 1: Auto-run script (Recommended)
```bash
./auto_run.sh
```

This automatically loads `.env` file if present and starts the server.

### Method 2: Manual environment setup
```bash
# Load environment variables
export $(grep -v '^#' .env | xargs)

# Run the server
./server
```

### Method 3: Direct execution with environment
```bash
ATTENDANCE_PORT=7080 \
ATTENDANCE_DB="postgresql://attendance:attendance@localhost:5432/attendance" \
KEYCLOAK_ISSUER="http://localhost:9080/realms/attendance" \
KEYCLOAK_CLIENT_ID="attendance-app" \
KEYCLOAK_CLIENT_SECRET="secret" \
LOG_LEVEL="info" \
./server
```

## API Documentation

The API is documented using Swagger/OpenAPI:

- **Swagger UI**: http://localhost:7080/swagger/index.html
- **OpenAPI JSON**: http://localhost:7080/swagger/doc.json

## Configuration

The server uses environment variables from `.env` file. See `.env.example` for all available options.

## API Endpoints

- `POST /api/auth/login` - Authenticate user
- `POST /api/auth/refresh` - Refresh authentication token
- `POST /api/auth/logout` - Logout user

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| ATTENDANCE_PORT | Server port | 7080 |
| ATTENDANCE_DB | PostgreSQL connection string | (required) |
| KEYCLOAK_ISSUER | Keycloak issuer URL | (required) |
| KEYCLOAK_CLIENT_ID | Keycloak client ID | (required) |
| KEYCLOAK_CLIENT_SECRET | Keycloak client secret | (required) |
| LOG_LEVEL | Log level (debug, info, warn, error) | info |

## Development

### Build
```bash
go build -o server cmd/server/main.go
```

### Test
```bash
go test ./...
```

## Troubleshooting

If you encounter database connection issues:
1. Ensure PostgreSQL is running
2. Check that the database URL in `.env` is correct
3. Verify database credentials and permissions
4. Make sure the database exists

## License

MIT