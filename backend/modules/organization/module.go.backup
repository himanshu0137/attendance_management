package organization

import (
	"context"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"

	"hunrx.com/attendance_management/pkg/app"
	"hunrx.com/attendance_management/pkg/db"
	"hunrx.com/attendance_management/pkg/httpx"
)

// init registers the module with the core.
func init() {
	app.Register(NewModule())
}

// Module is the feature module registration.
type Module struct {
	service *Service
}

// NewModule returns the module instance.
func NewModule() app.FeatureModule {
	return Module{}
}

// Name returns the module name.
func (m Module) Name() string {
	return "organization"
}

// Migrate runs the org schema migrations.
func (m Module) Migrate(ctx context.Context, db *db.DB) error {
	return db.Migrate(ctx, "org", "../migrations")
}

// Mount wires the handlers.
func (m Module) Mount(core *app.CoreContext, rg *gin.RouterGroup) {
	repo := NewRepository(core.Pool)
	svc := NewService(repo, core.Events, newKeycloak(core.Keycloak))
	m.service = svc
	rg.GET("/:tenant_id/policies", m.getPolicy)
	rg.POST("/:tenant_id/policies", m.setPolicy)
	rg.GET("/:tenant_id/employees", m.listEmployees)
	rg.POST("/:tenant_id/employees", m.createEmployee)
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
	employees, err := m.service.ListEmployees(c.Request.Context(), tenantID)
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
// @Param employee body CreateEmployeeReq true "Employee details"
// @Success 201 "Created"
// @Router /api/organization/{tenant_id}/employees [post]
func (m Module) createEmployee(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	req, ok := httpx.BindJSON[CreateEmployeeReq](c)
	if !ok {
		return
	}
	sub, err := m.service.CreateEmployee(c.Request.Context(), tenantID, req.Username, req.Role)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 201, gin.H{"id": sub})
}

type CreateEmployeeReq struct {
	Username string `json:"username"`
	Role     string `json:"role"`
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
