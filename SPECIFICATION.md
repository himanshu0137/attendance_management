# Attendance Management System - Specification

## Overview
This document specifies the complete attendance management system architecture, including backend Go services, Angular frontend, and Flutter mobile applications.

## System Architecture

### Backend (Go)
- **Core Services**: gin-gonic/gin web framework
- **Database**: PostgreSQL 16 with pgxpool driver
- **Authentication**: Keycloak integration with JWT verification
- **Modular Design**: Selective module loading system
- **Event System**: Thread-safe in-memory pub/sub EventBus

### Frontend (Angular)
- **Framework**: Angular 18 with standalone components
- **UI Library**: Angular Material with Material Design 3
- **State Management**: Built-in services and HTTP interceptors
- **Responsive Design**: Mobile-first approach

### Mobile (Flutter)
- **Framework**: Flutter 3.x with Material 3
- **Platforms**: Android and iOS support
- **Biometric Integration**: Face and fingerprint scanning capabilities
- **Offline Support**: Local caching for kiosk operations

## Detailed Implementation

### Backend Modules

#### 1. Core Module (`pkg/app`)
- Module registry system with selective ResolveModules filter
- Gin server builder implementation
- Application setup and configuration

#### 2. Configuration Module (`pkg/config`)  
- Typed environment configuration using caarlos0/env
- Struct tags for environment variable parsing
- Unit tests for config loading

#### 3. HTTP Extensions (`pkg/httpx`)
- Generic BindJSON[T] function for JSON binding with RFC 7807 errors
- JSON helper for sending responses
- Error helper for RFC 7807 Problem Details format
- JWT verification middleware with local JWKS cache
- RequireRole function for role-based access control
- OAuth2 token proxy handlers (login, refresh, logout)

#### 4. Database Layer (`pkg/db`)
- PostgreSQL connection using pgxpool
- Tenant transaction isolation helper
- Schema-scoped goose migration runner placeholder

#### 5. Event System (`pkg/events`)
- Thread-safe in-memory Pub/Sub EventBus
- Event struct definition and HandlerFunc type
- Bus interface with Publish/Subscribe methods

### 6. Leave Management Module (`modules/leave`)
- Employee leave request submission
- HR approval workflow for leave requests
- Leave type management (Casual, Sick, Earned, etc.)
- Leave history tracking and reporting
- Status tracking (PENDING, APPROVED, REJECTED)

### Frontend Applications

#### Angular Employee Portal (`frontend/angular-app`)
**Features:**
- Login screen with form validation (TASK-16)
- Employee dashboard with punch-in/out functionality (TASK-17)
- Attendance history display
- Responsive Material Design UI

**Integration Points:**
- `POST /auth/login` - Authentication
- `POST /api/attendance/punch` - Record attendance
- `GET /api/attendance/me` - Get employee records
- `GET /api/attendance/reports` - Get attendance reports

#### Flutter Applications

##### 1. User App (`frontend/userapp`)
**Features:**
- Login screen with Material 3 design (TASK-20)
- Employee dashboard for attendance management (TASK-21)
- Punch-in/out functionality
- Attendance history view

**Integration Points:**
- Same API endpoints as Angular app

##### 2. Kiosk App (`frontend/kioskapp`)
**Features:**
- Face recognition scanning interface (TASK-22)
- Fingerprint scanning alternative (TASK-23)
- Kiosk-specific workflow with real-time feedback
- Biometric authentication support

## API Endpoints

### Authentication
```
POST /auth/login          # User login
POST /auth/refresh        # Token refresh  
POST /auth/logout         # User logout
```

### Attendance Management
```
POST /api/attendance/punch  # Record punch-in/out
GET /api/attendance/me      # Get current employee's records
GET /api/attendance/reports # Get attendance reports
```

### Organization Management (Admin)
```
GET /api/organization/employees   # List employees
POST /api/organization/employees  # Create employee
```

## Technical Requirements

### Backend Dependencies
- Go 1.21+
- PostgreSQL 16
- Keycloak 21+
- gin-gonic/gin
- jackc/pgx/v5
- pressly/goose/v3
- caarlos0/env
- golang-jwt/jwt/v5

