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

4. **Stop services**:
```bash
docker compose down
```

### Backend Setup

1. **Build the backend server**:
```bash
cd backend
go build -o server cmd/server/main.go
```

2. **Set environment variables** (using .env file):
```bash
cd backend
export $(grep -v '^#' .env | xargs)
```

3. **Run the backend server**:
```bash
./server
```

### Frontend Setup

1. **Angular Web App**:
```bash
cd frontend/angular-app
npm install
ng serve
```

2. **User Mobile App**:
```bash
cd frontend/userapp
flutter pub get
flutter run
```

3. **Kiosk Mobile App**:
```bash
cd frontend/kioskapp
flutter pub get
flutter run
```

## API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/refresh` - Token refresh  
- `POST /auth/logout` - User logout

### Attendance Management
- `POST /api/attendance/punch` - Record punch-in/out
- `GET /api/attendance/me` - Get employee records
- `GET /api/attendance/reports` - Get reports

### Organization Management (Admin)
- `GET /api/organization/employees` - List employees
- `POST /api/organization/employees` - Create employee

## Testing

Run backend unit tests:
```bash
cd backend
go test ./pkg/...
```

## Configuration

Environment variables defined in `.env` file:
- `ATTENDANCE_PORT=7080`
- `ATTENDANCE_DB="postgresql://attendance:attendance@localhost:5432/attendance"`
- `KEYCLOAK_ISSUER="http://localhost:9080/realms/attendance"`
- `KEYCLOAK_CLIENT_ID="attendance-app"`
- `KEYCLOAK_CLIENT_SECRET="secret"`
- `LOG_LEVEL=info`
- `ENABLED_MODULES=` (empty = all modules)

## Folder Structure

```
.
├── backend/               # Go backend services
│   ├── cmd/               # Command line tools
│   ├── pkg/               # Package libraries
│   ├── modules/           # Feature modules
│   ├── docker-compose.yml # Docker configuration
│   └── .env               # Environment variables
├── frontend/              # Frontend applications  
│   ├── angular-app/       # Angular web app
│   ├── userapp/           # Employee mobile app
│   └── kioskapp/          # Kiosk mobile app
└── SPECIFICATION.md       # Complete system specification
```

## Security

- JWT-based authentication with Keycloak integration
- Role-based access control
- HTTPS support (recommended for production)
- Secure token handling

## Deployment

The system is designed to be containerized and can be deployed using Docker Compose or Kubernetes.