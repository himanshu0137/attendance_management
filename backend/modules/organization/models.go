package organization

// OrgPolicy represents the organization policy settings.
type OrgPolicy struct {
	ID         string `json:"id" db:"id"`
	TenantID   string `json:"tenant_id" db:"tenant_id"`
	MaxLeave   int    `json:"max_leave" db:"max_leave"`
	MaxOvertime int   `json:"max_overtime" db:"max_overtime"`
	CreatedAt  int64  `json:"created_at" db:"created_at"`
	UpdatedAt  int64  `json:"updated_at" db:"updated_at"`
}

// Employee represents an employee in the organization.
type Employee struct {
	ID          string `json:"id" db:"id"`
	TenantID    string `json:"tenant_id" db:"tenant_id"`
	EmployeeID  string `json:"employee_id" db:"employee_id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	Position    string `json:"position" db:"position"`
	Department  string `json:"department" db:"department"`
	JoiningDate int64  `json:"joining_date" db:"joining_date"`
	CreatedAt   int64  `json:"created_at" db:"created_at"`
	UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
}

// SalaryStructure represents an employee's salary structure.
type SalaryStructure struct {
	ID            string  `json:"id" db:"id"`
	EmployeeID    string  `json:"employee_id" db:"employee_id"`
	TenantID      string  `json:"tenant_id" db:"tenant_id"`
	BasicSalary   float64 `json:"basic_salary" db:"basic_salary"`
	HRA           float64 `json:"hra" db:"hra"`
	Allowances    float64 `json:"allowances" db:"allowances"`
	Deductions    float64 `json:"deductions" db:"deductions"`
	CreatedAt     int64   `json:"created_at" db:"created_at"`
	UpdatedAt     int64   `json:"updated_at" db:"updated_at"`
}