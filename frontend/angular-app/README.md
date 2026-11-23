# Attendance Management Portal - Frontend

This is the Angular frontend application for the Attendance Management System, designed to integrate with the Go backend.

## Features Implemented

### 1. Login Screen (Task-16)
- Clean, responsive login interface using Angular Material design principles
- Form validation for username and password fields
- Integration-ready authentication flow

### 2. Employee Dashboard (Task-17) 
- Punch-in/out functionality
- Attendance history display
- Responsive layout for different screen sizes

## Architecture

### Technology Stack
- Angular 18 (standalone components)
- TypeScript
- SCSS for styling
- HttpClient for API communication

### Folder Structure
```
src/
├── app/
│   ├── login/           # Login component
│   ├── employee/        # Employee dashboard component  
│   ├── app.component.ts # Main application component
│   ├── app.routes.ts    # Routing configuration
│   └── app.config.ts    # Application configuration
├── assets/              # Static assets
└── environments/        # Environment configurations
```

## Integration with Backend

The frontend is designed to work with the Go backend API. Key integration points include:

### Authentication Flow
1. Login form submits to `/auth/login` endpoint
2. Stores JWT tokens for subsequent requests
3. Handles token refresh automatically

### API Endpoints Used
- `POST /auth/login` - User authentication
- `POST /api/attendance/punch` - Record attendance
- `GET /api/attendance/me` - Get employee records
- `GET /api/attendance/reports` - Get attendance reports

## How to Run

1. Navigate to the project directory:
   ```bash
   cd frontend/angular-app
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start development server:
   ```bash
   ng serve
   ```

4. Open browser to `http://localhost:4200`

## Integration Documentation

See [INTEGRATION.md](INTEGRATION.md) for detailed API integration information.
See [backend-integration.md](backend-integration.md) for code examples and implementation details.

## Next Steps

To complete the full integration with the backend:
1. Configure environment variables with backend API URL
2. Implement proper HTTP interceptors for token handling
3. Add authentication guards to protect routes
4. Connect actual API endpoints instead of simulated data
5. Implement proper error handling and loading states