package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service handles business logic for organizations.
type Service struct {
	repo *Repository
}

// NewService creates a new service instance.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateOrganization creates a new organization.
func (s *Service) CreateOrganization(ctx context.Context, req CreateOrganizationRequest) (*Organization, error) {
	// Generate UUID for the organization
	orgID := uuid.New()

	// Set created timestamp
	req.CreatedAt = time.Now()

	// Create the organization in the repository
	org, err := s.repo.CreateOrganization(ctx, orgID, req)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrganization retrieves an organization by ID.
func (s *Service) GetOrganization(ctx context.Context, id uuid.UUID) (*Organization, error) {
	return s.repo.GetOrganization(ctx, id)
}

// UpdateOrganization updates an existing organization.
func (s *Service) UpdateOrganization(ctx context.Context, id uuid.UUID, req UpdateOrganizationRequest) (*Organization, error) {
	org, err := s.repo.UpdateOrganization(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// DeleteOrganization deletes an organization.
func (s *Service) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteOrganization(ctx, id)
}

// CreateBranch creates a new branch.
func (s *Service) CreateBranch(ctx context.Context, req CreateBranchRequest) (*Branch, error) {
	branchID := uuid.New()

	branch, err := s.repo.CreateBranch(ctx, branchID, req)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

// ListBranches lists branches for a tenant.
func (s *Service) ListBranches(ctx context.Context, tenantID string, page, limit int) ([]*Branch, int, error) {
	return s.repo.ListBranches(ctx, tenantID, page, limit)
}

// GetBranch retrieves a branch by ID.
func (s *Service) GetBranch(ctx context.Context, id uuid.UUID, tenantID string) (*Branch, error) {
	return s.repo.GetBranch(ctx, id, tenantID)
}

// UpdateBranch updates an existing branch.
func (s *Service) UpdateBranch(ctx context.Context, id uuid.UUID, req UpdateBranchRequest) (*Branch, error) {
	branch, err := s.repo.UpdateBranch(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

// DeleteBranch deletes a branch.
func (s *Service) DeleteBranch(ctx context.Context, id uuid.UUID, tenantID string) error {
	return s.repo.DeleteBranch(ctx, id, tenantID)
}

// CreateEmployee creates a new employee.
func (s *Service) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*Employee, error) {
	employeeID := uuid.New()

	employee, err := s.repo.CreateEmployee(ctx, employeeID, req)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

// ListEmployees lists employees for a tenant.
func (s *Service) ListEmployees(ctx context.Context, tenantID string, page, limit int) ([]*Employee, int, error) {
	return s.repo.ListEmployees(ctx, tenantID, page, limit)
}

// GetEmployee retrieves an employee by ID.
func (s *Service) GetEmployee(ctx context.Context, id uuid.UUID, tenantID string) (*Employee, error) {
	return s.repo.GetEmployee(ctx, id, tenantID)
}

// UpdateEmployee updates an existing employee.
func (s *Service) UpdateEmployee(ctx context.Context, id uuid.UUID, req UpdateEmployeeRequest) (*Employee, error) {
	employee, err := s.repo.UpdateEmployee(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

// DeleteEmployee deletes an employee.
func (s *Service) DeleteEmployee(ctx context.Context, id uuid.UUID, tenantID string) error {
	return s.repo.DeleteEmployee(ctx, id, tenantID)
}

// CreateEmployeeWithDocuments creates an employee with KYC documents and generates a unique employee code
func (s *Service) CreateEmployeeWithDocuments(ctx context.Context, req CreateEmployeeWithDocumentsRequest) (*EmployeeWithDocuments, error) {
	// Generate unique employee code
	employeeCode, err := s.generateUniqueEmployeeCode(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}

	// Set the employee code in the request
	req.EmployeeCode = employeeCode

	employeeID := uuid.New()

	// Create the employee with documents
	employee, err := s.repo.CreateEmployeeWithDocuments(ctx, employeeID, req)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

// BulkImportEmployees imports multiple employees from CSV/Excel
func (s *Service) BulkImportEmployees(ctx context.Context, tenantID string, employees []CreateEmployeeRequest) ([]*EmployeeImportResult, error) {
	var results []*EmployeeImportResult

	for _, empReq := range employees {
		result := &EmployeeImportResult{
			EmployeeCode: empReq.EmployeeCode,
			Status:       "SUCCESS",
		}

		// Check if employee code already exists
		existingEmp, err := s.repo.GetEmployeeByCode(ctx, tenantID, empReq.EmployeeCode)
		if err == nil && existingEmp != nil {
			result.Status = "CONFLICT"
			result.Error = "Employee code already exists"
			results = append(results, result)
			continue
		}

		// Create employee
		employeeID := uuid.New()
		emp, err := s.repo.CreateEmployee(ctx, employeeID, empReq)
		if err != nil {
			result.Status = "ERROR"
			result.Error = err.Error()
			results = append(results, result)
			continue
		}

		result.EmployeeID = emp.ID.String()
		results = append(results, result)
	}

	return results, nil
}

// GenerateUniqueEmployeeCode generates a unique employee code for the tenant
func (s *Service) generateUniqueEmployeeCode(ctx context.Context, tenantID string) (string, error) {
	// Simple approach: EMP-{tenant_id}-{counter}
	// In real implementation, this would be more sophisticated with database checks

	// For now, we'll use a simple pattern - in production this should check against existing codes
	return "EMP-" + tenantID + "-" + time.Now().Format("20060102150405"), nil
}

// GetEmployeeDocuments retrieves employee KYC documents
func (s *Service) GetEmployeeDocuments(ctx context.Context, employeeID uuid.UUID, tenantID string) (*EmployeeDocuments, error) {
	return s.repo.GetEmployeeDocuments(ctx, employeeID, tenantID)
}

// GenerateOfferLetter generates a dynamic offer letter PDF
func (s *Service) GenerateOfferLetter(ctx context.Context, employeeID uuid.UUID, tenantID string, template string) ([]byte, error) {
	employee, err := s.GetEmployee(ctx, employeeID, tenantID)
	if err != nil {
		return nil, err
	}

	// In a real implementation, this would use a templating engine to generate PDF
	// For now, we'll return mock data
	offerLetterContent := []byte("Offer Letter for " + employee.FirstName + " " + employee.LastName)
	return offerLetterContent, nil
}

// SetEmployeeExitStatus updates employee status to exited and handles final settlement
func (s *Service) SetEmployeeExitStatus(ctx context.Context, employeeID uuid.UUID, tenantID string, exitStatus EmployeeExitStatus) error {
	// Check if employee exists
	employee, err := s.GetEmployee(ctx, employeeID, tenantID)
	if err != nil {
		return err
	}

	// Update employee status to exited
	err = s.repo.UpdateEmployeeExitStatus(ctx, employeeID, tenantID, exitStatus)
	if err != nil {
		return err
	}

	// Process final settlement if needed
	if exitStatus.SettlementDetails != nil {
		// Handle F&F settlement logic here
		// This would include calculating net payable, updating balances, etc.
	}

	return nil
}

// CalculateFinalSettlement calculates employee's final settlement
func (s *Service) CalculateFinalSettlement(ctx context.Context, employeeID uuid.UUID, tenantID string) (*FinalSettlement, error) {
	// In a real implementation, this would calculate:
	// - Encashable leave days
	// - Pending advances/loans
	// - Net payable amount

	return &FinalSettlement{
		EncashableLeaveDays:   5,
		LeaveEncashmentAmount: 5000,
		PendingAdvanceBalance: 2000,
		NetPayable:            3000,
	}, nil
}

// CreateSalaryStructure creates a salary structure for an employee
func (s *Service) CreateSalaryStructure(ctx context.Context, employeeID uuid.UUID, tenantID string, salaryReq CreateSalaryStructureRequest) (*SalaryStructure, error) {
	// Create the salary structure in the repository
	salaryStruct, err := s.repo.CreateSalaryStructure(ctx, employeeID, tenantID, salaryReq)
	if err != nil {
		return nil, err
	}

	return salaryStruct, nil
}

// GetSalaryStructure retrieves an employee's salary structure
func (s *Service) GetSalaryStructure(ctx context.Context, employeeID uuid.UUID, tenantID string) (*SalaryStructure, error) {
	return s.repo.GetSalaryStructure(ctx, employeeID, tenantID)
}

// UpdateSalaryStructure updates an employee's salary structure
func (s *Service) UpdateSalaryStructure(ctx context.Context, employeeID uuid.UUID, tenantID string, salaryReq CreateSalaryStructureRequest) (*SalaryStructure, error) {
	salaryStruct, err := s.repo.UpdateSalaryStructure(ctx, employeeID, tenantID, salaryReq)
	if err != nil {
		return nil, err
	}

	return salaryStruct, nil
}

// OnboardEmployee handles the employee onboarding process including document upload and setup
func (s *Service) OnboardEmployee(ctx context.Context, employeeID uuid.UUID, tenantID string, req OnboardEmployeeReq) (*EmployeeWithDocuments, error) {
	// First check if employee exists
	employee, err := s.GetEmployee(ctx, employeeID, tenantID)
	if err != nil {
		return nil, err
	}

	// In a real implementation, we would:
	// 1. Store the documents in a secure location (S3, local storage, etc.)
	// 2. Update employee with onboarding details
	// 3. Set initial status to onboarded

	// For now, we'll mock this functionality and return what we have
	result := &EmployeeWithDocuments{
		Employee:    employee,
		Documents:   req.Documents,
		JoiningDate: req.JoiningDate,
		Position:    req.Position,
		Department:  req.Department,
	}

	return result, nil
}

// GenerateHRLetter generates a dynamic HR letter for an employee
func (s *Service) GenerateHRLetter(ctx context.Context, employeeID uuid.UUID, tenantID string, letterType string) (*HRLetterResponse, error) {
	employee, err := s.GetEmployee(ctx, employeeID, tenantID)
	if err != nil {
		return nil, err
	}

	// In a real implementation, this would:
	// 1. Generate a PDF using a templating engine (like Go templates or docxtemplater)
	// 2. Store it in cloud storage
	// 3. Return the URL to access it

	// For now, we'll create mock content based on letter type
	var content []byte
	switch letterType {
	case "offer":
		content = []byte("Offer Letter for " + employee.FirstName + " " + employee.LastName + "\nPosition: " + employee.Position)
	case "appointment":
		content = []byte("Appointment Letter for " + employee.FirstName + " " + employee.LastName + "\nPosition: " + employee.Position)
	default:
		content = []byte("HR Letter for " + employee.FirstName + " " + employee.LastName)
	}

	result := &HRLetterResponse{
		LetterType: letterType,
		Content:    content,
		URL:        "", // In real implementation, this would be the URL to the generated document
	}

	return result, nil
}
