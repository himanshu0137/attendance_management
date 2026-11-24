package attendance

// Punch represents an attendance punch record.
type Punch struct {
	ID           string `json:"id" db:"id"`
	EmployeeID   string `json:"employee_id" db:"employee_id"`
	TenantID     string `json:"tenant_id" db:"tenant_id"`
	PunchTime    int64  `json:"punch_time" db:"punch_time"`
	PunchType    string `json:"punch_type" db:"punch_type"` // "IN" or "OUT"
	Location     string `json:"location" db:"location"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
}

// Shift represents an employee's shift schedule.
type Shift struct {
	ID            string `json:"id" db:"id"`
	EmployeeID    string `json:"employee_id" db:"employee_id"`
	TenantID      string `json:"tenant_id" db:"tenant_id"`
	ShiftName     string `json:"shift_name" db:"shift_name"`
	StartTime     int64  `json:"start_time" db:"start_time"`
	EndTime       int64  `json:"end_time" db:"end_time"`
	WorkingDays   []int  `json:"working_days" db:"working_days"` // Array of day numbers (0=Sunday, 1=Monday, etc.)
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
}

// Report represents an attendance report.
type Report struct {
	ID           string `json:"id" db:"id"`
	EmployeeID   string `json:"employee_id" db:"employee_id"`
	TenantID     string `json:"tenant_id" db:"tenant_id"`
	ReportDate   int64  `json:"report_date" db:"report_date"`
	TotalHours   float64 `json:"total_hours" db:"total_hours"`
	Status       string `json:"status" db:"status"` // "PRESENT", "ABSENT", "LEAVE"
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
}