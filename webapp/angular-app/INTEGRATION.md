# Frontend-Backend Integration Documentation

This document describes how the Angular frontend integrates with the Go backend for the Attendance Management System.

## Authentication Flow

### 1. Login Endpoint
The login form sends credentials to the backend's authentication endpoint:

```
POST /auth/login
Content-Type: application/json

{
  "username": "employee@example.com",
  "password": "securePassword"
}
```

Backend response:
```
200 OK
Content-Type: application/json

{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

### 2. Protected Routes
All protected routes require a valid JWT token in the Authorization header:

```
GET /api/attendance/reports
Authorization: Bearer <access_token>
```

## API Endpoints

### Attendance Management
- `POST /api/attendance/punch` - Record punch-in/out
- `GET /api/attendance/me` - Get current employee's attendance records  
- `GET /api/attendance/reports` - Get attendance reports

### Organization Management (Admin)
- `GET /api/organization/employees` - List employees
- `POST /api/organization/employees` - Create employee

## Token Handling

The frontend should:
1. Store tokens in memory or secure HTTP-only cookies
2. Automatically refresh tokens when they expire using the refresh endpoint
3. Redirect to login when tokens are invalid or expired

## Environment Configuration

Frontend should read API base URL from environment variables:

```typescript
// environments/environment.ts
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api'
};
```

## Integration Points

### 1. Authentication Service
Handles login/logout, token storage, and refresh logic.

### 2. HTTP Interceptor
Automatically adds Authorization headers to API requests.

### 3. Guard Services
Protect routes that require authentication or specific roles.