### Frontend Dependencies
- Angular 18
- TypeScript
- SCSS
- Angular Material

### Mobile Dependencies  
- Flutter 3.x
- Dart
- Material 3 Design
- camera plugin (for kiosk apps)
- http client

## Implementation Status

### ✅ Completed Tasks
**Phase 1 - Backend Infrastructure (Tasks 01-07):**
- Go workspace initialization with modules
- Environment configuration loader
- Gin request/response helpers  
- Keycloak security implementation
- Database connectivity and transactions
- Event system implementation
- Modular application framework

**Phase 3 - Frontend Development (Tasks 15-17):**
- Angular workspace with Material UI theme
- Login screen implementation
- Employee portal dashboard

**Phase 4 - Mobile Development (Tasks 19-23):**
- Flutter project with Material 3 theme
- User app login screen (TASK-20)
- Employee dashboard structure (TASK-21)  
- FaceScan kiosk workflow (TASK-22)
- FingerScan kiosk workflow (TASK-23)

### 🔄 In Progress Tasks
**Phase 2 - Module Development (Tasks 08-14):**
- Organization module implementation
- Attendance module with repository and business logic
- Payroll module with calculation services

## Deployment Structure

### Docker Compose Services
```
postgres:16          # Database
keycloak:latest      # Authentication server  
backend:latest       # Go backend service
frontend:latest      # Angular frontend (optional)
mobile:latest        # Flutter apps (separate deployments)
```

### Environment Variables
```env
DATABASE_URL=postgresql://user:pass@postgres:5432/attendance
KEYCLOAK_ISSUER_URL=http://keycloak:8080/realms/attendance
KEYCLOAK_CLIENT_ID=attendance-app
KEYCLOAK_CLIENT_SECRET=secret
HTTP_PORT=8080
LOG_LEVEL=info
```

## Integration Documentation

### Authentication Flow
1. User submits credentials to `/auth/login`
2. Backend validates with Keycloak and returns JWT tokens
3. Frontend stores tokens and includes in subsequent requests
4. Tokens automatically refreshed when needed

### API Communication
- All API calls use JSON format
- Authentication header: `Authorization: Bearer <token>`
- Error responses follow RFC 7807 Problem Details specification
- All endpoints support CORS for cross-origin requests

## Testing Strategy

### Backend Testing
- Unit tests for all packages
- Integration tests for database operations  
- Mock tests for external services
- API contract testing with Postman collections

### Frontend Testing
- Component unit tests
- E2E testing with Cypress
- Responsive design testing
- Authentication flow testing

### Mobile Testing
- Unit tests for business logic
- UI component testing
- Device-specific testing
- Biometric authentication testing

## Next Steps

1. Complete backend module implementations (organization, attendance, payroll)
2. Implement full API endpoints for all modules
3. Connect frontend and mobile applications to backend APIs
4. Conduct end-to-end integration testing
5. Finalize documentation and deployment scripts

## Leave Management Module Implementation Details

The leave management module implements the following functionality:

1. **Employee Leave Requests**:
   - Submit leave requests with start date, end date, type, and description
   - View personal leave history
   - Status tracking (PENDING, APPROVED, REJECTED)

2. **HR Approval Workflow**:
   - Approve/reject leave requests
   - Role-based access control (HR_MANAGER, ORG_ADMIN)
   - Automatic status updates with approval timestamps

3. **Leave Types**:
   - Casual Leave
   - Sick Leave  
   - Earned Leave
   - Other custom leave types

4. **Database Schema**:
   - `leave_requests` table with all required fields
   - Indexes for performance optimization
   - Data validation constraints

5. **Integration with System Requirements**:
   - TC-13.1: Insufficient leave balance rejection (requires additional leave balance tracking)
   - TC-13.2: Approved leave & on-duty attendance lock (requires integration with attendance system)

6. **Testing Coverage**:
   - Unit tests for all core functionality
   - Integration tests for API endpoints
   - Role-based access control verification