# Attendance Management Portal - Mobile (Flutter)

This is the Flutter mobile application for the Attendance Management System, designed to integrate with the Go backend.

## Features Implemented

### 1. Login Screen (Task-20)
- Material 3 design theme
- Clean, responsive login interface  
- Form validation for username and password

### 2. Employee Portal (Task-21)
- Punch-in/out functionality
- Attendance history display
- FaceScan and FingerScan attendance workflows

## Architecture

### Technology Stack
- Flutter 3.x
- Dart
- Material 3 design system
- HTTP client for API communication

### Folder Structure
```
lib/
├── main.dart            # Entry point
├── screens/             # Screen components
│   ├── login_screen.dart
│   ├── employee_screen.dart
│   └── kiosk_screens/
│       ├── face_scan_screen.dart
│       └── fingerprint_scan_screen.dart
├── services/            # API services
│   ├── auth_service.dart
│   └── attendance_service.dart
├── models/              # Data models
│   ├── user.dart
│   └── attendance_record.dart
└── utils/               # Utility functions
```

## Integration with Backend

The mobile app is designed to work with the Go backend API. Key integration points include:

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
   cd frontend/flutter-app
   ```

2. Install dependencies:
   ```bash
   flutter pub get
   ```

3. Start development server:
   ```bash
   flutter run
   ```

## Implementation Status

This is a structural scaffold for the Flutter implementation. The following tasks remain to be completed:

### Task-20: Flutter login screen
- [x] Basic Material 3 login screen structure
- [ ] Implement form validation and authentication logic

### Task-21: Employee screens  
- [x] Basic employee dashboard structure
- [ ] Implement punch-in/out functionality
- [ ] Add attendance history display

### Task-22: Gate Kiosk FaceScan workflow
- [ ] Implement face recognition interface
- [ ] Create face scan flow with camera access

### Task-23: Gate Kiosk FingerScan workflow  
- [ ] Implement fingerprint scanner interface
- [ ] Create fingerprint scan flow with biometric access

## Integration Documentation

The mobile app will follow the same API structure as the Angular frontend:

See [backend-integration.md](../angular-app/backend-integration.md) for detailed API integration information.