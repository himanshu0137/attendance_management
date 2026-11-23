package organization

import (
	"time"
	"github.com/google/uuid"
)

// AllowedPunchMethods is the set of punch methods a tenant may enable.
const (
	// FACESCAN_KIOSK is the kiosk face recognition biometric.
	FACESCAN_KIOSK = "FACESCAN_KIOSK"
	// FINGERSCAN_KIOSK is the kiosk finger scan biometric.
	FINGERSCAN_KIOSK = "FINGERSCAN_KIOSK"
	// WEB_LOGIN is the Angular web portal login.
	WEB_LOGIN = "WEB_LOGIN"
	// MOBILE_LOGIN is the Flutter employee app login.
	MOBILE_LOGIN = "MOBILE_LOGIN"
)

// OrgPolicy is the per-tenant attendance policy.
type OrgPolicy struct {
	// Methods is the list of enabled punch methods.
	Methods []string `json:"methods"`
	// SalarySlipEnabled controls whether payroll may generate salary slips.
	SalarySlipEnabled bool `json:"salary_slip_enabled"`
	// ShiftStart/ShiftEnd bound the working day in local time.
	ShiftStart time.Time `json:"shift_start_time"`
	ShiftEnd   time.Time `json:"shift_end_time"`
}

// IsEnabled reports whether the given punch method is enabled for the policy.
func (p OrgPolicy) IsEnabled(method string) bool {
	for _, m := range p.Methods {
		if m == method {
			return true
		}
	}
	return false
}

// Employee is a Keycloak-backed org user.
type Employee struct {
	ID          string `json:"id"`
	KeycloakSub string `json:"keycloak_sub"`
	Role        string `json:"role"`
}

// CreateEmployeeRequest represents the request to create an employee
type CreateEmployeeRequest struct {
	TenantID     string `json:"tenant_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	EmployeeCode string `json:"employee_code"`
	Designation  string `json:"designation"`
	Department   string `json:"department"`
	ManagerID    *uuid.UUID `json:"manager_id,omitempty"`
	BranchID     *uuid.UUID `json:"branch_id,omitempty"`
}

// EmployeeWithDocuments represents an employee with their KYC documents
type EmployeeWithDocuments struct {
	*Employee
	Documents   []Document `json:"documents"`
	JoiningDate time.Time  `json:"joining_date"`
	Position    string     `json:"position"`
	Department  string     `json:"department"`
}

// EmployeeDocuments represents employee's KYC documents
type EmployeeDocuments struct {
	EmployeeID    uuid.UUID `json:"employee_id"`
	PanNumber     string    `json:"pan_number"`
	AadhaarNumber string    `json:"aadhaar_number"`
	UanNumber     string    `json:"uan_number"`
	EsicNumber    string    `json:"esic_number"`
	Documents     []Document `json:"documents"`
}

// Document represents an employee document
type Document struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// CreateSalaryStructureRequest represents the request to create a salary structure
type CreateSalaryStructureRequest struct {
	PayBasis           string    `json:"pay_basis"` // "monthly" or "daily"
	Basic              float64   `json:"basic"`
	HRA                float64   `json:"hra"`
	DA                 float64   `json:"da"`
	Conveyance         float64   `json:"conveyance"`
	SpecialAllowance   float64   `json:"special_allowance"`
	MedicalAllowance   float64   `json:"medical_allowance"`
	Bonus              float64   `json:"bonus"`
	OvertimeRate       float64   `json:"overtime_rate"`
	EffectiveDate      time.Time `json:"effective_date"`
	CreatedAt          time.Time `json:"created_at"`
}

// SalaryStructure represents an employee's salary structure
type SalaryStructure struct {
	ID                 uuid.UUID `json:"id"`
	EmployeeID         uuid.UUID `json:"employee_id"`
	TenantID           string    `json:"tenant_id"`
	PayBasis           string    `json:"pay_basis"`
	Basic              float64   `json:"basic"`
	HRA                float64   `json:"hra"`
	DA                 float64   `json:"da"`
	Conveyance         float64   `json:"conveyance"`
	SpecialAllowance   float64   `json:"special_allowance"`
	MedicalAllowance   float64   `json:"medical_allowance"`
	Bonus              float64   `json:"bonus"`
	OvertimeRate       float64   `json:"overtime_rate"`
	EffectiveDate      time.Time `json:"effective_date"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
