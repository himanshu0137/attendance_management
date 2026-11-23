package attendance

import (
	"context"
	"strconv"
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
	
	// Shift management endpoints (US-06)
	shiftGroup := rg.Group("/shifts")
	{
		shiftGroup.POST("", m.createShift)
		shiftGroup.GET("", m.listShifts)
		shiftGroup.GET("/:id", m.getShift)
		shiftGroup.PUT("/:id", m.updateShift)
		shiftGroup.DELETE("/:id", m.deleteShift)
	}
	
	// Attendance rules endpoints (US-07, US-08) 
	ruleGroup := rg.Group("/rules")
	{
		ruleGroup.POST("", m.createAttendanceRule)
		ruleGroup.GET("", m.listAttendanceRules)
		ruleGroup.GET("/:id", m.getAttendanceRule)
		ruleGroup.PUT("/:id", m.updateAttendanceRule)
		ruleGroup.DELETE("/:id", m.deleteAttendanceRule)
	}
}

// PunchIn punches in for an employee
// @Summary Punch in for employee
// @Description Record a punch-in or punch-out for an employee
// @Tags Attendance
// @Accept json
// @Produce json
// @Param punch body PunchReq true "Punch request"
// @Success 201 {object} PunchLog
// @Router /api/attendance/punch [post]
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

// GetEmployeeAttendance gets current employee attendance records
// @Summary Get employee attendance records
// @Description Retrieve attendance logs for the authenticated employee
// @Tags Attendance
// @Accept json
// @Produce json
// @Success 200 {array} PunchLog
// @Router /api/attendance/me [get]
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

// GetAttendanceReport gets attendance report for a date range
// @Summary Get attendance report
// @Description Generate an attendance report for a date range
// @Tags Attendance
// @Accept json
// @Produce json
// @Param date_from query string true "Start date (YYYY-MM-DD)"
// @Param date_to query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} Report
// @Router /api/attendance/reports [get]
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

// CreateShift creates a new shift
// @Summary Create shift
// @Description Create a new shift for a tenant
// @Tags Attendance - Shifts
// @Accept json
// @Produce json
// @Param shift body Shift true "Shift details"
// @Success 201 {object} Shift
// @Router /api/attendance/shifts [post]
func (m Module) createShift(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	req, ok := httpx.BindJSON[Shift](c)
	if !ok {
		return
	}
	
	shift, err := m.service.CreateShift(c.Request.Context(), tenantID, req)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 201, shift)
}

// ListShifts lists shifts for a tenant
// @Summary List shifts
// @Description List all shifts for a tenant with pagination
// @Tags Attendance - Shifts
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {array} Shift
// @Router /api/attendance/shifts [get]
func (m Module) listShifts(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	page := 1
	limit := 20
	
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	
	shifts, total, err := m.service.ListShifts(c.Request.Context(), tenantID, page, limit)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	// Return paginated response
	httpx.JSON(c, 200, gin.H{
		"shifts": shifts,
		"page":   page,
		"limit":  limit,
		"total":  total,
	})
}

// GetShift retrieves a shift by ID
// @Summary Get shift
// @Description Retrieve a specific shift by ID
// @Tags Attendance - Shifts
// @Accept json
// @Produce json
// @Param id path string true "Shift ID"
// @Success 200 {object} Shift
// @Router /api/attendance/shifts/{id} [get]
func (m Module) getShift(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	shiftID := c.Param("id")
	
	shift, err := m.service.GetShift(c.Request.Context(), tenantID, shiftID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 200, shift)
}

// UpdateShift updates an existing shift
// @Summary Update shift
// @Description Update a shift by ID
// @Tags Attendance - Shifts
// @Accept json
// @Produce json
// @Param id path string true "Shift ID"
// @Param shift body Shift true "Shift details"
// @Success 200 {object} Shift
// @Router /api/attendance/shifts/{id} [put]
func (m Module) updateShift(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	shiftID := c.Param("id")
	
	req, ok := httpx.BindJSON[Shift](c)
	if !ok {
		return
	}
	
	shift, err := m.service.UpdateShift(c.Request.Context(), tenantID, shiftID, req)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 200, shift)
}

// DeleteShift deletes a shift
// @Summary Delete shift
// @Description Delete a shift by ID
// @Tags Attendance - Shifts
// @Accept json
// @Produce json
// @Param id path string true "Shift ID"
// @Success 204 "No Content"
// @Router /api/attendance/shifts/{id} [delete]
func (m Module) deleteShift(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	shiftID := c.Param("id")
	
	err := m.service.DeleteShift(c.Request.Context(), tenantID, shiftID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 204, gin.H{})
}

// CreateAttendanceRule creates a new attendance rule
// @Summary Create attendance rule
// @Description Create a new attendance policy rule for a tenant
// @Tags Attendance - Rules
// @Accept json
// @Produce json
// @Param rule body AttendanceRule true "Attendance rule details"
// @Success 201 {object} AttendanceRule
// @Router /api/attendance/rules [post]
func (m Module) createAttendanceRule(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	req, ok := httpx.BindJSON[AttendanceRule](c)
	if !ok {
		return
	}
	
	rule, err := m.service.CreateAttendanceRule(c.Request.Context(), tenantID, req)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 201, rule)
}

// ListAttendanceRules lists attendance rules for a tenant
// @Summary List attendance rules
// @Description List all attendance rules for a tenant with pagination
// @Tags Attendance - Rules
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Success 200 {array} AttendanceRule
// @Router /api/attendance/rules [get]
func (m Module) listAttendanceRules(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	page := 1
	limit := 20
	
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	
	rules, total, err := m.service.ListAttendanceRules(c.Request.Context(), tenantID, page, limit)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	// Return paginated response
	httpx.JSON(c, 200, gin.H{
		"rules": rules,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetAttendanceRule retrieves an attendance rule by ID
// @Summary Get attendance rule
// @Description Retrieve a specific attendance rule by ID
// @Tags Attendance - Rules
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Success 200 {object} AttendanceRule
// @Router /api/attendance/rules/{id} [get]
func (m Module) getAttendanceRule(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	ruleID := c.Param("id")
	
	rule, err := m.service.GetAttendanceRule(c.Request.Context(), tenantID, ruleID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 200, rule)
}

// UpdateAttendanceRule updates an existing attendance rule
// @Summary Update attendance rule
// @Description Update an attendance rule by ID
// @Tags Attendance - Rules
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Param rule body AttendanceRule true "Attendance rule details"
// @Success 200 {object} AttendanceRule
// @Router /api/attendance/rules/{id} [put]
func (m Module) updateAttendanceRule(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	ruleID := c.Param("id")
	
	req, ok := httpx.BindJSON[AttendanceRule](c)
	if !ok {
		return
	}
	
	rule, err := m.service.UpdateAttendanceRule(c.Request.Context(), tenantID, ruleID, req)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 200, rule)
}

// DeleteAttendanceRule deletes an attendance rule
// @Summary Delete attendance rule
// @Description Delete an attendance rule by ID
// @Tags Attendance - Rules
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Success 204 "No Content"
// @Router /api/attendance/rules/{id} [delete]
func (m Module) deleteAttendanceRule(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	ruleID := c.Param("id")
	
	err := m.service.DeleteAttendanceRule(c.Request.Context(), tenantID, ruleID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 204, gin.H{})
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
