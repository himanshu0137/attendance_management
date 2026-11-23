# Testing Plan for Attendance Management System

## Overview
This document outlines the testing approach for the complete attendance management system with backend, frontend, and mobile components.

## Prerequisites

### Required Software:
1. Docker (for running PostgreSQL and Keycloak)
2. Go 1.21+ (for backend server)
3. Node.js & npm (for Angular frontend) 
4. Flutter SDK (for mobile apps)
5. PostgreSQL client tools (optional for database inspection)

## System Architecture Test Setup

### Backend Services
- **PostgreSQL 16**: Database for storing attendance data
- **Keycloak**: Authentication server with JWT support
- **Go Backend Server**: Main application serving APIs

### Frontend Applications  
- **Angular Web App**: Employee portal and admin dashboard
- **Flutter Mobile Apps**: 
  - User app: For employee attendance management
  - Kiosk app: For face/fingerprint scanning kiosks

## Testing Approach

### Phase 1: Backend Service Testing

#### 1. Database Connection Test
```bash
# Start database service (requires Docker)
docker run -d --name test-postgres \
  -e POSTGRES_USER=attendance \
  -e POSTGRES_PASSWORD=attendance \
  -e POSTGRES_DB=attendance \
  -p 5432:5432 \
  postgres:16

# Test connection
psql -h localhost -p 5432 -U attendance -d attendance
```

#### 2. Keycloak Service Test
```bash
# Start Keycloak service (requires Docker)
docker run -d --name test-keycloak \
  -e KC_BOOTSTRAP_ADMIN_USERS=admin \
  -e KC_HOSTNAME_STRICT=false \
  -e KC_HOSTNAME=localhost \
  -p 8080:8080 \
  quay.io/keycloak/keycloak:latest start-dev
```

#### 3. Backend Server Test
```bash
# Set environment variables
export ATTENDANCE_PORT=8080
export ATTENDANCE_DB="postgresql://attendance:attendance@localhost:5432/attendance"
export KEYCLOAK_ISSUER="http://localhost:8080/realms/attendance"
export KEYCLOAK_CLIENT_ID="attendance-app" 
export KEYCLOAK_CLIENT_SECRET="secret"
export LOG_LEVEL=debug

# Build and run backend
go build -o server cmd/server/main.go
./server
```

### Phase 2: API Testing

#### Authentication Endpoints
```bash
# Test login endpoint (after backend is running)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"employee@example.com","password":"password"}'
```

#### Attendance Endpoints  
```bash
# Test punch endpoint (after authentication)
curl -X POST http://localhost:8080/api/attendance/punch \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"type":"in"}'
```

### Phase 3: Frontend Testing

#### Angular Application
```bash
# Install dependencies
cd frontend/angular-app
npm install

# Run development server  
ng serve
```

#### Mobile Applications
```bash
# User app
cd frontend/userapp
flutter pub get
flutter run

# Kiosk app  
cd frontend/kioskapp
flutter pub get
flutter run
```

## Integration Test Scenarios

### 1. Authentication Flow
- User logs in via Angular web app
- Backend validates credentials with Keycloak
- JWT token is returned to client
- Subsequent requests include Authorization header

### 2. Attendance Recording
- Employee punches in via mobile app or web portal
- System records timestamp in database
- Response includes confirmation details

### 3. Report Generation
- Admin user requests attendance reports
- Backend queries database and returns structured data

## Test Automation

### Unit Tests
```bash
cd backend
go test ./pkg/...
```

### Integration Tests
Use Postman collections or similar tools to test API endpoints.

### End-to-End Tests
- Cypress for Angular frontend testing
- Flutter driver for mobile app testing

## Expected Results

### Backend Service Status:
- Server starts successfully on port 8080
- Database connection established
- Keycloak integration working
- All modules loaded correctly

### API Endpoints:
- `POST /auth/login` returns valid JWT tokens
- `POST /api/attendance/punch` records attendance
- `GET /api/attendance/reports` returns data
- Proper error responses for invalid requests

### Frontend Functionality:
- Login screen accepts credentials
- Employee dashboard displays correctly
- Punch-in/out buttons work
- Mobile apps show proper interfaces

## Troubleshooting

### Common Issues:
1. **Database Connection Refused**: Ensure PostgreSQL is running with correct credentials
2. **Keycloak Not Available**: Verify Keycloak service is running at configured URL  
3. **Port Conflicts**: Check that ports 8080, 5432, 8080 are available
4. **Environment Variables**: Verify all required environment variables are set

### Debugging Commands:
```bash
# Check if services are running
docker ps

# View logs
docker logs <container_name>

# Test database connectivity
psql -h localhost -p 5432 -U attendance -d attendance

# Test API endpoints  
curl http://localhost:8080/health
```

## Documentation

All system components are documented in:
- SPECIFICATION.md (complete system spec)
- COMPLETED_TASKS_SUMMARY.md (task completion tracking)
- INTEGRATION.md (API integration details)