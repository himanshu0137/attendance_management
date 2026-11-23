# Implementation Notes for Leave Management Module

## Test Case Coverage

### TC-13.1: Insufficient Leave Balance Rejection
This functionality requires integration with the leave balance management system. 
Current implementation:
- Basic leave request creation is implemented
- Leave balance validation needs to be added to service layer

Future enhancements needed:
- Add leave balance tracking per employee
- Implement balance validation logic in CreateLeaveRequest
- Return appropriate 400 Bad Request error when balance insufficient

### TC-13.2: Approved Leave & On-Duty Attendance Lock
This functionality requires integration with the attendance system.
Current implementation:
- Basic approval workflow is implemented
- No integration with attendance lock mechanisms

Future enhancements needed:
- Add integration with attendance system for leave locking
- Implement logic to prevent attendance punching during approved leave periods
- Add proper audit trail for leave/attendance interactions

## Current Status

The module currently provides:
✓ Core leave request functionality
✓ Approval/rejection workflow  
✓ Role-based access control
✓ Employee self-service capabilities
✓ Proper database schema with constraints
✓ API endpoints for all required operations

## Next Steps

1. Implement leave balance tracking system (for TC-13.1)
2. Integrate with attendance system for locking (for TC-13.2)
3. Add comprehensive error handling and validation
4. Implement additional test cases for edge scenarios