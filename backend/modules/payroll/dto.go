package payroll

// SalaryRecordRequest represents a request to create or update salary record.
type SalaryRecordRequest struct {
	EmployeeID    string  `json:"employee_id"`
	SalaryMonth   int64   `json:"salary_month"`
	BasicSalary   float64 `json:"basic_salary"`
	HRA           float64 `json:"hra"`
	Allowances    float64 `json:"allowances"`
	Deductions    float64 `json:"deductions"`
	NetSalary     float64 `json:"net_salary"`
	PayDate       int64   `json:"pay_date"`
}

// PayrollSummaryRequest represents a request to create payroll summary.
type PayrollSummaryRequest struct {
	TenantID     string  `json:"tenant_id"`
	SalaryMonth  int64   `json:"salary_month"`
	TotalSalary  float64 `json:"total_salary"`
}

// PayrollReportRequest represents a request for payroll report.
type PayrollReportRequest struct {
	TenantID    string `json:"tenant_id"`
	FromDate    int64  `json:"from_date"`
	ToDate      int64  `json:"to_date"`
}

// PayrollReportResponse represents a payroll report response.
type PayrollReportResponse struct {
	ID            string  `json:"id"`
	EmployeeID    string  `json:"employee_id"`
	SalaryMonth   int64   `json:"salary_month"`
	BasicSalary   float64 `json:"basic_salary"`
	HRA           float64 `json:"hra"`
	Allowances    float64 `json:"allowances"`
	Deductions    float64 `json:"deductions"`
	NetSalary     float64 `json:"net_salary"`
	PayDate       int64   `json:"pay_date"`
}