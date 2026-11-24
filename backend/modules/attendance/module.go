package attendance

import (
	"time"

	"github.com/gin-gonic/gin"

	"hunrx.com/attendance_management/pkg/app"
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
	return "attendance"
}

// Mount wires the handlers.
func (m *Module) Mount(core *app.CoreContext, rg *gin.RouterGroup) {
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

type PunchReq struct {
	EmployeeID string      `json:"employee_id"`
	Method     PunchMethod `json:"method"`
	Timestamp  time.Time   `json:"timestamp"`
}
