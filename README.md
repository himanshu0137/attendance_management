# Attendance Management System

## Overview
This is a complete attendance management system with:
- Go backend services with modular architecture
- Angular web frontend 
- Flutter mobile applications for employees and kiosks

## System Architecture

### Backend Components
- **Go Backend**: gin-gonic/gin web framework with modular design
- **Database**: PostgreSQL 16 with pgxpool driver
- **Authentication**: Keycloak integration with JWT verification
- **Event System**: Thread-safe in-memory pub/sub EventBus

### Frontend Components  
- **Angular Web App**: Employee portal with Material Design UI
- **Flutter Mobile Apps**:
  - User app: For employee attendance management
  - Kiosk app: For face/fingerprint scanning kiosks

## Getting Started

### Prerequisites
1. Docker (for running services)
2. Go 1.21+ (for backend compilation)  
3. Node.js & npm (for Angular frontend)
4. Flutter SDK (for mobile apps)

### Running Services with Docker

1. **Start PostgreSQL and Keycloak**:
```bash
cd backend
docker compose up -d
```

2. **Verify services are running**:
```bash
docker ps
```

3. **Check service status**:
```bash
# Check PostgreSQL
docker logs <postgres_container_id>

# Check Keycloak  
docker logs <keycloak_container_id>
```

### Backend Setup

1. **Build the backend server**:
```bash
cd backend
go build -o server cmd/server/main.go
```

2. **Set environment variables**:
```bash
export ATTENDANCE_PORT=7080
export ATTENDANCE_DB="postgresql://attendance:attendance@localhost:5432/attendance"
export KEYCLOAK_ISSUER="http://localhost:9080/realms/attendance" 
export KEYCLOAK_CLIENT_ID="attendance-app"
export KEYCLOAK_CLIENT_SECRET="secret"
export LOG_LEVEL=debug
```

3. **Run the backend server**:
```bash
./server
```

### Frontend Setup

1. **Angular Web App**:
```bash
cd webapp/angular-app
npm install
ng serve
```

2. **User Mobile App**:
```bash
cd userapp
flutter pub get
flutter run
```

3. **Kiosk Mobile App**:
```bash
cd kioskapp
flutter pub get
flutter run
```

## API Documentation

The system provides comprehensive API documentation through Swagger/OpenAPI:

- **Swagger UI**: http://localhost:7080/swagger/index.html
- **OpenAPI JSON**: http://localhost:7080/swagger/doc.json

### Available Endpoints

The API includes endpoints across all modules:

#### Authentication
- `POST /api/auth/login` - User login with Keycloak
- `POST /api/auth/refresh` - Refresh authentication token  
- `POST /api/auth/logout` - Logout user

#### Attendance Management
- `POST /api/attendance/punch` - Record punch-in/out
- `GET /api/attendance/me` - Get employee records
- `GET /api/attendance/reports` - Get attendance reports

#### Organization Management (Admin)
- `GET /api/organization/{tenant_id}/policies` - Get organization policy
- `POST /api/organization/{tenant_id}/policies` - Set organization policy
- `GET /api/organization/{tenant_id}/employees` - List employees
- `POST /api/organization/{tenant_id}/employees` - Create employee

#### Payroll Management
- `POST /api/payroll/slips/generate` - Generate payroll slip
- `GET /api/payroll/slips/me` - Get employee payroll slips
- `GET /api/payroll/slips` - Get all payroll slips

#### User Management
- `GET /api/user/me` - Get current user details
- `PUT /api/user/me` - Update current user details
- `GET /api/user/me/preferences` - Get user preferences
- `PUT /api/user/me/preferences` - Update user preferences

All endpoints are documented with Swagger comments for full API exploration and testing.

## Testing

Run backend unit tests:
```bash
cd backend
go test ./pkg/...
```

## Configuration

Environment variables:
- `ATTENDANCE_PORT=7080` (updated from 8080)
- `ATTENDANCE_DB="postgresql://attendance:attendance@localhost:5432/attendance"`
- `KEYCLOAK_ISSUER="http://localhost:9080/realms/attendance"` (updated port)
- `KEYCLOAK_CLIENT_ID="attendance-app"`
- `KEYCLOAK_CLIENT_SECRET="secret"`
- `LOG_LEVEL=info`
- `ENABLED_MODULES` = (all modules)

## Folder Structure

```
.
├── backend/               # Go backend services
│   ├── cmd/               # Command line tools
│   ├── pkg/               # Package libraries
│   ├── modules/           # Feature modules
│   └── docker-compose.yml # Docker configuration
├── webapp/                # Angular web frontend  
│   ├── angular-app/       # Angular web app components
├── flutter-app/           # Flutter mobile app (general)
├── kioskapp/              # Kiosk mobile app
└── userapp/               # User mobile app
└── SPECIFICATION.md       # Complete system specification
```

## Security

- JWT-based authentication with Keycloak integration
- Role-based access control
- HTTPS support (recommended for production)
- Secure token handling

## Deployment

The system is designed to be containerized and can be deployed using Docker Compose or Kubernetes.