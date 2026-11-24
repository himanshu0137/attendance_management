package organization

// CreateEmployeeRequest represents the request to create a new employee.
type CreateEmployeeRequest struct {
	TenantID   string `json:"tenant_id"`
	EmployeeID string `json:"employee_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Position   string `json:"position"`
	Department string `json:"department"`
	JoiningDate int64 `json:"joining_date"`
}

// OnboardEmployeeReq represents the request for employee onboarding.
type OnboardEmployeeReq struct {
	Documents   []Document `json:"documents"`
	JoiningDate int64      `json:"joining_date"`
	Position    string     `json:"position"`
	Department  string     `json:"department"`
}

// Document represents an employee document.
type Document struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// EmployeeWithDocuments represents an employee with their documents.
type EmployeeWithDocuments struct {
	EmployeeID  string    `json:"employee_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Position    string    `json:"position"`
	Department  string    `json:"department"`
	JoiningDate int64     `json:"joining_date"`
	Documents   []Document `json:"documents"`
}

// CreateSalaryStructureRequest represents the request to create a salary structure.
type CreateSalaryStructureRequest struct {
	BasicSalary float64 `json:"basic_salary"`
	HRA         float64 `json:"hra"`
	Allowances  float64 `json:"allowances"`
	Deductions  float64 `json:"deductions"`
}

// SalaryStructureResponse represents a salary structure response.
type SalaryStructureResponse struct {
	ID            string  `json:"id"`
	EmployeeID    string  `json:"employee_id"`
	BasicSalary   float64 `json:"basic_salary"`
	HRA           float64 `json:"hra"`
	Allowances    float64 `json:"allowances"`
	Deductions    float64 `json:"deductions"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
}

// HRLetterResponse represents a generated HR letter response.
type HRLetterResponse struct {
	LetterType string `json:"letter_type"`
	Content    string `json:"content"`
	URL        string `json:"url"`
}