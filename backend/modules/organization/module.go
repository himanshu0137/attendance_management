package organization

import (
	"errors"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hunrx.com/attendance_management/pkg/app"
	"hunrx.com/attendance_management/pkg/httpx"
)

// Module is the feature module registration.
type Module struct {
	service *Service
}

// NewModule returns the module instance.
func NewModule() app.FeatureModule {
	return &Module{}
}

// Name returns the module name.
func (m *Module) Name() string {
	return "organization"
}

// Mount wires the handlers.
func (m *Module) Mount(core *app.CoreContext, rg *gin.RouterGroup) {
	repo := NewRepository(core.Pool)
	svc := NewService(repo)
	m.service = svc
	rg.GET("/:tenant_id/policies", m.getPolicy)
	rg.POST("/:tenant_id/policies", m.setPolicy)
	rg.GET("/:tenant_id/employees", m.listEmployees)
	rg.POST("/:tenant_id/employees", m.createEmployee)

	// Employee Lifecycle Management endpoints (US-03, US-04, US-05)
	rg.POST("/:tenant_id/employees/:employee_id/onboard", m.onboardEmployee)
	rg.GET("/:tenant_id/employees/:employee_id/salary-structure", m.getSalaryStructure)
	rg.POST("/:tenant_id/employees/:employee_id/salary-structure", m.createSalaryStructure)
	rg.PUT("/:tenant_id/employees/:employee_id/salary-structure", m.updateSalaryStructure)
	rg.POST("/:tenant_id/employees/:employee_id/hr-letters/:letter_type", m.generateHRLetter)
}

// GetOrganizationPolicy gets organization policy settings
// @Summary Get organization policy
// @Description Retrieve organization policy settings for a tenant
// @Tags Organization
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {object} OrgPolicy
// @Router /api/organization/{tenant_id}/policies [get]
func (m Module) getPolicy(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	p, err := m.service.GetPolicy(c.Request.Context(), tenantID)
	if err != nil {
		if errors.Is(err, ErrPolicyNotFound{}) {
			httpx.JSON(c, 404, gin.H{"error": err.Error()})
			return
		}
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 200, p)
}

// SetOrganizationPolicy sets organization policy settings
// @Summary Set organization policy
// @Description Update organization policy settings for a tenant
// @Tags Organization
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param policy body OrgPolicy true "Policy settings"
// @Success 200 {object} OrgPolicy
// @Router /api/organization/{tenant_id}/policies [post]
func (m Module) setPolicy(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	p, ok := httpx.BindJSON[OrgPolicy](c)
	if !ok {
		return
	}
	if err := m.service.SetPolicy(c.Request.Context(), tenantID, p); err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 200, gin.H{"policy": p})
}

// ListEmployees lists organization employees
// @Summary List employees
// @Description Retrieve list of employees for a tenant
// @Tags Organization
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {array} Employee
// @Router /api/organization/{tenant_id}/employees [get]
func (m Module) listEmployees(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	employees, _, err := m.service.ListEmployees(c.Request.Context(), tenantID, 10, 10)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 200, employees)
}

// CreateEmployee creates a new employee
// @Summary Create employee
// @Description Create a new employee for a tenant
// @Tags Organization
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param employee body CreateEmployeeRequest true "Employee details"
// @Success 201 "Created"
// @Router /api/organization/{tenant_id}/employees [post]
func (m Module) createEmployee(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	req, ok := httpx.BindJSON[CreateEmployeeRequest](c)
	if !ok {
		return
	}

	// Set the tenant ID for the request
	req.TenantID = tenantID

	emp, err := m.service.CreateEmployee(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 201, gin.H{"id": emp.ID})
}

// OnboardEmployee handles employee onboarding process
// @Summary Onboard employee
// @Description Complete employee onboarding with documents and initial setup
// @Tags Organization - Employee Lifecycle
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param employee_id path string true "Employee ID"
// @Param onboarding body OnboardEmployeeReq true "Onboarding details"
// @Success 200 {object} EmployeeWithDocuments
// @Router /api/organization/{tenant_id}/employees/{employee_id}/onboard [post]
func (m Module) onboardEmployee(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}

	employeeIDStr := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		httpx.Error(c, 400, "invalid_employee_id", "Invalid employee ID format")
		return
	}

	req, ok := httpx.BindJSON[OnboardEmployeeReq](c)
	if !ok {
		return
	}

	// Call service method to onboard employee
	result, err := m.service.OnboardEmployee(c.Request.Context(), employeeID, tenantID, req)
	if err != nil {
		httpx.Error(c, 500, "onboarding_error", err.Error())
		return
	}

	httpx.JSON(c, 200, result)
}

