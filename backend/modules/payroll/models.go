package payroll

// SalaryRecord represents an employee's salary record.
type SalaryRecord struct {
	ID            string  `json:"id" db:"id"`
	EmployeeID    string  `json:"employee_id" db:"employee_id"`
	TenantID      string  `json:"tenant_id" db:"tenant_id"`
	SalaryMonth   int64   `json:"salary_month" db:"salary_month"` // Unix timestamp of the salary month
	BasicSalary   float64 `json:"basic_salary" db:"basic_salary"`
	HRA           float64 `json:"hra" db:"hra"`
	Allowances    float64 `json:"allowances" db:"allowances"`
	Deductions    float64 `json:"deductions" db:"deductions"`
	NetSalary     float64 `json:"net_salary" db:"net_salary"`
	PayDate       int64   `json:"pay_date" db:"pay_date"`
	CreatedAt     int64   `json:"created_at" db:"created_at"`
	UpdatedAt     int64   `json:"updated_at" db:"updated_at"`
}

// PayrollSummary represents a payroll summary for a period.
type PayrollSummary struct {
	ID           string  `json:"id" db:"id"`
	TenantID     string  `json:"tenant_id" db:"tenant_id"`
	SalaryMonth  int64   `json:"salary_month" db:"salary_month"`
	TotalEmployees int   `json:"total_employees" db:"total_employees"`
	TotalSalary  float64 `json:"total_salary" db:"total_salary"`
	Status       string  `json:"status" db:"status"` // "PENDING", "PAID"
	CreatedAt    int64   `json:"created_at" db:"created_at"`
	UpdatedAt    int64   `json:"updated_at" db:"updated_at"`
}

// TaxBracket represents a tax bracket for income tax calculation.
type TaxBracket struct {
	ID          string  `json:"id" db:"id"`
	TenantID    string  `json:"tenant_id" db:"tenant_id"`
	MinIncome   float64 `json:"min_income" db:"min_income"`
	MaxIncome   float64 `json:"max_income" db:"max_income"`
	TaxRate     float64 `json:"tax_rate" db:"tax_rate"`
	CreatedAt   int64   `json:"created_at" db:"created_at"`
	UpdatedAt   int64   `json:"updated_at" db:"updated_at"`
}