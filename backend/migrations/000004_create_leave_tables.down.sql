-- Drop indexes
DROP INDEX IF EXISTS idx_leave_requests_employee_id;
DROP INDEX IF EXISTS idx_leave_requests_tenant_id;
DROP INDEX IF EXISTS idx_leave_requests_status;
DROP INDEX IF EXISTS idx_leave_requests_start_date;
DROP INDEX IF EXISTS idx_leave_requests_end_date;

-- Drop the leave_requests table
DROP TABLE IF EXISTS leave_requests;