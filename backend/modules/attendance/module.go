package attendance

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
	return "attendance"
}

// Migrate runs the attendance schema migrations.
func (m Module) Migrate(ctx context.Context, db *db.DB) error {
	return db.Migrate(ctx, "attendance", "../migrations")
}

// Mount wires the handlers.
func (m Module) Mount(core *app.CoreContext, rg *gin.RouterGroup) {
	repo := NewRepository(core.Pool)
	svc := NewService(repo)
	m.service = svc
	rg.POST("/punch", m.punch)
	rg.GET("/me", m.me)
	rg.GET("/reports", m.report)
}

func (m Module) punch(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	req, ok := httpx.BindJSON[PunchReq](c)
	if !ok {
		return
	}
	log, err := m.service.Punch(c.Request.Context(), tenantID, req.EmployeeID, req.Method, req.Timestamp)
	if err != nil {
		if err == ErrMethodDisabled {
			httpx.Error(c, 403, "method_disabled", err.Error())
			return
		}
		if err == ErrOutsideShift {
			httpx.Error(c, 409, "outside_shift", err.Error())
			return
		}
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 201, log)
}

func (m Module) me(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	logs, err := m.service.GetLogs(c.Request.Context(), tenantID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 200, logs)
}

func (m Module) report(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	from, fromOK := c.Get("date_from")
	to, toOK := c.Get("date_to")
	if !fromOK || !toOK {
		httpx.Error(c, 400, "missing_range", "date_from and date_to are required")
		return
	}
	report, err := m.service.Report(c.Request.Context(), tenantID, parseDate(from.(string)), parseDate(to.(string)))
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	httpx.JSON(c, 200, report)
}

type PunchReq struct {
	EmployeeID string      `json:"employee_id"`
	Method     PunchMethod `json:"method"`
	Timestamp  time.Time   `json:"timestamp"`
}

func parseDate(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, _ = time.Parse("2006-01-02", s)
	}
	return t
}
