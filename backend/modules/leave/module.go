package leave

import (
	"time"

	"hunrx.com/attendance_management/pkg/httpx"

	"uuid"

	"github.com/gin-gonic/gin"
	"hunrx.com/attendance_management/pkg/app"
)

type Module struct {
	service *Service
}

func NewModule() app.FeatureModule {
	return &Module{}
}

func (m *Module) Name() string {
	return "user"
}

func (m *Module) Mount(core *app.CoreContext, group *gin.RouterGroup) {
	repo := NewRepository(core.Pool)
	svc := NewService(repo)
	m.service = svc
	// Public routes - leave requests can be submitted by employees
	group.POST("", m.createLeaveRequest)

	// Protected routes with tenant context and role restrictions
	group.Use(httpx.RequireTenant())
	group.Use(httpx.RequireRole("EMPLOYEE", "HR_MANAGER", "ORG_ADMIN"))
	group.GET("", m.listLeaveRequests)
	group.GET("/:id", m.getLeaveRequest)

	// Approval routes (HR/Manager only)
	approvalGroup := group.Group("/approvals")
	{
		approvalGroup.Use(httpx.RequireRole("HR_MANAGER", "ORG_ADMIN"))
		approvalGroup.POST("/:id/approve", m.approveLeaveRequest)
		approvalGroup.POST("/:id/reject", m.rejectLeaveRequest)
	}

	// Employee-specific routes (can only view their own requests)
	employeeGroup := group.Group("/my")
	{
		employeeGroup.Use(httpx.RequireRole("EMPLOYEE"))
		employeeGroup.GET("", m.listMyLeaveRequests)
	}
}

func (m *Module) createLeaveRequest(c *gin.Context) {
	var req CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tenantID, exists := c.Get("org_id")
	if !exists {
		c.JSON(400, gin.H{"error": "Tenant ID not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}

	req.TenantID = tenantID.(string)
	req.EmployeeID = userID.(string)

	leave, err := m.service.CreateLeaveRequest(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, leave)
}

func (m *Module) getLeaveRequest(c *gin.Context) {
	leaveID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid leave request ID"})
		return
	}

	leave, err := m.service.GetLeaveRequest(c.Request.Context(), leaveID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if leave == nil {
		c.JSON(404, gin.H{"error": "Leave request not found"})
		return
	}

	c.JSON(200, leave)
}

func (m *Module) listLeaveRequests(c *gin.Context) {
	tenantID, exists := c.Get("org_id")
	if !exists {
		c.JSON(400, gin.H{"error": "Tenant ID not found"})
		return
	}

	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		time.Sscanf(pageStr, "%d", &page)
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		time.Sscanf(limitStr, "%d", &limit)
	}

	status := c.Query("status")

	leaves, total, err := m.service.ListLeaveRequests(c.Request.Context(), tenantID.(string), status, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"leaves": leaves,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (m *Module) approveLeaveRequest(c *gin.Context) {
	leaveID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid leave request ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}

	err = m.service.ApproveLeaveRequest(c.Request.Context(), leaveID, userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Leave request approved successfully"})
}

func (m *Module) rejectLeaveRequest(c *gin.Context) {
	leaveID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid leave request ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}

	err = m.service.RejectLeaveRequest(c.Request.Context(), leaveID, userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Leave request rejected successfully"})
}

func (m *Module) listMyLeaveRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}

	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		time.Sscanf(pageStr, "%d", &page)
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		time.Sscanf(limitStr, "%d", &limit)
	}

	leaves, total, err := m.service.ListMyLeaveRequests(c.Request.Context(), userID.(string), page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"leaves": leaves,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// Leave structures

type CreateLeaveRequest struct {
	EmployeeID  string    `json:"-"` // Set from auth context
	TenantID    string    `json:"-"` // Set from auth context
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	LeaveType   string    `json:"leave_type"` // e.g., "Casual", "Sick", "Earned"
	Description string    `json:"description"`
}

type LeaveRequest struct {
	ID          uuid.UUID  `json:"id"`
	EmployeeID  string     `json:"employee_id"`
	TenantID    string     `json:"tenant_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	LeaveType   string     `json:"leave_type"`
	Description string     `json:"description"`
	Status      string     `json:"status"` // PENDING, APPROVED, REJECTED
	ApprovedBy  *string    `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
