# Leave Management Module

This module implements the leave management functionality for the Attendance Management System.

## Features

- Employee leave request submission
- HR approval workflow for leave requests  
- Leave type management (Casual, Sick, Earned, etc.)
- Leave history tracking and reporting
- Status tracking (PENDING, APPROVED, REJECTED)

## API Endpoints

### Submit Leave Request
```
POST /api/leaves
```
Submit a new leave request with:
- start_date: Date when leave starts
- end_date: Date when leave ends  
- leave_type: Type of leave (Casual, Sick, Earned, etc.)
- description: Optional description

### View Leave Requests
```
GET /api/leaves
```
List all leave requests (HR/Admin only)

```
GET /api/leaves/:id
```
Get specific leave request details

### Approve/Reject Leave Request
```
POST /api/leaves/:id/approve
```
Approve a leave request (HR/Admin only)

```
POST /api/leaves/:id/reject  
```
Reject a leave request (HR/Admin only)

### Employee Self-Service
```
GET /api/leaves/my
```
List employee's own leave requests

## Database Schema

The module uses the `leave_requests` table with the following fields:
- id: Unique identifier
- employee_id: Reference to employee
- tenant_id: Organization identifier  
- start_date: Leave start date
- end_date: Leave end date
- leave_type: Type of leave requested
- description: Optional description
- status: Current status (PENDING/APPROVED/REJECTED)
- approved_by: Approving user ID (if applicable)
- approved_at: Approval timestamp (if applicable)
- created_at: Record creation timestamp
- updated_at: Last update timestamp

## Authorization

- All endpoints require tenant context
- Leave submission: Any authenticated user (EMPLOYEE role)
- View all requests: HR_MANAGER or ORG_ADMIN roles  
- Approve/reject: HR_MANAGER or ORG_ADMIN roles
- Employee self-service: EMPLOYEE role only

## Testing

The module includes comprehensive unit tests covering:
- Create leave request functionality
- Approve/reject leave requests
- List employee's leave history
- Status tracking and updates

Tests are located in `service_test.go` and verify proper behavior for all core functionalities.