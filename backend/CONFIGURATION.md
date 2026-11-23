# Configuration Guide

## Environment Variables

The Attendance Management System uses environment variables for configuration. All required settings must be provided.

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `ATTENDANCE_PORT` | HTTP port for the server | 7080 |
| `ATTENDANCE_DB` | PostgreSQL connection string | postgresql://user:pass@localhost:5432/attendance |
| `KEYCLOAK_ISSUER` | Keycloak realm URL | http://localhost:9080/realms/attendance |
| `KEYCLOAK_CLIENT_ID` | OAuth2 client ID in Keycloak | attendance-app |
| `KEYCLOAK_CLIENT_SECRET` | OAuth2 client secret in Keycloak | secret |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | info |
| `ENABLED_MODULES` | Comma-separated list of enabled modules | (all modules) |

## Configuration File Example

```bash
# Server settings
export ATTENDANCE_PORT=7080

# Database configuration  
export ATTENDANCE_DB="postgresql://attendance:attendance@localhost:5432/attendance"

# Keycloak authentication
export KEYCLOAK_ISSUER="http://localhost:9080/realms/attendance"
export KEYCLOAK_CLIENT_ID="attendance-app" 
export KEYCLOAK_CLIENT_SECRET="secret"

# Logging
export LOG_LEVEL=info
```

## Keycloak Configuration

The system is configured to work with Keycloak running on port 9080. Make sure your Keycloak realm has:
1. A client with ID "attendance-app"
2. Client secret matching `KEYCLOAK_CLIENT_SECRET`
3. Proper roles defined for access control

## Database Connection

The PostgreSQL database should be accessible at the URL specified in `ATTENDANCE_DB`. The default setup uses:
- User: attendance
- Password: attendance  
- Database: attendance
- Port: 5432 (mapped from host)

## Module Configuration

The `ENABLED_MODULES` variable allows selective activation of modules. If not set, all modules are enabled.
Example: `export ENABLED_MODULES=attendance,organization`