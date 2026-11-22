package payroll

import (
	"context"
	"time"

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
	return "payroll"
}

// Migrate runs the payroll schema migrations.
func (m Module) Migrate(ctx context.Context, db *db.DB) error {
	return db.Migrate(ctx, "payroll", "../migrations")
}

// Mount wires the handlers.
func (m Module) Mount(core *app.CoreContext, rg *gin.RouterGroup) {
	svc := NewService(NewRepository(core.Pool))
	m.service = svc
	rg.POST("/slips/generate", m.generate)
	rg.GET("/slips/me", m.mySlips)
	rg.GET("/slips", m.slips)
}

func (m Module) generate(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	req, ok := httpx.BindJSON[GenerateSlipReq](c)
	if !ok {
		return
	}
	slip, err := m.service.GenerateSlip(c.Request.Context(), tenantID, req.EmployeeID, time.Month(req.Month), req.Year)
	if err != nil {
		if err == ErrNotEligible {
			httpx.Error(c, 403, "not_eligible", err.Error())
			return
		}
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	if err := m.service.InsertSlip(c.Request.Context(), tenantID, req.EmployeeID, slip.Period, slip); err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 201, slip)
}

func (m Module) mySlips(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	slips, err := m.service.GetSlips(c.Request.Context(), tenantID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 200, slips)
}

func (m Module) slips(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	slips, err := m.service.GetSlips(c.Request.Context(), tenantID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 200, slips)
}

type GenerateSlipReq struct {
	EmployeeID string `json:"employee_id"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
}
