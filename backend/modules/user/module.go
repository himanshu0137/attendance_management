package user

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"hunrx.com/attendance_management/backend/pkg/config"
	"hunrx.com/attendance_management/backend/pkg/middleware"
)

type Module struct {
	service *Service
}

func NewModule(db *pgxpool.Pool, cfg *config.Config) (*Module, error) {
	repo := NewRepository(db)
	service := NewService(repo)
	
	return &Module{
		service: service,
	}, nil
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		// Public routes - only for registration/invite handling
		userGroup.POST("", m.createUser)
		
		// Protected routes with tenant context
		userGroup.Use(middleware.RequireTenant())
		userGroup.Use(middleware.RequireRole("ORG_ADMIN", "HR_MANAGER"))
		userGroup.GET("/:id", m.getUser)
		userGroup.PUT("/:id", m.updateUser)
		userGroup.DELETE("/:id", m.deleteUser)
		userGroup.GET("", m.listUsers)
		
		// Role management (admin only)
		roleGroup := userGroup.Group("/roles")
		{
			roleGroup.Use(middleware.RequireRole("ORG_ADMIN"))
			roleGroup.POST("", m.createRole)
			roleGroup.GET("/:id", m.getRole)
		}
		
		// Role assignment
		assignmentGroup := userGroup.Group("/assignments")
		{
			assignmentGroup.Use(middleware.RequireRole("ORG_ADMIN", "HR_MANAGER"))
			assignmentGroup.POST("", m.assignRole)
		}
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