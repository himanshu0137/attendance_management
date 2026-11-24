# Attendance Management System Specification

## Overview

The Attendance Management System is a comprehensive solution designed to handle employee attendance tracking, organization management, payroll processing, and leave management. This system provides a modular approach to manage all aspects of workforce management within an organization.

## Core Modules

### 1. Organization Management
This module handles employee management and organizational policies.

#### Features:
- Employee creation and management
- Organization policy settings (leave limits, overtime rules)
- Salary structure definition and management
- Employee onboarding workflow
- HR letter generation (offer letters, appointment letters, etc.)

### 2. Attendance Tracking
This module manages real-time attendance records and reporting.

#### Features:
- Punch-in/punch-out functionality
- Location-based tracking
- Shift schedule management
- Daily attendance reports
- Overtime tracking

### 3. Payroll Processing
Handles salary calculations and payroll management.

#### Features:
- Salary record creation and management
- Payroll summary generation
- Tax bracket configuration
- Net salary calculation
- Payment date tracking

### 4. Leave Management
Manages employee leave requests and approvals.

#### Features:
- Leave type configuration
- Leave request submission
- Approval workflow
- Leave balance tracking
- Leave report generation

### 5. User Management
Handles authentication and user access control.

#### Features:
- User registration and login
- Role-based access control
- Token management
- Session handling

## Technical Architecture

The system follows a modular monolith architecture with the following components:

### Backend:
- Built with Go 1.21
- RESTful API using Gin framework
- PostgreSQL database with proper migrations
- Modular design with clearly separated concerns
- Swagger documentation auto-generated
- Structured into models.go and dto.go files for better organization

### Frontend:
- Flutter-based mobile application for employee interaction
- Web interface for administrators
- Responsive design for multiple device types

## Key Capabilities

### Attendance Management
- Real-time attendance tracking with location data
- Shift schedule management
- Automatic calculation of working hours
- Overtime tracking and reporting

### Organization Management
- Complete employee lifecycle management (onboarding, salary structure)
- HR document generation
- Policy configuration for organizational rules

### Payroll Processing
- Automated salary calculations
- Tax bracket configuration
- Payment tracking
- Comprehensive payroll reports

### Leave Management
- Flexible leave type definitions
- Approval workflow with notifications
- Leave balance management
- Historical leave tracking

### User Management
- Secure authentication and authorization
- Role-based access control
- Session management

## Integration Points

### Keycloak Authentication
- OAuth 2.0 based authentication
- Tenant isolation for multi-tenant deployments
- Role-based access control

### Database
- PostgreSQL with proper indexing and relationships
- Migration support for database schema changes

### External Services
- HR document generation services
- Notification systems for leave approvals and payroll processing

## Security Features

- JWT-based authentication
- Role-based access controls
- Tenant isolation
- Secure API endpoints
- Data encryption where necessary

## Performance Considerations

- Optimized database queries with proper indexing
- Caching strategies for frequently accessed data
- Asynchronous processing for heavy operations
- Scalable architecture supporting multiple tenants

## Code Structure

The backend code is organized into modules, each containing:
- `models.go` - Database entity definitions
- `dto.go` - Data Transfer Objects/Persistence Objects
- `service.go` - Business logic implementation
- `repository.go` - Data access layer
- `module.go` - Module registration and routing

Each module follows a consistent pattern for better maintainability and scalability.

## Module Organization

### Organization Module (`modules/organization`)
- `models.go` - OrgPolicy, Employee, SalaryStructure structs
- `dto.go` - CreateEmployeeRequest, OnboardEmployeeReq, Document, EmployeeWithDocuments, etc.
- `service.go` - Business logic for employee management
- `repository.go` - Database operations for organization data
- `module.go` - Gin handlers and module registration

### Attendance Module (`modules/attendance`)
- `models.go` - Punch, Shift, Report structs
- `dto.go` - PunchRequest, PunchResponse, ReportRequest, ReportResponse structs
- `service.go` - Business logic for attendance tracking
- `repository.go` - Database operations for attendance data
- `module.go` - Gin handlers and module registration

### Payroll Module (`modules/payroll`)
- `models.go` - SalaryRecord, PayrollSummary, TaxBracket structs
- `dto.go` - SalaryRecordRequest, PayrollSummaryRequest, PayrollReportRequest, PayrollReportResponse structs
- `service.go` - Business logic for payroll processing
- `repository.go` - Database operations for payroll data
- `module.go` - Gin handlers and module registration

### Leave Module (`modules/leave`)
- `models.go` - LeaveRequest, LeaveType, LeaveBalance structs
- `dto.go` - LeaveRequestDTO, LeaveResponseDTO, LeaveSummaryDTO structs
- `service.go` - Business logic for leave management
- `repository.go` - Database operations for leave data
- `module.go` - Gin handlers and module registration

### User Module (`modules/user`)
- `models.go` - User, UserRole, UserPreference structs
- `dto.go` - UserDTO, UserUpdateRequest, PreferenceDTO structs
- `service.go` - Business logic for user management
- `repository.go` - Database operations for user data
- `module.go` - Gin handlers and module registration