// GetSalaryStructure retrieves employee's salary structure
// @Summary Get salary structure
// @Description Retrieve an employee's salary structure details
// @Tags Organization - Employee Lifecycle
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} SalaryStructure
// @Router /api/organization/{tenant_id}/employees/{employee_id}/salary-structure [get]
func (m Module) getSalaryStructure(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}

	employeeIDStr := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		httpx.Error(c, 400, "invalid_employee_id", "Invalid employee ID format")
		return
	}

	result, err := m.service.GetSalaryStructure(c.Request.Context(), employeeID, tenantID)
	if err != nil {
		httpx.Error(c, 500, "salary_structure_error", err.Error())
		return
	}

	httpx.JSON(c, 200, result)
}

// CreateSalaryStructure creates a new salary structure for an employee
// @Summary Create salary structure
// @Description Create a new salary structure for an employee
// @Tags Organization - Employee Lifecycle
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param employee_id path string true "Employee ID"
// @Param salary_structure body CreateSalaryStructureRequest true "Salary structure details"
// @Success 201 {object} SalaryStructure
// @Router /api/organization/{tenant_id}/employees/{employee_id}/salary-structure [post]
func (m Module) createSalaryStructure(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}

	employeeIDStr := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		httpx.Error(c, 400, "invalid_employee_id", "Invalid employee ID format")
		return
	}

	req, ok := httpx.BindJSON[CreateSalaryStructureRequest](c)
	if !ok {
		return
	}

	result, err := m.service.CreateSalaryStructure(c.Request.Context(), employeeID, tenantID, req)
	if err != nil {
		httpx.Error(c, 500, "salary_structure_error", err.Error())
		return
	}

	httpx.JSON(c, 201, result)
}

// UpdateSalaryStructure updates an employee's salary structure
// @Summary Update salary structure
// @Description Update an employee's existing salary structure
// @Tags Organization - Employee Lifecycle
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param employee_id path string true "Employee ID"
// @Param salary_structure body CreateSalaryStructureRequest true "Salary structure details"
// @Success 200 {object} SalaryStructure
// @Router /api/organization/{tenant_id}/employees/{employee_id}/salary-structure [put]
func (m Module) updateSalaryStructure(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}

	employeeIDStr := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		httpx.Error(c, 400, "invalid_employee_id", "Invalid employee ID format")
		return
	}

	req, ok := httpx.BindJSON[CreateSalaryStructureRequest](c)
	if !ok {
		return
	}

	result, err := m.service.UpdateSalaryStructure(c.Request.Context(), employeeID, tenantID, req)
	if err != nil {
		httpx.Error(c, 500, "salary_structure_error", err.Error())
		return
	}

	httpx.JSON(c, 200, result)
}

// GenerateHRLetter generates a dynamic HR letter for an employee
// @Summary Generate HR letter
// @Description Generate a dynamic HR letter (offer letter, appointment letter, etc.)
// @Tags Organization - Employee Lifecycle
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param employee_id path string true "Employee ID"
// @Param letter_type path string true "Letter type (offer, appointment, etc.)"
// @Success 200 {object} HRLetterResponse
// @Router /api/organization/{tenant_id}/employees/{employee_id}/hr-letters/{letter_type} [post]
func (m Module) generateHRLetter(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}

	employeeIDStr := c.Param("employee_id")
	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		httpx.Error(c, 400, "invalid_employee_id", "Invalid employee ID format")
		return
	}

	letterType := c.Param("letter_type")

	// Call service method to generate HR letter
	result, err := m.service.GenerateHRLetter(c.Request.Context(), employeeID, tenantID, letterType)
	if err != nil {
		httpx.Error(c, 500, "letter_generation_error", err.Error())
		return
	}

	httpx.JSON(c, 200, result)
}

// OnboardEmployeeReq represents the request for employee onboarding
type OnboardEmployeeReq struct {
	Documents   []Document `json:"documents"`
	JoiningDate time.Time  `json:"joining_date"`
	Position    string     `json:"position"`
	Department  string     `json:"department"`
}

// newKeycloak adapts the core KeycloakConfig into the org KeycloakClient.
func newKeycloak(cfg *httpx.KeycloakConfig) *KeycloakClient {
	return &KeycloakClient{
		Issuer:       cfg.Issuer,
		Realm:        realmFromIssuer(cfg.Issuer),
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	}
}

// realmFromIssuer extracts the realm from a Keycloak issuer URL.
func realmFromIssuer(issuer string) string {
	m := regexp.MustCompile(`realms/([^/]+)/protocol/openid-connect`)
	if mm := m.FindStringSubmatch(issuer); mm != nil {
		return mm[1]
	}
	return "master"
}
