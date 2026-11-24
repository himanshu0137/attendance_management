package user

import (
	"fmt"
	"uuid"

	"github.com/gin-gonic/gin"
	"hunrx.com/attendance_management/pkg/app"
	"hunrx.com/attendance_management/pkg/httpx"
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
	app.Register(m)
	// Public routes - only for registration/invite handling
	group.POST("", m.createUser)

	group.Use(httpx.RequireRole("ORG_ADMIN", "HR_MANAGER"))
	group.GET("/:id", m.getUser)
	group.PUT("/:id", m.updateUser)
	group.DELETE("/:id", m.deleteUser)
	group.GET("", m.listUsers)

	// Role management (admin only)
	roleGroup := group.Group("/roles")
	{
		roleGroup.Use(httpx.RequireRole("ORG_ADMIN"))
		roleGroup.POST("", m.createRole)
		roleGroup.GET("/:id", m.getRole)
	}

	// Role assignment
	assignmentGroup := group.Group("/assignments")
	{
		assignmentGroup.Use(httpx.RequireRole("ORG_ADMIN", "HR_MANAGER"))
		assignmentGroup.POST("", m.assignRole)
	}
}

func (m *Module) createUser(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(400, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req.TenantID = tenantID.(uuid.UUID)

	user, err := m.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

func (m *Module) getUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := m.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

func (m *Module) updateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = m.service.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

func (m *Module) deleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	err = m.service.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func (m *Module) listUsers(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(400, gin.H{"error": "Tenant ID not found"})
		return
	}

	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	users, total, err := m.service.ListUsers(c.Request.Context(), tenantID.(uuid.UUID), page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (m *Module) createRole(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(400, gin.H{"error": "Tenant ID not found"})
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	role, err := m.service.CreateRole(c.Request.Context(), tenantID.(tenantID).ID, req.Name, req.Description, req.Permissions)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, role)
}

func (m *Module) getRole(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := m.service.GetRole(c.Request.Context(), roleID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if role == nil {
		c.JSON(404, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(200, role)
}

func (m *Module) assignRole(c *gin.Context) {
	var req struct {
		UserID uuid.UUID `json:"user_id"`
		RoleID uuid.UUID `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(400, gin.H{"error": "Tenant ID not found"})
		return
	}

	err := m.service.AssignRole(c.Request.Context(), req.UserID, req.RoleID, tenantID.(tenantID).ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Role assigned successfully"})
}
